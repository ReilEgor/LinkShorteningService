package shortener

import "testing"

func TestEncodeDecode(t *testing.T) {
	tests := []struct {
		id   uint64
		code string
	}{
		{0, "a"},
		{1, "b"},
		{61, "9"},
		{12345, "dnh"},
	}

	for _, tt := range tests {
		gotCode := Encode(tt.id)
		if gotCode != tt.code {
			t.Errorf("Encode(%d) = %v; want %v", tt.id, gotCode, tt.code)
		}

		gotID, _ := Decode(tt.code)
		if gotID != tt.id {
			t.Errorf("Decode(%s) = %v; want %v", tt.code, gotID, tt.id)
		}
	}
}
