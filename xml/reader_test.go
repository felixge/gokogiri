package xml

import (
	"fmt"
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
		t.Fatalf("expected: %s, got: %s", expected, got)
	}
}

func TestReader_ExpandFlow(t *testing.T) {
	pr, pw := io.Pipe()
	go pw.Write([]byte("<root>"))

	r, err := NewReader(pr, DefaultEncodingBytes, nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	expr := xpath.Compile("self::child/subchild")

	results := []string{}
	expected := []string{}
	for i := 0; i < 10; i++ {
		expected = append(expected, fmt.Sprintf("%d", i))

		go func(i int) {
			fmt.Printf("Writing: %d\n", i)
			pw.Write([]byte(fmt.Sprintf("<child><subchild>%d</subchild></child>", i)))
			// Without this, the test does not pass : (
			//pw.Write([]byte("<wat/>"))
		}(i)

		for {
			if err := r.Read(); err == io.EOF {
				t.Fatal(err)
			}

			depth := r.Depth()
			nodeType := r.NodeType()
			name := r.Name()

			if nodeType == XML_ELEMENT_NODE && depth == 1 && name == "child"  {
				fmt.Printf("Expanding: %d\n", i)
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
				break
			}
		}
	}

	if strings.Join(expected, ", ") != strings.Join(results, ", ") {
		t.Fatalf("expected: %+v, got: %+v", expected, results)
	}
}
