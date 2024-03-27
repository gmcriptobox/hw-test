package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrCreate                = errors.New("unable to create file")
	ErrOpen                  = errors.New("unable to open file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	file, err := os.Create(toPath)
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	if err != nil {
		return ErrCreate
	}

	fileRead, err := os.Open(fromPath)
	defer func() {
		err := fileRead.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return ErrOpen
	}

	fi, err := os.Stat(fromPath)
	if err != nil {
		return ErrOpen
	}
	if fi.Size() == 0 {
		return ErrUnsupportedFile
	}
	if fi.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	fileRead.Seek(offset, io.SeekStart)
	data := make([]byte, 64)
	var count int64
	maxWriteLimit := limit
	if maxWriteLimit > fi.Size() || limit == 0 {
		maxWriteLimit = fi.Size()
	}
	bar := pb.Start64(maxWriteLimit)
	for {
		n, err := fileRead.Read(data)
		if count+int64(n) <= maxWriteLimit {
			count += int64(n)
			file.Write(data[:n])
			bar.Add(n)
		} else {
			file.Write(data[:maxWriteLimit-count])
			bar.Add(int(maxWriteLimit - count))
			break
		}
		if err == io.EOF { // если конец файла
			break // выходим из цикла
		}
	}
	bar.Finish()

	return nil
}
