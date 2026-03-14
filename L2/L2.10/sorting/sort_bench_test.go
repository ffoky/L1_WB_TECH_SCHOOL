package sorting_test

import (
	"bytes"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"L2.10/sorting"
)

func loadBigTxt(b *testing.B) []byte {
	b.Helper()
	data, err := os.ReadFile("testdata/big.txt")
	require.NoError(b, err)
	return data
}

func BenchmarkSortSingle(b *testing.B) {
	data := loadBigTxt(b)

	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)

	for b.Loop() {
		var buf bytes.Buffer
		err := sorting.Sort(bytes.NewReader(data), &buf, sorting.Options{Numeric: true})
		require.NoError(b, err)
	}
}

func BenchmarkSortParallel(b *testing.B) {
	data := loadBigTxt(b)

	for b.Loop() {
		var buf bytes.Buffer
		err := sorting.Sort(bytes.NewReader(data), &buf, sorting.Options{Numeric: true})
		require.NoError(b, err)
	}
}
