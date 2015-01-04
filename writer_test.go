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
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// TestNewWriter tests creating a new writer
func TestNewWriter(t *testing.T) {
	// test creating a writer
	w, err := NewWriter("testing.wpress")
	defer os.Remove("testing.wpress")
	if err != nil {
		t.Errorf("Failed to create a new Writer because %s", err)
	}

	if w.Filename != "testing.wpress" {
		t.Errorf("Failed to set filename to testing.wpress")
	}
}

// TestInit tests Writer constructor
func TestInit(t *testing.T) {
	// test creating a writer
	w, err := NewWriter("testing.wpress")
	defer os.Remove("testing.wpress")
	if err != nil {
		t.Errorf("Failed to create a new Writer because %s", err)
	}

	if reflect.TypeOf(w) != reflect.TypeOf(&Writer{}) {
		t.Errorf("NewWriter does not return instance of Writer")
	}
}

// TestAddFile tests adding af file
func TestAddFile(t *testing.T) {
	// obtain cwd
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current working dir: %s", err)
	}

	// path to our test data
	path := _getPathToTests(t)

	// create a temporary folder for our tests
	tempPath, err := ioutil.TempDir(cwd, "wpressTest")
	if err != nil {
		t.Errorf("Failed to create temporary folder %s", err)
	}

	// make sure to clean after ourselves
	defer os.Remove(tempPath)
	filename := tempPath + string(os.PathSeparator) + "output.wpress"
	defer os.Remove(filename)

	w, err := NewWriter(filename)
	if err != nil {
		t.Errorf("Failed to create a new Writer because %s", err)
	}

	// create an array of the files that we want to add
	filesToAdd := [...]string{"logo.svg", "lipsum.txt", "logo3.png"}
	for _, file := range filesToAdd {
		err = w.AddFile(path + string(os.PathSeparator) + file)
		if err != nil {
			t.Errorf("Failed to add `%s` because: %s", file, err)
		}
	}

	// close the archive
	w.Close()

	// verify that the file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("The archive file was not created: %s", filename)
	}

	// create a new reader instance with the archive we just created
	r, err := NewReader(filename)
	if err != nil {
		t.Errorf("Failed to create a new Reader instace: %s", err)
	}

	// get the files inside the archive
	filesCount, err := r.GetFilesCount()
	if err != nil {
		t.Errorf("Unable to get files count: %s", err)
	}
	// see if all of our files were added
	if len(filesToAdd) != filesCount {
		t.Errorf(
			"The archive contains %d files instead of %d",
			filesCount,
			len(filesToAdd))
	}
}

// TestAddDirectory tests adding a directory
func TestAddDirectory(t *testing.T) {
	// obtain cwd
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current working dir: %s", err)
	}

	// path to our test data
	path := _getPathToTests(t)

	// create a temporary folder for our tests
	tempPath, err := ioutil.TempDir(cwd, "wpressTest")
	if err != nil {
		t.Errorf("Failed to create temporary folder %s", err)
	}

	// make sure to clean after ourselves
	defer os.Remove(tempPath)
	filename := tempPath + string(os.PathSeparator) + "output.wpress"
	defer os.Remove(filename)

	w, err := NewWriter(filename)
	if err != nil {
		t.Errorf("Failed to create a new Writer because %s", err)
	}

	w.AddDirectory(path)

	w.Close()

	// verify that the file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("The archive file was not created: %s", filename)
	}

	// create a new reader instance with the archive we just created
	r, err := NewReader(filename)
	if err != nil {
		t.Errorf("Failed to create a new Reader instace: %s", err)
	}

	// get the files inside the archive
	filesCount, err := r.GetFilesCount()
	if err != nil {
		t.Errorf("Unable to get files count: %s", err)
	}
	// see if all of our files were added
	if 5 != filesCount {
		t.Errorf("The archive contains %d files instead of 5", filesCount)
	}
}

// TestClose tests closing archive
func TestClose(t *testing.T) {
	// obtain cwd
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current working dir: %s", err)
	}

	// path to our test data
	path := _getPathToTests(t)

	// create a temporary folder for our tests
	tempPath, err := ioutil.TempDir(cwd, "wpressTest")
	if err != nil {
		t.Errorf("Failed to create temporary folder %s", err)
	}

	// make sure to clean after ourselves
	defer os.Remove(tempPath)
	filename := tempPath + string(os.PathSeparator) + "output.wpress"
	defer os.Remove(filename)

	w, err := NewWriter(filename)
	if err != nil {
		t.Errorf("Failed to create a new Writer because %s", err)
	}

	w.Close()

	file, err := os.Open(filename)
	data := make([]byte, 100)
	_, err = file.Read(data)
	if err != io.EOF {
		t.Errorf("Close method wrote data to the file when it shouldn't have to.")
	}

	os.Remove(filename)
	w, err = NewWriter(filename)

	err = w.AddFile(path + string(os.PathSeparator) + "lipsum.txt")
	if err != nil {
		t.Errorf("Failed to add a file: %s", err)
	}
	err = w.Close()
	if err != nil {
		t.Errorf("Failed to close the archive: %s", err)
	}

	file, err = os.Open(filename)
	if err != nil {
		t.Errorf("Unable to open archive for reading: %s", err)
	}
	h := &Header{}
	data = make([]byte, headerSize)
	_, err = file.Seek(-headerSize, 2)
	if err != nil {
		t.Errorf("Unable to set the offset: %s", err)
	}
	_, err = file.Read(data)
	if err != nil {
		t.Errorf("Unable to read from the archive: %s", err)
	}

	if bytes.Compare(data, h.GetEOFBlock()) != 0 {
		t.Errorf("EOF sequence was not found at the end of the file")
	}
}
