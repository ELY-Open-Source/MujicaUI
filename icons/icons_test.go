package icons

import "testing"

func TestIconsParse(t *testing.T) {
	if len(Names()) < 60 {
		t.Fatalf("got %d icons", len(Names()))
	}
	if _, ok := Get("check"); !ok {
		t.Fatal("check missing")
	}
	if _, ok := Get("nope"); ok {
		t.Fatal("unknown name found")
	}
}
