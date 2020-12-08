package indexer

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var (
	uintSink       uint64
	exampleIndexer *Indexer
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

func BenchmarkIndexer_Next(b *testing.B) {
	var (
		indexer *Indexer
		err     error
	)

	if indexer, err = New("testfile.idb"); err != nil {
		b.Fatal(err)
	}
	defer os.Remove("testfile.idb")
	defer indexer.Close()

	// Now that we've initialized, reset timer
	b.ResetTimer()

	for i := uint64(0); i < uint64(b.N); i++ {
		uintSink = indexer.Next()
	}
}

/*
// If you want to benchmark against dbutils and/or with bolt.DB, please uncomment this section.
// I comment this out so that the go.mod file doesn't include unnecessary imports
func BenchmarkIndexer_Next_inside_bolt_txn(b *testing.B) {
	var (
		indexer *Indexer
		db      *bolt.DB
		err     error
	)

	if indexer, err = New("testfile.idb"); err != nil {
		b.Fatal(err)
	}
	defer os.Remove("testfile.idb")
	defer indexer.Close()

	if db, err = bolt.Open("./test.bdb", 0744, nil); err != nil {
		b.Fatal(err)
	}
	defer os.Remove("./test.bdb")
	defer db.Close()

	// Now that we've initialized, reset timer
	b.ResetTimer()

	if err = db.Update(func(txn *bolt.Tx) (err error) {
		for i := uint64(0); i < uint64(b.N); i++ {
			uintSink = indexer.Next()
		}

		return
	}); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkDBUtils_Next(b *testing.B) {
	var (
		db  *bolt.DB
		err error
	)

	if db, err = bolt.Open("./test.bdb", 0744, nil); err != nil {
		b.Fatal(err)
	}
	defer os.Remove("./test.bdb")
	defer db.Close()

	dbu := dbutils.New(8)
	if err = db.Update(func(txn *bolt.Tx) (err error) {
		if err = dbu.Init(txn); err != nil {
			b.Fatal(err)
		}

		return
	}); err != nil {
		b.Fatal(err)
	}

	// Now that we've initialized, reset timer
	b.ResetTimer()

	if err = db.Update(func(txn *bolt.Tx) (err error) {
		for i := uint64(0); i < uint64(b.N); i++ {
			if uintSink, err = dbu.NextIndex(txn, []byte("index")); err != nil {
				b.Fatal(err)
			}
		}

		return
	}); err != nil {
		b.Fatal(err)
	}
}
*/

func ExampleNew() {
	var err error
	if exampleIndexer, err = New("./indexer.idb"); err != nil {
		log.Fatal(err)
	}
}

func ExampleIndexer_Next() {
	value := exampleIndexer.Next()
	fmt.Println("Indexer next value is", value)
}

func ExampleIndexer_Set() {
	exampleIndexer.Set(1337)
}

func ExampleIndexer_Close() {
	if err := exampleIndexer.Close(); err != nil {
		log.Fatal(err)
	}
}
