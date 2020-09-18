# Indexer
Indexer is a persistent indexing library which utilizes MMAP and atomic values for fast and thread-safe incrementing.

## Usage

### New
```go

func ExampleNew() {
	var err error
	if exampleIndexer, err = New("./indexer.idb"); err != nil {
		log.Fatal(err)
	}
}
```

### Indexer.Next
```go
func ExampleIndexer_Next() {
	value := exampleIndexer.Next()
	fmt.Println("Indexer next value is", value)
}
```

### Indexer.Set
```go
func ExampleIndexer_Set() {
	exampleIndexer.Set(1337)
}
```

### Indexer.Close
```go
func ExampleIndexer_Close() {
	if err := exampleIndexer.Close(); err != nil {
		log.Fatal(err)
	}
}
```

## Performance
```
BenchmarkIndexer_Next-4						187052113		6.43 ns/op
BenchmarkIndexer_Next_inside_bolt_txn-4		183006438		6.46 ns/op
BenchmarkDBUtils_Next-4						2148330			485 ns/op
```