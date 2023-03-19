package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

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

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		fmt.Println()
		done <- true
	}()

	<-done
	fmt.Println("exiting")
}
