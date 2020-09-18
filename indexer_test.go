package indexer

import (
	"os"
	"testing"
)

func TestIndexer_Next(t *testing.T) {
	var (
		indexer *Indexer
		err     error
	)

	if indexer, err = New("testfile.idb"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("testfile.idb")
	defer indexer.Close()

	for i := uint64(0); i < 1000; i++ {
		value := indexer.Next()
		if value != i {
			t.Fatalf("invalid index, expected %d and received %d", i, value)
		}
	}
}
