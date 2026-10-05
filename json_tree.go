package mujica

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

// jsonNode is one value of a parsed document, in document order.
type jsonNode struct {
	key, path string
	// value is a scalar's JSON text, or "{}" / "[]" for containers.
	value  string
	depth  int
	parent int
	// end is the index after this node's subtree; count is its child count.
	end, count int
}

func (n jsonNode) container() bool { return n.value == "{}" || n.value == "[]" }

func (n jsonNode) label() string {
	if n.key == "" {
		return n.value
	}
	return n.key + ": " + n.value
}

// jsonTree is a document flattened depth-first.
type jsonTree struct{ nodes []jsonNode }

var jsonIdent = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

// parseJSONTree flattens data keeping key order; it fails on invalid or trailing JSON.
func parseJSONTree(data string) (jsonTree, error) {
	dec := json.NewDecoder(bytes.NewReader([]byte(data)))
	dec.UseNumber()
	var t jsonTree
	if err := jsonParseValue(dec, &t, "", "$", 0, -1); err != nil {
		return jsonTree{}, err
	}
	if _, err := dec.Token(); err != io.EOF {
		return jsonTree{}, fmt.Errorf("unexpected data after the value")
	}
	return t, nil
}

func jsonParseValue(dec *json.Decoder, t *jsonTree, key, path string, depth, parent int) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	at := len(t.nodes)
	t.nodes = append(t.nodes, jsonNode{key: key, path: path, depth: depth, parent: parent})
	switch v := tok.(type) {
	case json.Delim:
		t.nodes[at].value = map[json.Delim]string{'{': "{}", '[': "[]"}[v]
		for i := 0; dec.More(); i++ {
			childKey, childPath := "", fmt.Sprintf("%s[%d]", path, i)
			if v == '{' {
				kt, err := dec.Token()
				if err != nil {
					return err
				}
				childKey = kt.(string)
				childPath = path + "." + childKey
				if !jsonIdent.MatchString(childKey) {
					childPath = path + "[" + strconv.Quote(childKey) + "]"
				}
			}
			if err := jsonParseValue(dec, t, childKey, childPath, depth+1, at); err != nil {
				return err
			}
			t.nodes[at].count++
		}
		if _, err := dec.Token(); err != nil { // the closing delimiter
			return err
		}
	case string:
		t.nodes[at].value = strconv.Quote(v)
	case json.Number:
		t.nodes[at].value = v.String()
	case bool:
		t.nodes[at].value = strconv.FormatBool(v)
	case nil:
		t.nodes[at].value = "null"
	}
	t.nodes[at].end = len(t.nodes)
	return nil
}
