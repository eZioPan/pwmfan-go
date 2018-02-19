package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	CPUTempPath = "/sys/class/thermal/thermal_zone0/temp"
)

//PwmFan struct
type PwmFan struct {
	Pin         uint
	CPUTempPath string
	StartCount  uint
	StopCount   uint
	StartTemp   uint
	LowTemp     uint
	HighTemp    uint
	PwmFreq     uint
	StartCycle  uint
	LowCycle    uint
	HighCycle   uint
	FullCycle   uint
}

func main() {
	pwmfancfg := parsejson("config.json")
	fmt.Println(pwmfancfg)
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()
	pin := rpio.Pin(pwmfancfg.Pin)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-sigChan:
			fmt.Println("Catch Signal, prepare to exit")
			pin.Output()
			pin.Low()
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				panic(err)
			}
			p.Signal(os.Kill)
		}
	}()

	pin.Pwm()
	pin.Freq(int(pwmfancfg.PwmFreq))
	pin.DutyCycle(uint32(pwmfancfg.StartCycle), uint32(pwmfancfg.FullCycle))
	time.Sleep(time.Second)

	tempChan := make(chan float64)
	go readCPUTemp(tempChan)
	defer func() {
		pin.Output()
		pin.Low()
	}()
	cycle := float64(0)
	temp := float64(0)
	LowCount := int(0)
	StartCount := int(0)
	StopState := false
	for {
		temp = <-tempChan
		cycle = linearClampRemap(<-tempChan, float64(pwmfancfg.LowTemp), float64(pwmfancfg.HighTemp), float64(pwmfancfg.LowCycle), float64(pwmfancfg.HighCycle))
		if cycle == float64(pwmfancfg.LowCycle) && LowCount < 5 {
			LowCount++
		} else if cycle > float64(pwmfancfg.LowCycle) && LowCount > 0 {
			LowCount--
		}
		if LowCount >= 5 {
			StopState = true
			LowCount = 0
		}
		if StopState == true && temp >= float64(pwmfancfg.StartTemp) && StartCount < 5 {
			StartCount++
		} else if StopState == true && temp < float64(pwmfancfg.StartTemp) && StartCount > 0 {
			StartCount--
		}
		if StopState == true && StartCount >= 5 {
			StopState = false
			StartCount = 0
			cycle = float64(pwmfancfg.StartCycle)
			pin.DutyCycle(uint32(cycle), uint32(pwmfancfg.FullCycle))
			time.Sleep(time.Second)
		}
		if StopState == true {
			cycle = 0
		}
		fmt.Println("Temp:", temp, "Cycle:", cycle, "StartCount:", StartCount, "LowCount:", LowCount)
		pin.DutyCycle(uint32(cycle), uint32(pwmfancfg.FullCycle))
		time.Sleep(time.Millisecond * 500)
	}
}

func pingpong(small int, big int, ch chan<- int) {
	if small > big {
		big, small = small, big
	}
	for {
		for i := small; i <= big; i++ {
			ch <- i
		}
		for i := big; i >= small; i-- {
			ch <- i
		}
	}
}

func parsejson(cfgFilePath string) (cfg *PwmFan) {
	cfgFile, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer cfgFile.Close()
	jsd := json.NewDecoder(cfgFile)
	cfg = new(PwmFan)
	err = jsd.Decode(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func readCPUTemp(temp chan<- float64) {
	for {
		tempBuf, err := ioutil.ReadFile(CPUTempPath)
		if err != nil {
			panic(err)
		}
		raw, err := strconv.ParseFloat(string(tempBuf[:len(tempBuf)-1]), 64)
		temp <- raw / float64(1000)
	}
}

func linearClampRemap(input, xLow, xHigh, yLow, yHigh float64) (y float64) {
	if input <= xLow {
		return yLow
	} else if input <= xHigh {
		y = (yHigh-yLow)/(xHigh-xLow)*input + (yLow - xLow*(yHigh-yLow)/(xHigh-xLow))
		return y
	} else {
		return yHigh
	}
}
