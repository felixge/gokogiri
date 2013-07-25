package xml

/*
#cgo pkg-config: libxml-2.0
#include "helper.h"
*/
import "C"
import "fmt"
import "io"
import "reflect"

import (
	. "github.com/moovweb/gokogiri/util"
	"unsafe"
)

func NewReader(r io.Reader, inEncoding, url []byte, options int) (*Reader, error) {
	inEncoding = AppendCStringTerminator(inEncoding)

	var urlPtr, encodingPtr, goReaderPtr unsafe.Pointer
	if len(url) > 0 {
		url = AppendCStringTerminator(url)
		urlPtr = unsafe.Pointer(&url[0])
	}
	if len(inEncoding) > 0 {
		encodingPtr = unsafe.Pointer(&inEncoding[0])
	}

	reader := &Reader{r: r}
	goReaderPtr = unsafe.Pointer(reader)
	reader.ptr = C.xmlNewReader(goReaderPtr, urlPtr, encodingPtr, C.int(options), nil, 0)

	if reader.ptr == nil {
		return nil, ERR_FAILED_TO_PARSE_XML
	} else {
		return reader, nil
	}
}

//export readerReadCallback
func readerReadCallback(context unsafe.Pointer, buffer *C.char, length int) int {
	r := (*Reader)(context)
	h := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(buffer)),
		Len:  length,
		Cap:  length,
	}

	buf := *(*[]byte)(unsafe.Pointer(h))
	n, err := r.r.Read(buf)
	if err == io.EOF {
		return 0
	} else if err != nil {
		// @TODO: Can we also return the error message somehow?
		return -1
	}
	return n
}

type Reader struct {
	ptr C.xmlTextReaderPtr
	r   io.Reader
}

func (r *Reader) Read() error {
	ret := C.xmlTextReaderRead(r.ptr)
	if ret == 0 {
		return io.EOF
	} else if ret == 1 {
		return nil
	} else {
		return fmt.Errorf("C.xmlTextReaderRead() returned: %d", ret)
	}
}

func (r *Reader) Name() string {
	p := (*C.char)(unsafe.Pointer(C.xmlTextReaderName(r.ptr)))
	defer C.xmlFreeChars(p)
	name := C.GoString(p)
	return name
}
