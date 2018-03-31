package pwmfan

import (
	"encoding/binary"
	"math"
	"net"
	"strconv"

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

// String implements fmt.Stringer interface
func (tp TempPair) String() (res string) {
	var tmpstr, cylstr, cntstr string
	tmpstr = strconv.FormatFloat(tp.Temp, 'f', -1, 64)
	cylstr = strconv.FormatUint(uint64(tp.Cycle), 10)
	cntstr = strconv.FormatUint(uint64(tp.Count), 10)
	res = "Temp: " + tmpstr + "\tCycle: " + cylstr + "\tCount: " + cntstr
	return res
}

// MarshalBinary implements encoding.BinaryMarshaler interface
//
// As for TempPair, this method will output a 24 bytes binary flow.
//
// First 8 bytes for temperature, second 8 bytes for cycle, last 8 bytes for count.
func (tp TempPair) MarshalBinary() (data []byte, err error) {
	var tmpbin, cylbin, cntbin []byte
	binary.BigEndian.PutUint64(tmpbin, math.Float64bits(tp.Temp))
	binary.BigEndian.PutUint64(cylbin, uint64(tp.Cycle))
	binary.BigEndian.PutUint64(cntbin, uint64(tp.Count))
	data = append(append(tmpbin, cylbin...), cntbin...)
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
//
// As for TempPair, this method will try to parse a 24 bytes binary flow.
//
// First 8 bytes for temperature, second 8 bytes for cycle, last 8 bytes for count.
func (tp *TempPair) UnmarshalBinary(data []byte) error {
	var (
		tmp      float64
		cyl, cnt uint
	)
	tmp = math.Float64frombits(binary.BigEndian.Uint64(data[:8]))
	cyl = uint(binary.BigEndian.Uint64(data[8:16]))
	cnt = uint(binary.BigEndian.Uint64(data[16:]))
	tp.Temp = tmp
	tp.Cycle = cyl
	tp.Count = cnt
	return nil
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
