package main

import (
	"io"
	"io/ioutil"
	"os"
)

// START OMIT

func main() {
	// notice how we are able to seamlessly use higher-level libraries
	err := ioutil.WriteFile("/tmp/hello.txt", []byte("Example 06\n"), 0777)
	if err != nil {
		panic(err)
	}

	f, err := os.Open("/tmp/hello.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
	f.Close()
}

// END OMIT
