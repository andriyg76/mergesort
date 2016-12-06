package mergesort

import (
	"io"
	log "github.com/andriyg76/glogger"
)

type StrLess func(a, b string) bool

func AbcStrLess(a, b string) bool {
	return a < b
}

type eofReader struct {
}

func (i *eofReader) ReadLine() (error, string) {
	return io.EOF, ""
}

var eof = &eofReader{}

type combinedReaders struct {
	left, right Reader
	one, two    string
	err1, err2  error
	cmp         StrLess
}

func (i *combinedReaders) ReadLine() (error, string) {
	log.Trace("State: %s %s errors: %v %v", i.one, i.two, i.err1, i.err2)
	if i.err1 != nil && i.err1 != io.EOF {
		return i.err1, ""
	} else if i.err2 != nil && i.err2 != io.EOF {
		return i.err2, ""
	}

	var value string
	var err error
	if i.err1 == io.EOF && i.err2 == io.EOF {
		err = io.EOF
	} else if i.cmp(i.one, i.two) || i.err2 == io.EOF {
		err, value = i.err1, i.one
		i.err1, i.one = i.left.ReadLine()
	} else if !i.cmp(i.one, i.two) || i.err1 == io.EOF {
		err, value = i.err2, i.two
		i.err2, i.two = i.right.ReadLine()
	} else {
		log.Panic("%v %s %v %s", i.err1, i.one, i.err2, i.right)
	}
	log.Trace("Returning value: %s error: %v", value, err)
	return err, value
}

func MergeTwoReaders(left, right Reader, cmp StrLess) Reader {
	i := &combinedReaders{
		left: left,
		right: right,
		cmp: cmp,
	}
	i.err1, i.one = i.left.ReadLine()
	i.err2, i.two = i.right.ReadLine()
	return i
}

func MergeSort(cmp StrLess, readers... Reader) Reader {
	if len(readers) == 0 {
		return eof
	} else if len(readers) == 1 {
		return readers[0]
	}
	middle := len(readers) / 2
	return MergeTwoReaders(MergeSort(cmp, readers[:middle]...),
		MergeSort(cmp, readers[middle:]...), cmp)
}