package handlers

import (
	"bytes"
	"io"
)

func createBody(jsonStr string) io.Reader {
	return bytes.NewBufferString(jsonStr)
}
