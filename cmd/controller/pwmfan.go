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
	flag.Parse()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	err := rpio.Open()
	pwmfan.HandleErr(err)
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
	pwmfan.HandleErr(err)
	go pwmfan.SignalProcess(p, sigChan, action)

	srvConn := fan0.CreateServer()
	go fan0.HandleRequest(srvConn)

	fan0.Monitor()

}
