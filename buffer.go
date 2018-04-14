package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
)

// buffer ...
type buffer struct {
	head     int
	tail     int
	buf      []byte
	spills   int
	spillDir string
}

// Len ...
func (b *buffer) Len() int {
	return b.head / 8
}

// Swap ...
func (b *buffer) Swap(i, j int) {
	swap(b.buf, i, j)
}

// Less ...
func (b *buffer) Less(i, j int) bool {
	return compare(b.buf, i, j) < 0
}

// add ...
func (b *buffer) add(record []byte) error {
	recordSize := 2*8 + len(record)
	if b.free() < recordSize {
		if len(b.buf) < recordSize {
			return fmt.Errorf(
				"record is too large to fit in memory - required: %db but "+
					"buffer memory can only hold %db",
				recordSize,
				len(b.buf),
			)
		}
		if err := b.spill(); err != nil {
			return err
		}
	}
	b.appendRecord(record)
	return nil
}

// sort ...
func (b *buffer) sort() {
	sort.Sort(b)
}

// spill ...
func (b *buffer) spill() error {
	defer func() {
		b.head = 0
		b.tail = len(b.buf)
		b.spills++
	}()
	b.sort()
	if err := os.MkdirAll(b.spillDir, 0700); err != nil {
		return err
	}
	filename := path.Join(b.spillDir, fmt.Sprintf("spill-%d", b.spills))
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	wb := bufio.NewWriter(w)
	for i := 0; i < b.Len(); i++ {
		p := readInt(b.buf, i*8)
		s := readInt(b.buf, p)
		if _, err := wb.Write(b.buf[p : p+8+s]); err != nil {
			return err
		}
	}
	if err := wb.Flush(); err != nil {
		return err
	}
	return w.Close()
}

// extSort ...
func (b *buffer) externalSort() error {
	// During the final merge phase we will have at most mappers*reducers open files
	// so use this here as well. With a hard minimum of 16 for any situation where we
	// have < 16 mappers.
	ways := mappers
	if ways < 16 {
		ways = 16
	}
	buf := make([]byte, 8)
	for b.spills > 1 {
		newSpills := 0
		for i := 0; i <= b.spills/ways; i++ {
			start := i * ways
			end := start + ways
			if end >= b.spills {
				end = b.spills
			}
			if end-start == 0 {
				continue
			}
			newSpills++
			scanners := make([]recordScanner, end-start)
			for j := 0; j < end-start; j++ {
				filename := path.Join(b.spillDir, fmt.Sprintf("spill-%d", j+start))
				scanners[j] = newFileScanner(filename)
			}
			m, err := newMerger(scanners)
			if err != nil {
				return err
			}
			mergeFilename := path.Join(b.spillDir, "merge")
			f, err := os.Create(mergeFilename)
			if err != nil {
				return err
			}
			w := bufio.NewWriter(f)
			for m.next() {
				r := m.record()
				n := len(r)
				writeInt(buf, 0, n)
				if _, err := w.Write(buf); err != nil {
					return err
				}
				if _, err := w.Write(r); err != nil {
					return err
				}
			}
			if err := m.err(); err != nil {
				return err
			}
			if err := w.Flush(); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
			for j := 0; j < end-start; j++ {
				filename := path.Join(b.spillDir, fmt.Sprintf("spill-%d", j+start))
				if err := os.Remove(filename); err != nil {
					return err
				}
			}
			filename := path.Join(b.spillDir, fmt.Sprintf("spill-%d", i))
			if err := os.Rename(mergeFilename, filename); err != nil {
				return err
			}
		}
		b.spills = newSpills
	}
	return nil
}

// free ...
func (b *buffer) free() int {
	return b.tail - b.head
}

func (b *buffer) appendRecord(record []byte) {
	b.tail -= len(record)
	copy(b.buf[b.tail:b.tail+len(record)], record)
	b.tail -= 8
	writeInt(b.buf, b.tail, len(record))
	writeInt(b.buf, b.head, b.tail)
	b.head += 8
}

// newBuffer ...
func newBuffer(bufMem int, spillDir string) *buffer {
	return &buffer{
		head:     0,
		tail:     bufMem,
		buf:      make([]byte, bufMem),
		spills:   0,
		spillDir: spillDir,
	}
}
