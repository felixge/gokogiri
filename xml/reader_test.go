package xml

import (
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	xml := strings.NewReader("<foo></foo>")
	reader, err := NewReader(xml, []byte("utf-8"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	if err := reader.Read(); err != nil {
		t.Fatal(err)
	}

	name := reader.Name()
	if name != "foo" {
		t.Errorf("unexpected name: %s", name)
	}
}
