# **pwmfan-go**
A Pulse Width Modulation Fan For RaspberryPi written in Go  
Fan speed is adjust by CPU temperature, to reduce fan noise and power consumption  
  
[![GitHub version](https://img.shields.io/github/release/eZioPan/pwmfan-go.svg)](https://github.com/eZioPan/pwmfan-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/eZioPan/pwmfan-go)](https://goreportcard.com/report/github.com/eZioPan/pwmfan-go)
[![GoDoc](https://godoc.org/github.com/eZioPan/pwmfan-go?status.svg)](https://godoc.org/github.com/eZioPan/pwmfan-go)  

### **Warning:**
This program still in **EARLY** development, function may missing or not work properly, use with caution.  
This library still in **EARLY** development, APIs may change, use with caution.  

### **Pros**
1. Easy to use, just a binary file and a simple configuration file  
2. Easy install/uninstall with handy script  
3. Run as systemd service  
4. remote fan state monitor  

## **Hardware Requirement**
1. A RaspberryPi 1/2/3 serial  
2. A fan that support PWM signal, or you can build one by add a transistor (and resistor) to you origianl fan  

### **Remember:**
Raspberry Pi support max **+5V** for power pin, and use **+3.3V** for other gpio pin, **double check before connect anything to you pi**  

## **Software Requirement**
1. (If you want to use inbox system daemon configuration,) A [Systemd controlled system](https://en.wikipedia.org/wiki/Systemd#Availability).  
2. (If you want to build programe yourself,) A [Go SDK](https://golang.org/dl/).  

## **Usage**
Since controller program uses pwm to control fan, it only works with **root** privillage.  

You can get a very simple **installation help** when you run ./install.sh with no argument.  

### **Install Controller From Binary Release:**
1. Download release package from [Releases](https://github.com/eZioPan/pwmfan-go/releases) to your Pi  
```bash
cd controller
```  
2. Modify configuration file, see [Controller Configuration File](#controller-configuration-file) below  
3. Use
```bash
./install.sh install
```
to install controller program to your system  

### **Install Controller From Source:**
1. On your Pi  
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
cd $/GOPATH/src/github.com/eZioPan/pwmfan-go/cmd/controller
./install.sh build
```
2. Modify configuration file, see [Controller Configuration File](#controller-configuration-file) below  
3. Use
```bash
./install.sh install
```
to install controller program to your system  

### **Controller Configuration File:**
By default, the controller program use JSON style **"config.json"** file in the same directory of binary file, you can change it using -conifg command line parameter.  
Here is an example of controller configuration file:  
```json
{
	"Pin":18,
	"CPUTempPath":"/sys/class/thermal/thermal_zone0/temp",
	"StartCount":5,
	"StopCount":5,
	"StartTemp":38,
	"LowTemp":30,
	"HighTemp":45,
	"PwmFreq":600,
	"StartCycle":50,
	"LowCycle":50,
	"HighCycle":100,
	"FullCycle":100,
	"SampleRate":2,
	"NetworkInterfaceName":"lo",
	"ListenPort":2334,
	"Token":"123456789"
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
**SampleRate** is the rate that program check CPU temperature and change fan speed in.  
**NetworkInterfaceName** is the network interface which will be use to listen request and send fan data from  
**ListenPort**  is the network port which will be use to listen request and send fan data by  
**Token** is the token that will be check during fan data request  

### **Uninstall Controller:**
```bash
cd controller
./install.sh uninstall
```

### **Install Monitor From Binary Release:**
1. Download release package from [Releases](https://github.com/eZioPan/pwmfan-go/releases)  
```bash
cd monitor
```
2. Modify configuration file, see [Monitor Configuration File](#monitor-configuration-file) below  
3. Use
```bash
./monitor
```
to monitor your fan state

### **Install Monitor From Source:**
1. On your Pi  
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
cd $/GOPATH/src/github.com/eZioPan/pwmfan-go/cmd/monitor
go build -v -ldflags "-s -w"
```
2. Modify configuration file, see [Monitor Configuration File](#monitor-configuration-file) below  
3. Use
```bash
./monitor
```
to monitor your fan state  

### **Monitor Configuration File:**
By default, the monitor program use JSON style **"config.json"** file in the same directory of binary file, you can change it using -conifg command line parameter.  
Here is an example of monitor configuration file:  
```json
{
	"NetworkInterface":"lo",
	"remoteHost":"localhost",
	"remotePort": 2334,
	"Token":"123456789",
	"SampleRate":1
}
```
**NetworkInterface** is the network interface which will be use to send request and receive fan data from  
**remoteHost** is IP/hostname that your controller program exist  
**remotePort** is the port your controller program use to listen request and send data  
**Token** is the token that will be check during fan data request  
**SampleRate** is the rate that monitor will request fan state in  

### **Uninstall Monitor:**
Just remove binary file and configuration file

### **Tips:**
1. If you want to use fan monitor in a diffrent machine, please check your network/firewall settings.  
pwmfan use UDP to send and receive message
2. **Token** must be **THE SAME** both in controller configuration file and monitor configuration file.  
As for now, token is sent in plain text, **DO NO PUT ANY SENSITIVE DATA IN THIS FIELD**
3. As for now, fan data is sent in plain text, will change in future release

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
- Redesign configuration for more clear layer style  
- Add support for user defined fan data  
- Linux man page for commandline  
- send network data with encryption  
- [Done in v0.2.0] ~~Network function to read fan state~~  
- [Done in v0.2.0] ~~run as a Systemd service~~  
- [WON'T DO] ~~Network function to change fan parameter~~(Too dangerous)  
  
## **Changelog**  
- v0.2.0  
Support systemd service  
Support remote fan monitor  
Support SampleRate parameter  

- v0.1.0  
Reconstruct code  
Support StartCount and StopCount parameter  
  
- v0.0.1  
Init Release  
Basic function  
