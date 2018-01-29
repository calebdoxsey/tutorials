package main

import "fmt"

// START OMIT

func main() {
	fmt.Print("enter your name: ")
	var name string
	fmt.Scanln(&name)
	fmt.Println("\n\nyou entered:", name)
}

// END OMIT
