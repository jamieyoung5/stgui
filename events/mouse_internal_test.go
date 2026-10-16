package events

import "testing"

func TestParseMouse(t *testing.T) {
	cases := []struct {
		seq                 string
		wantCol, wantRow    int
		wantPressed, wantOK bool
	}{
		{"\x1b[<0;1;1M", 0, 0, true, true},
		{"\x1b[<0;15;7M", 14, 6, true, true},
		{"\x1b[<0;15;7m", 14, 6, false, true},
		{"\x1b[<2;100;40M", 99, 39, true, true},
		{"\x1b[<35;15;7M", 0, 0, false, false}, // motion, bit 5
		{"\x1b[<64;15;7M", 0, 0, false, false}, // wheel, bit 6
		{"\x1b[A", 0, 0, false, false},         // an arrow key
		{"\x1b[<0;15M", 0, 0, false, false},    // too few fields
		{"\x1b[<0;x;7M", 0, 0, false, false},   // not a number
		{"\x1b[<0;15;7", 0, 0, false, false},   // no terminator
	}

	for _, c := range cases {
		got, ok := parseMouse(c.seq)
		if ok != c.wantOK {
			t.Errorf("parseMouse(%q) ok = %v want %v", c.seq, ok, c.wantOK)
			continue
		}
		if ok && (got.Col != c.wantCol || got.Row != c.wantRow || got.Pressed != c.wantPressed) {
			t.Errorf("parseMouse(%q) = col %d row %d pressed %v, want %d %d %v",
				c.seq, got.Col, got.Row, got.Pressed, c.wantCol, c.wantRow, c.wantPressed)
		}
	}
}
