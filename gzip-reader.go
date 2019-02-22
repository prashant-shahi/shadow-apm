package main
/*
import (
	"fmt"
	"io"
	"log"
	"os"
	"compress/gzip"
)

func gzipReader (compressedData io.Reader) {
	var reader io.ReadCloser
	log.Output(0, "Function: gzipReader")
	zr, err := gzip.NewReader(&compressedData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("zr: %s\n", String(zr))
	if x, err := io.Copy(os.Stdout, zr); err != nil {
		log.Fatal(err)
	}
	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}
	return 
}
/*
var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	// Setting the Header fields is optional.
	zw.Name = "a-new-hope.txt"
	zw.Comment = "an epic space opera by George Lucas"
	zw.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)

	_, err := zw.Write([]byte("A long time ago in a galaxy far, far away..."))
	if err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}*/