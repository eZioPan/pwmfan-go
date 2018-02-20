# **pwmfan-go**
A Pulse Width Modulation Fan For RaspberryPi written in Go  
Fan speed is adjust by CPU temperature, to reduce fan noise and power consumption  

### **Warning:**  
This program still in **VERY EARLY** development, function may missing or not work properly, use with caution  

## **Usage:**
Because this program use pwm to control fan, it only works with **root** privillage

### **Install from binary release:**
Download release package from [Releases](https://github.com/eZioPan/pwmfan-go/releases), modify configuration file and run

### **Install from source:**
```bash
go get -v -u -d github.com/eZioPan/pwmfan-go
cd $/GOPATH/src/github.com/eZioPan/pwmfan-go
GOOS=linux GOARCH=arm go build
```
Then modify configuration file and run

## **Confugration File**
Currently we only support **config.json** file in the **same directory** of binary file  
Here is an example of config.json file
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
**~~StartCount~~** , **~~StopCount~~** not use in this version  
**StartTemp** when temperature higher than this value, fan will start  
**LowTemp** when temperature lower than this value, fan will stop
**HighTemp** when temperature higher than this value, fan will run in **HighCycle** defined speed  
**PwmFreq** PWM signal frequence, should be set high enough for PWM signal  
**StartCycle** when fan start from full stop, it needs a little higher speed to kick it start  
**LowCycle** is the lowest speed fan will run, before fan fully stopped  
**HighCycle** is the highest speed fan will run. This value shoudn't greater than **FullCycle**  
**FullCycle** is maximum speed which fan is able to run at

## **Dependencices:**
[stianeikeland/go-rpio](https://github.com/stianeikeland/go-rpio) for accessing RasberryPi gpio in pure go
