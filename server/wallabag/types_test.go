package wallabag

import (
	"encoding/json"
	"testing"
)

func TestIntBoolUnmarshal(t *testing.T) {
	tests := []struct {
		in   string
		want IntBool
	}{
		{"1", true},
		{"0", false},
	}

	var b IntBool
	for _, it := range tests {
		if err := json.Unmarshal([]byte(it.in), &b); err != nil {
			t.Fatal(err)
		}
		if b != it.want {
			t.Fatalf("got %v, want %v", b, it.want)
		}
	}
}

func TestIntBoolMarshal(t *testing.T) {
	tests := []struct {
		in   IntBool
		want string
	}{
		{true, "1"},
		{false, "0"},
	}

	for _, it := range tests {
		b, err := json.Marshal(it.in)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != it.want {
			t.Fatalf("got %v, want %v", string(b), it.want)
		}
	}
}

func TestMagicInt(t *testing.T) {
	tests := []struct {
		in   string
		want MagicInt
	}{
		{"null", MagicInt{true, 0}},
		{"1", MagicInt{false, 1}},
		{"\"1\"", MagicInt{false, 1}},
	}

	var i MagicInt
	for _, it := range tests {
		if err := json.Unmarshal([]byte(it.in), &i); err != nil {
			t.Fatal(err)
		}
		if i != it.want {
			t.Fatalf("got %v, want %v", i, it.want)
		}
	}
}
