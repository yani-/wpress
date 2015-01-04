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
	"os"
	"reflect"
	"testing"
)

const (
	TestDataName    = "testdata"
	TestArchiveName = "test_archive.wpress"
)

var pathToTestData string

// _getPathToTests returns absolute path to tests folder
func _getPathToTests(t *testing.T) string {
	if pathToTestData != "" {
		return pathToTestData
	}
	// obtain cwd
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current working dir: %s", err)
	}

	pathToTestData = cwd + string(os.PathSeparator) + TestDataName
	return pathToTestData
}

// TestNewReader tests creating a new reader
func TestNewReader(t *testing.T) {
	path := _getPathToTests(t)
	filename := path + string(os.PathSeparator) + TestArchiveName
	r, err := NewReader(filename)
	if err != nil {
		t.Errorf("Unable to create a new Reader instance: %s", err)
	}
	if reflect.TypeOf(r) != reflect.TypeOf(&Reader{}) {
		t.Errorf("NewWriter does not return instance of Reader")
	}
	if r.Filename != filename {
		t.Errorf("`%s` doesn't match `%s`", r.Filename, filename)
	}
}

// TestReaderInit tests Reader contrustor
func TestReaderInit(t *testing.T) {
	path := _getPathToTests(t)
	filename := path + string(os.PathSeparator) + TestArchiveName
	r, err := NewReader(filename)
	if err != nil {
		t.Errorf("Unable to create a new Reader instance: %s", err)
	}
	if reflect.TypeOf(r.File) != reflect.TypeOf(&os.File{}) {
		t.Errorf("Unable to initialize the reader")
	}
}

// TestExtractFile tests extracting a file from archive
func TestExtractFile(t *testing.T) {
	// TODO: add test
}

// TestExtract tests extracting all files from archive
func TestExtract(t *testing.T) {
	// TODO: add test
}

// TestGetFilesCount tests enumerating the files in archive
func TestGetFilesCount(t *testing.T) {
	path := _getPathToTests(t)
	// create a new reader instance with the test archive
	r, err := NewReader(path + string(os.PathSeparator) + "test_archive.wpress")
	if err != nil {
		t.Errorf("Failed to create a new Reader instace: %s", err)
	}

	// get the files inside the archive
	filesCount, err := r.GetFilesCount()
	if err != nil {
		t.Errorf("Unable to get files count: %s", err)
	}
	// check if total files equals 3
	if 3 != filesCount {
		t.Errorf(
			"The archive contains %d files instead of 3",
			filesCount)
	}
}
