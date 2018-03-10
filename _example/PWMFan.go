package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eZioPan/pwmfan-go"
	"github.com/stianeikeland/go-rpio"
)

var (
	configPath = ""
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "system temperature file")
}

func main() {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()
	fan0 := pwmfan.NewFan(pwmfan.ParseJSON(configPath))
	action := func() {
		fmt.Println("Catch Signal, perpare to exit")
		fan0.Pin.Output()
		fan0.Pin.Low()
		rpio.Close()
		fmt.Println("Terminating")
	}
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}
	go pwmfan.SignalProcess(p, sigChan, action)

	fan0.Monitor()

}
