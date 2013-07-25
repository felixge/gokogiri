package xml

import (
	"github.com/moovweb/gokogiri/xpath"
	"io"
	"os"
	"strings"
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

	expr := xpath.Compile("self::book/title")

	results := []string{}
	for {
		if err := r.Read(); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		depth := r.Depth()
		nodeType := r.NodeType()

		if nodeType == XML_ELEMENT_NODE && depth == 1 {
			node, err := r.Expand()
			if err != nil {
				t.Fatal(err)
			}

			nodes, err := node.Search(expr)
			if err != nil {
				t.Fatal(err)
			}

			for _, node := range nodes {
				results = append(results, node.Content())
			}
		}
	}

	got := strings.Join(results, ", ")
	expected := strings.Join([]string{"XML Developer's Guide", "Midnight Rain", "Maeve Ascendant", "Oberon's Legacy", "The Sundered Grail", "Lover Birds", "Splish Splash", "Creepy Crawlies", "Paradox Lost", "Microsoft .NET: The Programming Bible", "MSXML3: A Comprehensive Guide", "Visual Studio 7: A Comprehensive Guide"}, ", ")

	if expected != got {
		t.Fatal("expected: %s, got: %s", expected, got)
	}
}
