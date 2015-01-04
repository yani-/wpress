[![Build Status](https://travis-ci.org/yani-/wpress.svg?branch=master)](https://travis-ci.org/yani-/wpress)
[![Coverage Status](https://coveralls.io/repos/yani-/wpress/badge.png?branch=master)](https://coveralls.io/r/yani-/wpress?branch=master)
[![GoDoc](https://godoc.org/github.com/yani-/wpress?status.svg)](https://godoc.org/github.com/yani-/wpress)

# wpress

WordPress archive format

## Quick Start
```
// create new archive and add a file to it
archiver, _ := NewWriter("test.wpress")
archiver.AddFile("file-to-add.txt")
archiver.Close()

// create a new archive and add a directory to it
archiver, _ := NewWriter("test.wpress")
archiver.AddDirectory("/path/to/directory/to/add")
archiver.Close()
```

## License

This project is licensed under the MIT open source license.
