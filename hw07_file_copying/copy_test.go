package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		from   string
		offset int64
		limit  int64
	}
	tests := []struct {
		testName       string
		args           args
		expectError    error
		expectedResult string
	}{
		{
			testName: "full-copy",
			args: args{
				from:   "testdata/input.txt",
				offset: 0,
				limit:  0,
			},
			expectError:    nil,
			expectedResult: "testdata/out_offset0_limit0.txt",
		},
		{
			testName: "copy-with-offset-0-limit-10",
			args: args{
				from:   "testdata/input.txt",
				offset: 0,
				limit:  10,
			},
			expectError:    nil,
			expectedResult: "testdata/out_offset0_limit10.txt",
		},
		{
			testName: "copy-with-offset-0-limit-1000",
			args: args{
				from:   "testdata/input.txt",
				offset: 0,
				limit:  1000,
			},
			expectError:    nil,
			expectedResult: "testdata/out_offset0_limit1000.txt",
		},
		{
			testName: "copy-with-offset-0-limit-10000",
			args: args{
				from:   "testdata/input.txt",
				offset: 0,
				limit:  10000,
			},
			expectError:    nil,
			expectedResult: "testdata/out_offset0_limit10000.txt",
		},
		{
			testName: "copy-with-offset-0-limit-10",
			args: args{
				from:   "testdata/input.txt",
				offset: 100,
				limit:  1000,
			},
			expectError:    nil,
			expectedResult: "testdata/out_offset100_limit1000.txt",
		},
		{
			testName: "copy-with-offset-0-limit-10",
			args: args{
				from:   "testdata/input.txt",
				offset: 6000,
				limit:  1000,
			},
			expectError:    nil,
			expectedResult: "testdata/out_offset6000_limit1000.txt",
		},
		{
			testName: "offset exceeds file",
			args: args{
				from:   "testdata/input.txt",
				offset: 9999,
				limit:  0,
			},
			expectError: ErrOffsetExceedsFileSize,
		},
		{
			testName: "unsupported file (directory)",
			args: args{
				from:   "testdata",
				offset: 0,
				limit:  0,
			},
			expectError: ErrUnsupportedFile,
		},
		{
			testName: "file not exists",
			args: args{
				from:   "file-not-exists",
				offset: 0,
				limit:  0,
			},
			expectError: ErrOpeningSourceFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tmpFile := filepath.Join(os.TempDir(), "test-"+tt.testName+".out")
			defer os.Remove(tmpFile)

			err := Copy(tt.args.from, tmpFile, tt.args.offset, tt.args.limit)

			if !errors.Is(err, tt.expectError) {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			if tt.expectError == nil && len(tt.expectedResult) > 0 {
				isEqual, err := filesAreEqual(tt.expectedResult, tmpFile)
				if err != nil {
					t.Fatalf("failed to read output file: %v", err)
				}
				if !isEqual {
					t.Errorf("%q - unexpected content", tt.testName)
				}
			}
		})
	}
}

func filesAreEqual(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	const chunkSize = 4096
	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)

	for {
		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		if n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}

		if err1 == io.EOF && err2 == io.EOF {
			return true, nil // оба файла закончились, и они совпадают
		}

		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}
	}
}
