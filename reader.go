/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2014 Yani Iliev <yani@iliev.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package wpress

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

// Reader structure
type Reader struct {
	Filename      string
	File          *os.File
	NumberOfFiles int
}

// NewReader creates a new Reader instance and calls its constructor
func NewReader(filename string) (*Reader, error) {
	// create a new instance of Reader
	r := &Reader{filename, nil, 0}

	// call the constructor
	err := r.Init()
	if err != nil {
		return nil, err
	}

	// return Reader instance
	return r, nil
}

// Init is the constructor of Reader struct
func (r *Reader) Init() error {
	// try to open the file
	file, err := os.Open(r.Filename)
	if err != nil {
		return err
	}

	// file was openned, assign the handle to the holding variable
	r.File = file

	return nil
}

// ExtractFile extracts file that matches tha filename and path from archive
func (r Reader) ExtractFile(filename string, path string) ([]byte, error) {
	// TODO: implement
	return nil, nil
}

// Extract all files from archive
func (r Reader) Extract() (int, error) {
	outputPath := "."
	NumberOfFiles, err := r.ExtractToPath(outputPath)
	return NumberOfFiles, err
}

// Extract all files from archive
func (r Reader) ExtractToPath(outputPath string) (int, error) {
	// put pointer at the beginning of the file
	r.File.Seek(0, 0)

	// loop until end of file was reached
	for {
		// read header block
		block, err := r.GetHeaderBlock()
		if err != nil {
			return 0, err
		}

		// initialize new header
		h := &Header{}

		// check if block equals EOF sequence
		if bytes.Compare(block, h.GetEOFBlock()) == 0 {
			// EOF reached, stop the loop
			break
		}

		// populate header from our block bytes
		h.PopulateFromBytes(block)

		pathToFile := path.Clean(outputPath + string(os.PathSeparator) + string(bytes.Trim(h.Prefix, "\x00")) + string(os.PathSeparator) + string(bytes.Trim(h.Name, "\x00")))

		err = os.MkdirAll(path.Dir(pathToFile), 0755)
		if err != nil {
			fmt.Println(err)
			return r.NumberOfFiles, err
		}

		// try to open the file
		file, err := os.Create(pathToFile)
		if err != nil {
			return r.NumberOfFiles, err
		}

		totalBytesToRead, _ := h.GetSize()
		for {
			bytesToRead := 512
			if bytesToRead > totalBytesToRead {
				bytesToRead = totalBytesToRead
			}

			if bytesToRead == 0 {
				break
			}

			content := make([]byte, bytesToRead)
			bytesRead, err := r.File.Read(content)
			if err != nil {
				return r.NumberOfFiles, err
			}

			totalBytesToRead -= bytesRead
			contentRead := content[0:bytesRead]

			_, err = file.Write(contentRead)
			if err != nil {
				return r.NumberOfFiles, err
			}
		}

		file.Close()

		// increment file counter
		r.NumberOfFiles++
	}

	return r.NumberOfFiles, nil

}

// GetHeaderBlock reads and returns header block from archive
func (r Reader) GetHeaderBlock() ([]byte, error) {
	// create buffer to keep the header block
	block := make([]byte, headerSize)

	// read the header block
	bytesRead, err := r.File.Read(block)
	if err != nil {
		return nil, err
	}

	if bytesRead != headerSize {
		return nil, errors.New("unable to read header block size")
	}

	return block, nil
}

// GetFilesCount returns the number of files in archive
func (r Reader) GetFilesCount() (int, error) {
	// test if we have enumerated the archive already
	if r.NumberOfFiles != 0 {
		return r.NumberOfFiles, nil
	}

	// put pointer at the beginning of the file
	r.File.Seek(0, 0)

	// loop until end of file was reached
	for {
		// read header block
		block, err := r.GetHeaderBlock()
		if err != nil {
			return 0, err
		}

		// initialize new header
		h := &Header{}

		// check if block equals EOF sequence
		if bytes.Compare(block, h.GetEOFBlock()) == 0 {
			// EOF reached, stop the loop
			break
		}

		// populate header from our block bytes
		h.PopulateFromBytes(block)

		// set pointer after file content, to the next header block
		size, err := h.GetSize()
		if err != nil {
			return 0, err
		}
		r.File.Seek(int64(size), 1)

		// increment file counter
		r.NumberOfFiles++
	}

	return r.NumberOfFiles, nil
}

// Header and other necessary imports and structs should be defined above this.
// Added by Slavi Marinov so no need to extract to view files.
// List lists all files in the archive without extracting them.
func (r *Reader) List() ([]string, error) {
	var fileList []string

	// Reset the file counter as we'll be re-iterating the archive.
	r.NumberOfFiles = 0

	// Ensure we start from the beginning of the file.
	_, err := r.File.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	for {
		// Read the header block.
		block, err := r.GetHeaderBlock()
		if err != nil {
			// If an error occurs (e.g., EOF), break the loop.
			break
		}

		// Initialize a new header to hold the data.
		h := &Header{}

		// Check if the block is an EOF marker.
		if bytes.Compare(block, h.GetEOFBlock()) == 0 {
			break
		}

		// Populate the header with data from the block.
		h.PopulateFromBytes(block)

		// Step 1 & 2: Convert the string to an integer
		timestampStr := string(bytes.Trim(h.Mtime, "\x00"))
		unixTimestamp, errTs := strconv.ParseInt(timestampStr, 10, 64)
		formattedDate := timestampStr // defaults to timestamp if conversion fails

		if errTs == nil {
			// Step 3: Convert integer to time.Time object
			t := time.Unix(unixTimestamp, 0)

			// Step 4: Format the time.Time object to "YYYY-MM-DD HH:MM:SS"
			formattedDate = t.Format("2006-01-02 15:04:05")
		}

		// Create a line SIZE Mtime path
		filePath := string(bytes.Trim(h.Size, "\x00")) + " " + formattedDate + " " + path.Clean("."+string(os.PathSeparator)+string(bytes.Trim(h.Prefix, "\x00"))+string(os.PathSeparator)+string(bytes.Trim(h.Name, "\x00")))

		// Add the file path to the list of files.
		fileList = append(fileList, filePath)

		// Calculate the size of the content and skip over it to the next header.
		size, _ := h.GetSize()
		_, err = r.File.Seek(int64(size), 1)
		if err != nil {
			return fileList, err
		}

		// Increment the file counter.
		r.NumberOfFiles++
	}

	return fileList, nil
}
