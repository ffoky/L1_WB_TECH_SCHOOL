package cut

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func loadBigFile(tb testing.TB) []byte {
	tb.Helper()
	const path = "testdata/big.txt"
	if _, err := os.Stat(path); err != nil {
		tb.Skip("big.txt not found, generate it first")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		tb.Fatal(err)
	}
	return data
}

func TestCutBigFile(t *testing.T) {
	data := loadBigFile(t)
	opts := Options{Fields: []int{1, 3, 4, 5}, Delimiter: "\t"}

	var buf bytes.Buffer
	if err := Cut(bytes.NewReader(data), &buf, opts); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Fatal("empty output")
	}
}

func BenchmarkCut(b *testing.B) {
	data := loadBigFile(b)
	opts := Options{Fields: []int{1, 3, 4, 5}, Delimiter: "\t"}

	for b.Loop() {
		if err := Cut(bytes.NewReader(data), io.Discard, opts); err != nil {
			b.Fatal(err)
		}
	}
}
