package xml

import (
	"testing"
)

func TestReader(t *testing.T) {
	xml := "<foo></foo>"
	reader, err := NewReader([]byte(xml), []byte("utf-8"), nil, 0);
	if err != nil {
		t.Error(err)
	}

	if err := reader.Read(); err != nil {
		t.Error(err)
	}

	name := reader.Name()
	if name != "foo" {
		t.Errorf("unexpected name: %s", name)
	}
}
