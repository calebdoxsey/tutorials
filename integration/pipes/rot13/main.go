package main

import "os"

func rot13(b byte) byte {
	var a, z byte
	switch {
	case 'a' <= b && b <= 'z':
		a, z = 'a', 'z'
	case 'A' <= b && b <= 'Z':
		a, z = 'A', 'Z'
	default:
		return b
	}
	return (b-a+13)%(z-a+1) + a
}

func main() {
	bs := []byte{0}
	for {
		_, err := os.Stdin.Read(bs)
		if err != nil {
			break
		}
		bs[0] = rot13(bs[0])
		_, err = os.Stdout.Write(bs)
		if err != nil {
			break
		}
	}
}
