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

// Composition of reader and close, os.File have to match
type DisposableIoReader interface {
	io.Reader
	io.Closer
}

type arrayPos struct {
	array []string
	pos   int
}

func (i *arrayPos) Close() error {
	i = nil
	return nil
}

func (i *arrayPos) ReadLine() (error, string) {
	if i.pos >= len(i.array) {
		return io.EOF, ""
	}

	val := i.array[i.pos]
	i.pos++
	return nil, val
}

func NewArrayReader(array []string) DisposableReader {
	return &arrayPos{
		array: array,
		pos:   0,
	}
}

type DisposableReader interface {
	Reader
	Close() error
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

func NewAsyncFileReader(file DisposableIoReader, trace *stdlog.Logger) (error, DisposableReader) {
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
				reader.Close()
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

func (i *fileReader) Close() error {
	var errors []error
	if i.file != nil {
		if err := i.file.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	close(i.channel)
	if len(errors) != 0 {
		return MultipleErrors{
			errors: errors,
		}
	}
	return nil
}

func (i *fileReader) ReadLine() (error, string) {
	lae, ok := <-i.channel
	if !ok {
		return io.EOF, ""
	}
	return lae.error, lae.string
}
