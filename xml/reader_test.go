package xml

import (
	"io"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	file, err := os.Open("tests/reader/books.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	r, err := NewReader(file, DefaultEncodingBytes, nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	for {
		if err := r.Read(); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		depth := r.Depth()
		nodeType := r.NodeType()
		//name := r.Name()

		if nodeType == XML_ELEMENT_NODE && depth == 1 {
			node, err := r.Expand()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("%#v", node.Name())
		}
	}
}
