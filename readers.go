package mergesort

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	stdlog "log"
	"strings"
)

type Reader interface {
	ReadLine() (error, string)
}

type arrayPos struct {
	array []string
	pos   int
}

func (i *arrayPos) ReadLine() (error, string) {
	if i.pos >= len(i.array) {
		return io.EOF, ""
	}

	val := i.array[i.pos]
	i.pos++
	return nil, val
}

func NewArrayReader(array []string) Reader {
	return &arrayPos{
		array: array,
		pos:   0,
	}
}

type stringAndErr struct {
	string string
	error  error
}

type fileReader struct {
	file    io.Reader
	channel chan stringAndErr
	gotEOF  bool
}

func NewAsyncFileReader(file io.Reader, trace *stdlog.Logger) (error, Reader) {
	if file == nil {
		return errors.New("null pointer exception: file"), nil
	}

	fileRrd := bufio.NewReader(file)

	reader := &fileReader{
		file:    file,
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
				error:  err,
			}
			if err == io.EOF {
				break
			}
		}
	}()

	return nil, reader
}

type MultipleErrors struct {
	errors []error
}

func (i MultipleErrors) Error() string {
	return fmt.Sprintf("Multiple errors: %s", i.errors)
}

func (i *fileReader) ReadLine() (error, string) {
	lae, ok := <-i.channel
	if !ok {
		return io.EOF, ""
	}
	return lae.error, lae.string
}
