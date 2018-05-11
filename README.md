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
1. A Raspberry Pi 1/2/3 serial  
2. A fan that support PWM signal, or you can build one by add a transistor (and resistor) to you origianl fan  

### **Remember:**
Raspberry Pi support max **+5V** for power pin, and use **+3.3V** for other gpio pin, **double check before connect anything to you pi**.  

## **Software Requirement**
1. Raspberry Pi runs with linux.
2. [OPTIONAL] If you want to use inbox system daemon configuration, a [Systemd controlled system](https://en.wikipedia.org/wiki/Systemd#Availability).  
3. [OPTIONAL] If you want to build programe yourself, a [Go SDK](https://golang.org/dl/).  

## **Basic Principle**
1. This program uses PWM signal to adjust hardware fan speed and run state.
2. This program reads temperature data from system file, and compare with configuration file, then change fan speed and state at **SampleRate**.
3. This program uses **ratio speed** to control fan speed, not hardware real speed.  
Which means **FullCycle** in configuration file corresponds to the maximum capablility of fan speed.  
And the **ratio** of **xxxCycle** divided by **FullCycle** define fan's speed.
4. There 3 states for fan, **Stop**, **Start**, **Run**.
5. In **Stop** state, fan will run at **StopCycle** speed, which usually be 0.  
If temperature is higher than **Temp**, this program will add one into count, otherwise remove one from count.  
If count reaches **Count** settings, fan will shift into **Start** state.
6. In **Start** state, fan will run at **Cycle** speed, then shift into **Run** state.  
This state designs for starting fan to run properly, not too fast nor too slow.
7. In **Run** state, fan's speed is adjust between a high **Cycle** and a low **Cycle**, with temperature as a ruler.  
If temperature is lower than **Temp**, this program will add one into count, otherwise remove one from count.  
If count reaches **Count** settings, fan will shift into **Stop** state.  
Then fan will loop into step 5.

## **Usage**
Since controller program uses pwm to control fan, it only works with **root** privillage.  

You can get a very simple **installation help** when you run ./install.sh with no argument.  

### **Install Controller From Binary Release:**
1. Download release package from [Releases](https://github.com/eZioPan/pwmfan-go/releases) to your Pi.
```bash
cd controller
```  
2. Modify configuration file, see [Controller Configuration File](#controller-configuration-file) below.
3. Use
```bash
./install.sh install
```
to install controller program to your system.

### **Install Controller From Source:**
1. On your Pi  
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
cd $/GOPATH/src/github.com/eZioPan/pwmfan-go/cmd/controller
./install.sh build
```
2. Modify configuration file, see [Controller Configuration File](#controller-configuration-file) below.
3. Use
```bash
./install.sh install
```
to install controller program to your system.

### **Controller Configuration File:**
By default, the controller program use JSON style **"config.json"** file in the same directory of binary file, you can change it using **-conifg** command line parameter.  
Here is an example of controller configuration file:  
```json
{
	"Pin":18,
	"CPUTempPath":"/sys/class/thermal/thermal_zone0/temp",
	"SampleRate":2,
	"PwmFreq":600,
	"FullCycle":100,
	"StopCycle":0,
	"Start":{
		"Temp":38,
		"Cycle":50,
		"Count":5
	},
	"Low":{
		"Temp":30,
		"Cycle":50,
		"Count":5
	},
	"High":{
		"Temp":45,
		"Cycle":100,
		"Count":0
	},
	"NetworkSettings":{
		"InterfaceName":"lo",
		"ListenPort":2334,
		"Token":"123456789"
	}
}
```
**Pin** defines the gpio pin number of PWM. This uses BCM Pinout number, you can check [here](https://pinout.xyz).  
**CPUTempPath** defines the system file path to read cpu's temperature.  
**SampleRate** is the rate that program check CPU temperature and change fan speed in.  
**PwmFreq** PWM signal frequence, should be set high enough for PWM signal.  
**FullCycle** is maximum speed which fan is able to run at.  
**StopCycle** is the speed when fan stays in **Stop** state, usually this value should be 0.  
**Stop**,**Start**,**Run** as described in [Basic Principle](#basic-principle) earlier, is the configurations for each state.  
**NetworkSettings** contains all network configuration for sending and receiving information.  
**InterfaceName** is the network interface which will be use to listen request and send fan data from.  
**ListenPort**  is the network port which will be use to listen request and send fan data by.  
**Token** is the token that will be check during fan data request.  

### **Uninstall Controller:**
```bash
cd controller
./install.sh uninstall
```

### **Install Monitor From Binary Release:**
1. Download release package from [Releases](https://github.com/eZioPan/pwmfan-go/releases).  
```bash
cd monitor
```
2. Modify configuration file, see [Monitor Configuration File](#monitor-configuration-file) below.  
3. Use
```bash
./monitor
```
to monitor your fan state.

### **Install Monitor From Source:**
1. On your Pi  
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
cd $/GOPATH/src/github.com/eZioPan/pwmfan-go/cmd/monitor
go build -v -ldflags "-s -w"
```
2. Modify configuration file, see [Monitor Configuration File](#monitor-configuration-file) below.  
3. Use
```bash
./monitor
```
to monitor your fan state.  

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
**NetworkInterface** is the network interface which will be use to send request and receive fan data from.  
**remoteHost** is IP/hostname that your controller program exist.  
**remotePort** is the port your controller program use to listen request and send data.  
**Token** is the token that will be check during fan data request.  
**SampleRate** is the rate that monitor will request fan state in.  

### **Uninstall Monitor:**
Just remove binary file and configuration file.

### **Tips:**
1. If you want to use fan monitor in a diffrent machine, please check your network/firewall settings.  
pwmfan use **UDP** to send and receive message.
2. **Token** must be **THE SAME** both in controller configuration file and monitor configuration file.  
As for now, token is sent in plain text, **DO NO PUT ANY SENSITIVE DATA IN THIS FIELD**.
3. As for now, fan data is sent in plain text, will change in future release.

## **For Developers**
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
```
```go
import "github.com/eZioPan/pwmfan-go/common"
import "github.com/eZioPan/pwmfan-go"
```

## **Dependencices**
[stianeikeland/go-rpio](https://github.com/stianeikeland/go-rpio) for accessing RasberryPi gpio in pure go.  
  
## **TODO** 
- Write fan log to system log  
- Add support for user defined fan data  
- Linux man page for commandline  
- send network data with encryption  
- [Done in v0.3.0] ~~Redesign configuration for more clear layer style~~  
- [Done in v0.2.0] ~~Network function to read fan state~~  
- [Done in v0.2.0] ~~run as a Systemd service~~  
- [WON'T DO] ~~Network function to change fan parameter~~(Too dangerous)  
  
## **Changelog**
- v0.3.0  
Redesign configuration file into more clear layer style  
Extract common code into sub directory for reusing  
Bugs fix

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
