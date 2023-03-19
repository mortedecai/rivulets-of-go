package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/mortedecai/rivulets-of-go/server/connection"
	"github.com/mortedecai/rivulets-of-go/server/info"
)

func printLogWelcome() {
	fmt.Println(info.Name)
	fmt.Println("")
	fmt.Println("Version: " + info.Version)
	fmt.Println("Commit:  " + info.Commit)
}

func printHello(conn *connection.Data) *connection.Data {
	helloString := fmt.Sprintf("%s\n\r\n\rVersion: %s\n\rCommit: %s\n\r\n\rUnder Construction. Good Bye\n\r\n\r", info.Name, info.Version, info.Commit)
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
	const methodName = "main"
	var logger *zap.SugaredLogger
	var mgr connection.Manager
	printLogWelcome()

	//if dl, err := zap.NewDevelopment(); err == nil {
	if dl, err := zap.NewProduction(); err == nil {
		logger = dl.Sugar().Named("RoG")
	} else {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("ERROR:  Could not create logger:  ", err.Error(), ".")
		fmt.Println("")
		fmt.Println("Terminating.")
		os.Exit(2)
	}

	logger.Debugw(methodName, "Port", portVal)

	if m, err := connection.NewManager(portVal, logger); err != nil {
		logger.Errorw(methodName, "Manager Error", err)
	} else {
		mgr = m
		logger.Debugw(methodName, "Manager", mgr)
	}

	mgr.SetMaintenanceHandler(printHello)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		logger.Infow(methodName, "Signal", sig)
		done <- true
	}()
	if err := mgr.MaintenanceStart(); err != nil {
		logger.Errorw(methodName, "Error", err)
		os.Exit(2)
	}
	info.UpDate = time.Now().Local()
	logger.Infow(methodName, info.Name, "ONLINE", "Time", info.UpDate.Format(time.RFC1123))
	<-done
	fmt.Println("exiting")
}
