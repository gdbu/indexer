package indexer

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/edsrzf/mmap-go"
	"github.com/hatchify/atoms"
	"github.com/hatchify/errors"
)

// New will return a new Indexer
func New(filename string) (ip *Indexer, err error) {
	var i Indexer
	if i.f, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0744); err != nil {
		err = fmt.Errorf("error opening file \"%s\": %v", filename, err)
		return
	}
	defer func() {
		if err != nil {
			i.f.Close()
		}
	}()

	if err = i.setSize(); err != nil {
		err = fmt.Errorf("error setting size: %v", err)
		return
	}

	if i.mm, err = mmap.Map(i.f, os.O_RDWR, 0); err != nil {
		err = fmt.Errorf("error initializing MMAP: %v", err)
		return
	}

	// Set underlying index bytes as MMAP bytes
	i.index = (*atoms.Uint64)(unsafe.Pointer(&i.mm[0]))

	// Associate returning pointer to created Indexer
	ip = &i
	return
}

// Indexer manages indexes
type Indexer struct {
	f     *os.File
	mm    mmap.MMap
	index *atoms.Uint64

	closed atoms.Bool
}

// Get will get the current Indexer value
func (i *Indexer) Get() (value uint64) {
	return i.index.Load()
}

// Next will increment the Indexer value
func (i *Indexer) Next() (next uint64) {
	return i.index.Add(1) - 1
}

// Set will set the current Indexer value
func (i *Indexer) Set(value uint64) {
	i.index.Store(value)
}

// Flush will force a flush
// Note: The OS handles this automatically for MMAP data. This isn't necessary for most use-cases.
// This can be used to for situations where ACID compliance needs to be 100% guaranteed
func (i *Indexer) Flush(value uint64) {
	i.index.Store(value)
}

// Close will close an Indexer
func (i *Indexer) Close() (err error) {
	if !i.closed.Set(true) {
		return errors.ErrIsClosed
	}

	var errs errors.ErrorList
	errs.Push(i.f.Close())
	errs.Push(i.mm.Flush())
	errs.Push(i.mm.Unmap())
	return errs.Err()
}

func (i *Indexer) setSize() (err error) {
	var fi os.FileInfo
	if fi, err = i.f.Stat(); err != nil {
		err = fmt.Errorf("error getting file information: %v", err)
		return
	}

	if fi.Size() == 64 {
		return
	}

	return i.f.Truncate(64)
}
