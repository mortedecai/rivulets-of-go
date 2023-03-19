package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/mortedecai/rivulets-of-go/internal/version"
	"github.com/mortedecai/rivulets-of-go/server/connection"
)

func printLogWelcome() {
	fmt.Println("Rivulets of Go")
	fmt.Println("")
	fmt.Println("Version: " + version.Version)
	fmt.Println("Commit:  " + version.Commit)
}

func printHelloAndExit(conn *connection.Data) *connection.Data {
	helloString := fmt.Sprintf("Rivulets of Go\n\r\n\rVersion: %s\n\rCommit: %s\n\r\n\rUnder Construction. Good Bye\n\r\n\r", version.Version, version.Commit)
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
	return conn
}

const portVal = ":3160"

func main() {
	var logger *zap.SugaredLogger
	var mgr connection.Manager
	printLogWelcome()

	if dl, err := zap.NewDevelopment(); err == nil {
		logger = dl.Sugar().Named("RoG")
	} else {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("ERROR:  Could not create logger:  ", err.Error(), ".")
		fmt.Println("")
		fmt.Println("Terminating.")
		os.Exit(2)
	}

	logger.Debugw("Creating Connection Manager", "Port", portVal)

	if m, err := connection.NewManager(portVal, logger); err != nil {
		logger.Errorw("Creating Connection Manager", "Error", err)
	} else {
		mgr = m
		logger.Debugw("Connection Manager Created", "Manager", mgr)
	}

	mgr.SetMaintenanceHandler(printHelloAndExit)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		fmt.Println()
		done <- true
	}()
	mgr.MaintenanceStart()
	<-done
	fmt.Println("exiting")
}
