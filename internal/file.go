package internal

import "fmt"

func newCachedFile(path string, data []byte, contentType string) *cachedFile {
	return &cachedFile{
		path:          path,
		data:          data,
		contentType:   contentType,
		contentLength: fmt.Sprintf("%v", len(data)),
	}
}

type cachedFile struct {
	path          string
	data          []byte
	contentLength string
	contentType   string
}

func (f *cachedFile) Path() string {
	return f.path
}
func (f *cachedFile) Data() []byte {
	return f.data
}
func (f *cachedFile) ContentType() string {
	return f.contentType
}
func (f *cachedFile) ContentLength() string {
	return f.contentLength
}
