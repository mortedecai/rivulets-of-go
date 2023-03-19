package main

import (
	"fmt"
	"net"
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

func printHelloAndExit(conn net.Conn) {
	helloString := fmt.Sprintf("Rivulets of Go\n\r\n\rVersion: %s\n\rCommit: %s\n\r\n\rUnder Construction. Good Bye\n\r\n\r", Version, Commit)
	data := []byte(helloString)

	totalBytes := 0
	bw := 0

	for (totalBytes + bw) < len(helloString) {
		bw, err := conn.Write(data[totalBytes:])
		totalBytes += bw
		if err != nil {
			fmt.Println("Error writing hello string:  ", err.Error())
		}
	}
	conn.Close()
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

	//portString := ":3160"

	address := net.TCPAddr{Port: 3160}

	listener, err := net.ListenTCP("tcp", &address)
	if err != nil {
		fmt.Println("Error creating TCP Listener:  ", err.Error())
		os.Exit(1)
	}

	terminateMUD := false

	go func() {
		for {
			if conn, err := listener.Accept(); err != nil {
				fmt.Println("Error accepting connection:  ", err.Error())
			} else {
				defer conn.Close()
				fmt.Println("Accepting connection and printing hello.")
				printHelloAndExit(conn)
			}
			if terminateMUD {
				break
			}
		}
	}()

	<-done
	fmt.Println("exiting")
}
