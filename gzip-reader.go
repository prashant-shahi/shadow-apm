package main

import (
	"compress/gzip"
	"io"
	"log"
)

func gzipReader(compressedData io.Reader) (io.Reader, error) {
	log.Output(0, "gzip: Reached here")
	reader, err := gzip.NewReader(compressedData)
	if err != nil {
		log.Output(0, "Error: "+err.Error())
		return nil, err
	}
	return reader, nil
}
