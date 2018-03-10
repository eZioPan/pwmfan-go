# **pwmfan-go**
A Pulse Width Modulation Fan For RaspberryPi written in Go  
Fan speed is adjust by CPU temperature, to reduce fan noise and power consumption  
[![Go Report Card](https://goreportcard.com/badge/github.com/eZioPan/pwmfan-go)](https://goreportcard.com/report/github.com/eZioPan/pwmfan-go)
[![GoDoc](https://godoc.org/github.com/eZioPan/pwmfan-go?status.svg)](https://godoc.org/github.com/eZioPan/pwmfan-go)  
  
### **Warning:**
This program still in **VERY EARLY** development, function may missing or not work properly, use with caution.  
This library still in **VERY EARLY** development, APIs may change, use with caution.  
  
### **Pros**
1. Easy to use, just a binary file and a simple configuration file  
2. Easy to compile and/or cross compile  
3. Easy development  

### **Cons**
1. (Compare to pure C/C++ program,) A slightly big binary file size
2. (Compare to pure C/C++ program,) A slightly big memory consumption

## **Hardware Requirement**
1. A Raspberry Pi2 or Raspberry Pi3 serial  
2. A fan that support PWM signal, or you can build one by add a transistor to you origianl fan  
  
### **Remember:**
Raspberry Pi support max **+5V** for power pin, and use **+3.3V** for other gpio pin, **double check before connect anything to you pi**  
  
## **Usage**
Because this program uses pwm to control fan, it only works with **root** privillage  
  
### **Install from binary release:**
Download release package from [Releases](https://github.com/eZioPan/pwmfan-go/releases) to your Pi, then modify configuration file and run.  
  
### **Install from source:**
On your Pi  
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
cd $/GOPATH/src/github.com/eZioPan/pwmfan-go/_example
GOOS=linux GOARCH=arm go build -o PWMFan
```
Then modify configuration file and run.  
  
## **Confugration File**
By default, the program use JSON style **"config.json"** file in the same directory of binary file, you can change it using -conifg command line parameter.  
Here is an example of configuration file:  
```json
{
	"Pin":18,
	"CPUTempPath":"/sys/class/thermal/thermal_zone0/temp",
	"StartCount":5,
	"StopCount":5,
	"StartTemp":33,
	"LowTemp":30,
	"HighTemp":40,
	"PwmFreq":600,
	"StartCycle":50,
	"LowCycle":80,
	"HighCycle":100,
	"FullCycle":100
}
```  

**Pin** defines the gpio pin number of PWM. This uses BCM Pinout number, you can check [here](https://pinout.xyz).  
**CPUTempPath** defines the system file path to read cpu's temperature.  
**StartCount** after counter reaches this count fan will start.  
**StopCount** after counter reaches this count fan will stop.  
**StartTemp** when temperature higher than this value, fan will start.  
**LowTemp** when temperature lower than this value, fan will stop.  
**HighTemp** when temperature higher than this value, fan will run in **HighCycle** defined speed.  
**PwmFreq** PWM signal frequence, should be set high enough for PWM signal.  
**StartCycle** when fan start from full stop, it needs a little higher speed to kick it start.  
**LowCycle** is the lowest speed fan will run, before fan fully stopped.  
**HighCycle** is the highest speed fan will run. This value shoudn't greater than **FullCycle**.  
**FullCycle** is maximum speed which fan is able to run at.  
  
## **For Developers**
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
```
```go
import "github.com/eZioPan/pwmfan-go"
```  
  
## **Dependencices**
[stianeikeland/go-rpio](https://github.com/stianeikeland/go-rpio) for accessing RasberryPi gpio in pure go  
  
## **TODO**
- Write fan log to system log  
- Network function to read fan state  
- Network function to change fan parameter  
  
## **Changelog**  
- v0.1.0  
Reconstruct code  
Support StartCount and StopCount parameter  
  
- v0.0.1  
Init Release  
Basic function  
