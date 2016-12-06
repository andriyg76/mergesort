package mergesort

import (
	"errors"
	"bufio"
	"strings"
	"io"
	"log"
)

type Reader interface {
	ReadLine() (error, string)
}

// Composition of reader and close, os.File have to match
type DisposableIoReader interface {
	io.Reader
	io.Closer
}

type arrayPos struct {
	array []string
	pos   int
}

func (i *arrayPos) ReadLine() (error, string)  {
	if i.pos >= len(i.array) {
		return io.EOF, ""
	}

	val:= i.array[i.pos]
	i.pos ++
	return nil, val
}

func NewArrayReader(array []string) Reader {
	return &arrayPos{
		array: array,
		pos: 0,
	}
}

type DisposableReader interface {
	Reader
	Close()
}

type stringAndErr struct {
	string string
	error  error
}

type fileReader struct {
	file    DisposableIoReader
	channel chan stringAndErr
	gotEOF  bool
}

func NewAsyncFileReader(file DisposableIoReader, trace* log.Logger) (error, DisposableReader) {
	if file == nil {
		return errors.New("null pointer exception: file"), nil
	}

	fileRrd := bufio.NewReader(file)

	reader := &fileReader{
		file: file,
		channel: make(chan stringAndErr),
	}

	go func() {
		for {
			line, err := fileRrd.ReadString('\n')
			if trace != nil {
				trace.Printf("Read line: %q error: %v", line, err)
			}
			line = strings.TrimRight(line, "\n\r")
			if err == io.EOF && !reader.gotEOF && line != "" {
				reader.gotEOF = true
				err = nil
			}
			reader.channel <- stringAndErr{
				string: line,
				error: err,
			}
			if err == io.EOF {
				reader.Close()
				break
			}
		}
	}();

	return nil, reader
}

func (i *fileReader) Close() {
	if i.file != nil {
		i.file.Close()
	}

	close(i.channel)
}

func (i *fileReader) ReadLine() (error, string) {
	lae, ok := <-i.channel
	if !ok {
		return io.EOF, ""
	}
	return lae.error, lae.string
}
