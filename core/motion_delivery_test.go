package core

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMotionDocsMatchExampleBodies(t *testing.T) {
	src, err := os.ReadFile("../feedback/example_motion_test.go")
	if err != nil {
		t.Fatal(err)
	}
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "../feedback/example_motion_test.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	examples := map[string]string{}
	for _, d := range file.Decls {
		fn, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		body := string(src[fs.Position(fn.Body.Lbrace).Offset+1 : fs.Position(fn.Body.Rbrace).Offset])
		lines := strings.Split(strings.Trim(body, "\n"), "\n")
		for i := range lines {
			lines[i] = strings.TrimPrefix(lines[i], "\t")
		}
		examples[strings.TrimPrefix(fn.Name.Name, "Example")] = strings.TrimSpace(strings.Join(lines, "\n"))
	}
	count := 0
	for path, doc := range readGlob(t, "../docs/components/*.md") {
		id := filepath.Base(path)[:3]
		if id < "281" || id > "300" {
			continue
		}
		name := strings.Fields(strings.SplitN(doc, "\n", 2)[0])[2]
		block := docsGoBlock.FindStringSubmatch(doc)
		if len(block) != 2 || strings.TrimSpace(block[1]) != examples[name] {
			t.Errorf("%s differs from Example%s body", path, name)
		}
		for _, section := range []string{"用途", "最小示例", "参数", "状态", "事件", "键盘操作", "限制"} {
			if !strings.Contains(doc, "## "+section) {
				t.Errorf("%s missing %s", path, section)
			}
		}
		count++
	}
	if count != 20 || len(examples) != 20 {
		t.Fatalf("docs %d examples %d", count, len(examples))
	}
}
