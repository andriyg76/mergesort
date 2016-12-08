package mergesort
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/andriyg76/glogger"
	"io"
)

type arrayReader struct {
	position int
	slice    []byte
}

var log = glogger.Create(glogger.TRACE)

func (i *arrayReader) Close() error {
	i.slice = nil
	return nil
}

func (i *arrayReader) Read(target []byte) (int, error) {
	if i.slice == nil {
		return 0, io.EOF
	}
	read := copy(target, i.slice[i.position:])
	i.position += read
	log.Debug("Read %s bytes %s", read, target[:read])
	if i.position >= len(i.slice) {
		i.slice = nil
	}
	return read, nil
}

func TestAsyncFileReader(t *testing.T) {
	err, _ := NewAsyncFileReader(nil, log.TraceLogger())

	t.Log(err)

	assert.NotNil(t, err)

	err, reader := NewAsyncFileReader(&arrayReader{
		slice: []byte("string\r\nst\rring2\r\n"),
	}, log.TraceLogger())

	assert.Nil(t, err)
	assert.NotNil(t, reader)

	err, r := reader.ReadLine()
	t.Log("1-th read: ", r, " err: ", err)
	assert.Nil(t, err)
	assert.Equal(t, "string", r)

	err, r = reader.ReadLine()
	t.Log("2-th read: ", r, " err: ", err)
	assert.Nil(t, err)
	assert.Equal(t, "st\rring2", r)

	err, r = reader.ReadLine()
	t.Log("3-th read: ", r, " err: ", err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, "", r)
}

func TestCombineReaders(t *testing.T) {
	err, reader := NewAsyncFileReader(&arrayReader{
		slice: []byte("\n67\n8\n99\n"),
	}, log.TraceLogger())

	assert.Nil(t, err)
	assert.NotNil(t, reader)

	err2, reader2 := NewAsyncFileReader(&arrayReader{
		slice: []byte("7\n9\n"),
	}, log.TraceLogger())

	assert.Nil(t, err2)
	assert.NotNil(t, reader2)

	r := MergeTwoReaders(reader, reader2, AbcStrLess)

	var res []string

	for {
		log.Debug("MergeReaders state %s", r)
		e, s := r.ReadLine()
		if e != nil && e != io.EOF {
			assert.Fail(t, "error ", e)
			break
		}
		if e == io.EOF {
			break
		}
		res = append(res, s)
	}

	assert.Equal(t, []string{"", "67", "7", "8", "9", "99"}, res)
}

func TestCombineReaderWithEmpty(t *testing.T) {
	err, reader := NewAsyncFileReader(&arrayReader{
		slice: []byte("\n8\n67\n99\n"),
	}, log.TraceLogger())

	assert.Nil(t, err)
	assert.NotNil(t, reader)

	err2, reader2 := NewAsyncFileReader(&arrayReader{
		slice: []byte(""),
	}, log.TraceLogger())

	assert.Nil(t, err2)
	assert.NotNil(t, reader2)

	r := MergeTwoReaders(reader, reader2, AbcStrLess)

	var res []string

	for {
		e, s := r.ReadLine()
		if e != nil && e != io.EOF {
			assert.Fail(t, "error ", e)
			break
		}
		if e == io.EOF {
			break
		}
		res = append(res, s)
	}

	assert.Equal(t, []string{"", "8", "67", "99"}, res)
}

func TestMergeSort(t *testing.T) {
	err, reader := NewAsyncFileReader(&arrayReader{
		slice: []byte("67\n8\n99"),
	}, log.TraceLogger())

	assert.Nil(t, err)
	assert.NotNil(t, reader)

	err2, reader2 := NewAsyncFileReader(&arrayReader{
		slice: []byte("7\n9"),
	}, log.TraceLogger())

	assert.Nil(t, err2)
	assert.NotNil(t, reader2)

	err3, reader3 := NewAsyncFileReader(&arrayReader{
		slice: []byte("77\n88"),
	}, log.TraceLogger())

	assert.Nil(t, err3)
	assert.NotNil(t, reader3)

	r := MergeSort(AbcStrLess, reader, reader2, reader3)

	var res []string

	for {
		log.Debug("MergeReaders state %s", r)
		e, s := r.ReadLine()
		if e != nil && e != io.EOF {
			assert.Fail(t, "error ", e)
			break
		}
		res = append(res, s)
		if e == io.EOF {
			break
		}
	}

	assert.Equal(t, []string{"67", "7", "77", "8", "88", "9", "99", ""}, res)
}