package pwmfan

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
)

// HandleErr will try to handle an error
func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseJSON parse json file into a Config structure
func ParseJSON(cfgFilePath string) (cfg Config) {
	cfgFile, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0644)
	HandleErr(err)
	defer cfgFile.Close()
	jsd := json.NewDecoder(cfgFile)
	cfg = Config{}
	err = jsd.Decode(&cfg)
	HandleErr(err)
	return cfg
}

// ReadCPUTemperature read temprature once from a file and divie raw data by a divider
/*
Note: As for raspberry Pi,
CPUTempPath is /sys/class/thermal/thermal_zone0/temp
Divider should be set to 1000
*/
func ReadCPUTemperature(CPUTempPath string, Divider float64) (Temperature float64) {
	tempBuf, err := ioutil.ReadFile(CPUTempPath)
	HandleErr(err)
	raw, err := strconv.ParseFloat(string(tempBuf[:len(tempBuf)-1]), 64)
	HandleErr(err)
	Temperature = raw / Divider
	return Temperature
}

// PullCPUTemp read from a system file constantanious and put temperature  into a predifined channel
/*
Note: As for raspberry Pi,
CPUTempPath is /sys/class/thermal/thermal_zone0/temp
Divider should be set to 1000
*/
func PullCPUTemp(CPUTempPath string, Divider float64, Temp chan<- float64) {
	for {
		Temp <- ReadCPUTemperature(CPUTempPath, Divider)
	}
}

// SignalProcess listen from a os.Notify function's output, when catch signal,
// Process the Action, then send original to the Process
func SignalProcess(Process *os.Process, SigChan <-chan os.Signal, Action func()) {
	sig := <-SigChan
	Action()
	signal.Reset()
	Process.Signal(sig)
}

// LinearRemap is a linear remap function
func LinearRemap(inputs []float64, opt ...float64) (outputs []float64) {
	res := (opt[3]-opt[2])/(opt[1]-opt[0])*inputs[0] + (opt[2] - opt[0]*(opt[3]-opt[2])/(opt[1]-opt[0]))
	outputs[0] = res
	return outputs
}

// LinearClampRemap is a linear remap function that will clamp output in between opt[2] and opt[3]
func LinearClampRemap(inputs []float64, opt ...float64) (outputs []float64) {
	var res float64
	if inputs[0] <= opt[0] {
		res = opt[2]
	} else if inputs[0] <= opt[1] {
		res = LinearRemap(inputs, opt[0], opt[1], opt[2], opt[3])[0]
	} else {
		res = opt[3]
	}
	outputs[0] = res
	return outputs
}

// TriangularWave is a linear remap will read input and generate a oscillation wave
func TriangularWave(small float64, big float64, step float64, ch chan<- float64) {
	if small > big {
		big, small = small, big
	}
	for {
		for i := small; i <= big; i += step {
			ch <- i
		}
		for i := big; i >= small; i -= step {
			ch <- i
		}
	}
}
