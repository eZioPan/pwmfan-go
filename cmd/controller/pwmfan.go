package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/eZioPan/pwmfan-go"
	"github.com/eZioPan/pwmfan-go/common"
	"github.com/stianeikeland/go-rpio"
)

var (
	configPath = ""
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "controller configuration file path")
}

func main() {
	flag.Parse()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	err := rpio.Open()
	common.HandleErr(err)
	defer rpio.Close()
	cfg := &common.Config{}
	common.ParseJSON(configPath, cfg)
	fan0 := pwmfan.NewFan(*cfg)
	action := func() {
		fmt.Println("Catch Signal, perpare to exit")
		rpio.Pin(fan0.Pin).Output()
		rpio.Pin(fan0.Pin).Low()
		rpio.Close()
		fmt.Println("Terminating")
	}
	p, err := os.FindProcess(os.Getpid())
	common.HandleErr(err)
	go common.SignalProcess(p, sigChan, action)

	srvConn := fan0.CreateServer()
	go fan0.HandleRequest(srvConn)

	pwmfan.Monitor(fan0)

}
