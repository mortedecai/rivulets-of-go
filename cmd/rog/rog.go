package main

import "fmt"

var (
	Version string = "N/A"
	Commit  string = "N/A"
)

func printLogWelcome() {
	fmt.Println("Rivulets of Go")
	fmt.Println("")
	fmt.Println("Version: " + Version)
	fmt.Println("Commit:  " + Commit)
}

func main() {
	printLogWelcome()
}
