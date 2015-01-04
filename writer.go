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
	"io/ioutil"
	"os"
)

type Writer struct {
	Filename   string
	File       *os.File
	FilesAdded int
}

/**
 * Create new Writer instance
 */
func NewWriter(filename string) (*Writer, error) {
	// create a new instance of Writer
	w := &Writer{filename, nil, 0}

	// call the constructor
	err := w.Init()
	if err != nil {
		return nil, err
	}

	// return Writer instance
	return w, nil
}

/**
 * Writer constructor
 */
func (w *Writer) Init() error {
	// try to create the file
	file, err := os.Create(w.Filename)
	if err != nil {
		return err
	}

	// file was created, assign it to its holding variable
	w.File = file

	return nil
}

/**
 * Add a file to the archive
 */
func (w *Writer) AddFile(filename string) error {
	// populate header block from the filename passed
	h := &Header{}
	err := h.PopulateFromFilename(filename)
	if err != nil {
		return err
	}

	// write header block
	_, err = w.File.Write(h.GetHeaderBlock())
	if err != nil {
		return err
	}
	// write file content
	// open the file for reading
	input, err := os.Open(filename)
	if err != nil {
		return err
	}

	for {
		bytesToRead := 512
		content := make([]byte, bytesToRead)
		bytesRead, err := input.Read(content)
		if err != nil {
			return err
		}

		// if we have read less than 100 or 0 bytes, we reached end of file
		if bytesRead < bytesToRead {
			// obtain only the bytes that were read
			contentRead := content[0:bytesRead]
			_, err = w.File.Write(contentRead)
			if err != nil {
				return err
			}

			// exit the loop, we reached end of file
			break
		}

		// we write the content we just read to the archive
		_, err = w.File.Write(content)
		if err != nil {
			return err
		}
	}

	// done reading from the file, let's close it
	err = input.Close()
	if err != nil {
		return err
	}

	// file was added to the archive, increment fileAdded
	w.FilesAdded += 1

	return nil
}

/**
 * Add a directory to the archive
 */
func (w *Writer) AddDirectory(path string) error {
	fiArray, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// go over every directory entry and add it
	// files are added using AddFile, directories are parsed recursevely
	for _, fi := range fiArray {
		if fi.IsDir() {
			w.AddDirectory(path + string(os.PathSeparator) + fi.Name())
		} else {
			err = w.AddFile(path + string(os.PathSeparator) + fi.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/**
 * Closes the archive, appends EOF sequence to the end of the file
 */
func (w Writer) Close() error {
	// if we haven't added any files, we don't append EOF sequence
	if w.FilesAdded == 0 {
		return nil
	}
	// create new header instance
	h := &Header{}

	// write eof sequence
	_, err := w.File.Write(h.GetEofBlock())
	if err != nil {
		return err
	}

	// close the archive
	err = w.File.Close()
	if err != nil {
		return err
	}

	return nil
}
