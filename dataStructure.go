package pwmfan

import (
	"net"

	"github.com/stianeikeland/go-rpio"
)

// TempPair represent a temperature and correspond Cycle and Count.
//
// TempPair struct can be static, as a configuration. Also can be dynamic, as a representation of current fan state.
type TempPair struct {
	Temp  float64
	Cycle uint
	Count uint
}
// NetworkSettings represents a basic network settings
type NetworkSettings struct {
	InterfaceName string
	ListenPort    uint
	Token         string
}
// Config contain all information that define a PWM fan's user configuration
//
// TODO: try to implement String() method
type Config struct {
	Pin         uint
	CPUTempPath string
	PwmFreq     uint
	FullCycle   uint
	SampleRate  uint
	Start       TempPair
	Stop        TempPair
	High        TempPair
	NetworkSettings
}

// State represent Fan's current running state
//
// TODO: try to implement String() method
type State uint

const (
	// Stop represent Pwm Fan stay in a full stop state
	Stop State = iota
	// Start represent Pwm Fan enter in s Start state
	Start
	// Run represent Pwm Fan stay in a running state
	Run
)

// Fan represent a Fan
//
// State, Cycle, Temp should be updated 'in realtime'
// Pin better not be modified after call NewFan()
//
// TODO: try to implement String() method
type Fan struct {
	Pin     rpio.Pin
	Current TempPair
	StateRecord
	Cfg     Config
	UDPAddr *net.UDPAddr
}

// StateRecord represents a Fan's state and state's counter
type StateRecord struct {
	State
	StopCounter  uint
	StartCounter uint
}

func (state State) String() string {
	var str string
	switch state {
	case Stop:
		str = "Stop"
	case Start:
		str = "Start"
	case Run:
		str = "Run"
	}
	return str
}

// RemapFunc is a function that read at least one float64 input and output a float64 data
//
// More specific,this function should take 1st argument as an input data(s),then read more input, and output result data(s)
//
// Example: LinearRemap and LinearClampRemap
type RemapFunc func([]float64, ...float64) []float64
