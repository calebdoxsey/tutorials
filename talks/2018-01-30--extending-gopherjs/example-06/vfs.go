package main

var minusOne = -1

// Our virtual file system contains files and references to files
// A file is just a slice of bytes
// A reference also tracks the position within the file

// START OMIT

type (
	virtualFile struct {
		data []byte
	}
	virtualFileReference struct {
		file *virtualFile
		pos  int
	}
	virtualFileSystem struct {
		files  map[string]*virtualFile
		fds    map[uintptr]*virtualFileReference
		nextFD uintptr
	}
)

func newVirtualFileSystem() *virtualFileSystem {
	return &virtualFileSystem{
		files:  make(map[string]*virtualFile),
		fds:    make(map[uintptr]*virtualFileReference),
		nextFD: 1000,
	}
}

// END OMIT
