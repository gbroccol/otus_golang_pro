package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOpeningSourceFile     = errors.New("error opening source file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrOpeningSourceFile
	}
	defer fromFile.Close()

	info, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat source file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	fileSize := info.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = fromFile.Seek(offset, io.SeekStart) // сместиться до начала чтения
	if err != nil {
		return fmt.Errorf("cannot seek in source file: %w", err)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer toFile.Close()

	bytesToCopy := getBytesToCopy(offset, limit, fileSize)
	err = copyInternal(bytesToCopy, *fromFile, *toFile)
	if err != nil {
		return fmt.Errorf("error copping file: %w", err)
	}

	fmt.Println("\nCopy complete.")
	return nil
}

func getBytesToCopy(offset, limit, fileSize int64) int64 {
	var bytesToCopy int64
	if limit == 0 || limit > fileSize-offset {
		bytesToCopy = fileSize - offset
	} else {
		bytesToCopy = limit
	}
	return bytesToCopy
}

func copyInternal(bytesToCopy int64, fromFile, toFile os.File) error {
	bufSize := 32 * 1024 // 32 KB — стандартный буфер
	buf := make([]byte, bufSize)
	var copied int64

	for copied < bytesToCopy {
		remaining := bytesToCopy - copied
		if int64(len(buf)) > remaining {
			buf = buf[:remaining]
		}
		n, readErr := fromFile.Read(buf) // n - кол-во прочитанных символов
		if n > 0 {
			written, writeErr := toFile.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("write error: %w", writeErr)
			}
			if written != n {
				return fmt.Errorf("incomplete write: written %d, expected %d", written, n)
			}
			copied += int64(written)
			fmt.Printf("\rCopying... %d%%", copied*100/bytesToCopy)
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return fmt.Errorf("read error: %w", readErr)
		}
	}
	return nil
}
