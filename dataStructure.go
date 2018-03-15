package pwmfan

import (
	"net"

	"github.com/stianeikeland/go-rpio"
)

//Config contain all information that define a PWM fan's user configuration
// TODO: try to implement String() method
type Config struct {
	Pin                  uint
	CPUTempPath          string
	StartCount           uint
	StopCount            uint
	StartTemp            float64
	LowTemp              float64
	HighTemp             float64
	PwmFreq              uint
	StopCycle            uint
	StartCycle           uint
	LowCycle             uint
	HighCycle            uint
	FullCycle            uint
	SampleRate           uint
	NetworkInterfaceName string
	ListenPort           int
	Token                string
}

// FanState represent a PWM Fan's running state
// TODO: try to implement String() method
type FanState uint

const (
	// Stop represent Pwm Fan stay in a full stop state
	Stop FanState = iota
	// Start represent Pwm Fan enter in s Start state
	Start
	// Run represent Pwm Fan stay in a running state
	Run
)

// Fan represent a Fan
/*
State, Cycle, Temp should be updated 'in realtime'
Pin better not be modified after call NewFan()
*/
// TODO: try to implement String() method
type Fan struct {
	State        FanState
	Cycle        uint
	Temp         float64
	Cfg          Config
	Pin          rpio.Pin
	StopCounter  uint
	StartCounter uint
	UDPAddr      *net.UDPAddr
}

// RemapFunc is a function that read at least one float64 input and output a float64 data
/*
More specific,
this function should take 1st argument as an input data,
then read more input, and output an data
Example: LinearRemap and LinearClampRemap
*/
type RemapFunc func(float64, ...float64) float64
