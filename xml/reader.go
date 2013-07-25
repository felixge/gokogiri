package xml

/*
#cgo pkg-config: libxml-2.0
#include "helper.h"
*/
import "C"
import "fmt"
import "io"

import(
	. "github.com/moovweb/gokogiri/util"
	"unsafe"
)

func NewReader(content, inEncoding, url []byte, options int) (*Reader, error) {
	inEncoding = AppendCStringTerminator(inEncoding)

	var readerPtr C.xmlTextReaderPtr
	contentLen := len(content)

	if contentLen > 0 {
		var contentPtr, urlPtr, encodingPtr unsafe.Pointer
		contentPtr = unsafe.Pointer(&content[0])

		if len(url) > 0 {
			url = AppendCStringTerminator(url)
			urlPtr = unsafe.Pointer(&url[0])
		}
		if len(inEncoding) > 0 {
			encodingPtr = unsafe.Pointer(&inEncoding[0])
		}

		readerPtr = C.xmlNewReader(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

		if readerPtr == nil {
			return nil, ERR_FAILED_TO_PARSE_XML
		} else {
			return &Reader{Ptr: readerPtr}, nil
		}
	} else {
		panic("not done")
	}
}

type Reader struct{
	Ptr C.xmlTextReaderPtr
}

func (r *Reader) Read() (error) {
	ret := C.xmlTextReaderRead(r.Ptr)
	if ret == 0 {
		return io.EOF
	} else if ret == 1 {
		return nil
	} else {
		return fmt.Errorf("C.xmlTextReaderRead() returned: %d", ret)
	}
}

func (r *Reader) Name() string {
	p := (*C.char)(unsafe.Pointer(C.xmlTextReaderName(r.Ptr)))
	defer C.xmlFreeChars(p)
	name := C.GoString(p)
	return name
}
