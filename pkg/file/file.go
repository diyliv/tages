package file

import "bytes"

type FileBuffer struct {
	Buffer *bytes.Buffer
}

func NewFileBuffer() *FileBuffer {
	return &FileBuffer{
		Buffer: &bytes.Buffer{},
	}
}

func (fb *FileBuffer) Write(chunk []byte) error {
	_, err := fb.Buffer.Write(chunk)
	return err
}
