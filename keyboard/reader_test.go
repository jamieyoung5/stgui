package keyboard_test

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/jamieyoung5/stgui/keyboard"
)

func TestReadInput_BasicKeys(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	reader := bufio.NewReader(r)

	go func() {
		w.WriteString("a")
		w.WriteString("\n")
	}()

	got, err := keyboard.ReadInput(reader, r)
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	if got != "a" {
		t.Errorf("Expected 'a', got %q", got)
	}

	got, err = keyboard.ReadInput(reader, r)
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}
	if got != keyboard.KeyEnter {
		t.Errorf("Expected EnterKey, got %q", got)
	}
}

func TestReadInput_NamedKeys(t *testing.T) {
	cases := map[string]keyboard.Key{
		"\t":   keyboard.KeyTab,
		"\x03": keyboard.KeyCtrlC,
		"\x7f": keyboard.KeyBackspace,
		"\b":   keyboard.KeyBackspace,
	}

	for typed, want := range cases {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		reader := bufio.NewReader(r)
		go func() {
			w.WriteString(typed)
			w.Close()
		}()

		got, err := keyboard.ReadInput(reader, r)
		if err != nil {
			t.Fatalf("Failed to read %q: %v", typed, err)
		}
		if got != want {
			t.Errorf("Reading %q: expected %q, got %q", typed, want, got)
		}

		r.Close()
	}
}

func TestReadInput_Sequences(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	reader := bufio.NewReader(r)

	go func() {
		w.WriteString("\x1b[A")
	}()

	got, err := keyboard.ReadInput(reader, r)
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}

	if got != keyboard.KeyUp {
		t.Errorf("Expected UpArrowKey, got %q", got)
	}
}

func TestReadInput_EscapeTimeout(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	reader := bufio.NewReader(r)

	go func() {
		w.WriteString("\x1b")
	}()

	start := time.Now()
	got, err := keyboard.ReadInput(reader, r)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}

	if got != keyboard.KeyEscape {
		t.Errorf("Expected EscapeKey, got %q", got)
	}

	if elapsed < 200*time.Millisecond {
		t.Errorf("Expected timeout wait (~200ms), but returned in %v", elapsed)
	}
}

func TestReadInput_UnknownSequence(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()

	reader := bufio.NewReader(r)

	rawSeq := "\x1b[Z"
	go func() {
		w.WriteString(rawSeq)
	}()

	got, err := keyboard.ReadInput(reader, r)
	if err != nil {
		t.Fatalf("Failed to read input: %v", err)
	}

	if string(got) != rawSeq {
		t.Errorf("Expected raw sequence %q, got %q", rawSeq, got)
	}
}
