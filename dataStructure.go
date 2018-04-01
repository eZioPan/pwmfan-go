package pwmfan

import (
	"encoding/binary"
	"errors"
	"math"
	"net"
	"strconv"
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
// As for TempPair, this method will output a 25 bytes binary flow.
//
// First 1 byte for total length of the binary flow, its 24 for TempPair, First 8 bytes for temperature, second 8 bytes for cycle, last 8 bytes for count.
func (tp TempPair) MarshalBinary() (data []byte, err error) {
	const totalLength = uint64(3 * 8)
	var total, tmpbin, cylbin, cntbin []byte
	binary.BigEndian.PutUint64(tmpbin, math.Float64bits(tp.Temp))
	binary.BigEndian.PutUint64(cylbin, uint64(tp.Cycle))
	binary.BigEndian.PutUint64(cntbin, uint64(tp.Count))
	binary.BigEndian.PutUint64(total, totalLength)
	data = append(append(append(total, tmpbin...), cylbin...), cntbin...)
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

// String implements fmt.Stringer interface
func (ns NetworkSettings) String() (res string) {
	lpt := strconv.FormatUint(uint64(ns.ListenPort), 10)
	res = "InterfaceName: " + ns.InterfaceName + "\tListenPort: " + lpt + "\tToken: " + ns.Token
	return res
}

// MarshalBinary implements encoding.BinaryMarshaler interface
//
// As for NetworkSettings, this method will output a dynamic length binary flow.
//
// First part for NetworkSettings.InterfaceName, using StringToBinary() for converting.
// Second part for NetworkSettings.ListenPort, 8 bytes.
// Last part for NetworkSettings.Token, using StringToBinary() for converting.
func (ns NetworkSettings) MarshalBinary() (data []byte, err error) {
	var ifn, lpt, tkn []byte
	ifn, _, err = StringToBinary(ns.InterfaceName)
	if err != nil {
		return nil, err
	}
	binary.BigEndian.PutUint64(lpt, uint64(ns.ListenPort))
	tkn, _, err = StringToBinary(ns.Token)
	if err != nil {
		return nil, err
	}
	data = append(append(ifn, lpt...), tkn...)
	return data, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler interface
//
// As for NetworkSettings, this method will read a binary flow, parse it into a NetworkSettings
//
// First part for NetworkSettings.InterfaceName, using BinaryToString() for converting.
// Second part for NetworkSettings.ListenPort, 8 bytes.
// Last part for NetworkSettings.Token, using BinaryToString() for converting.
func (ns *NetworkSettings) UnmarshalBinary(data []byte) (err error) {
	var (
		ifn, tkn string
		lpt      uint
	)
	ifn, n, _ := BinaryToString(data)
	lpt = uint(binary.BigEndian.Uint64(data[n : n+8]))
	tkn, _, _ = BinaryToString(data[n+8:])
	ns.InterfaceName = ifn
	ns.ListenPort = lpt
	ns.Token = tkn
	return nil
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
	Pin     uint8
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

// StringToBinary will read a string and convert it into binary flow data, whole data bytes size n, and any occured error err.
//
// First byte for total length of output(including the first byte and the str bytes), after first byte, the rest bytes represent the str itself.
//
// As for now, the str should not greater than 256 bytes. If str great than 256 bytes, the first 256 bytes and a "overflow" error will return.
func StringToBinary(str string) (data []byte, n uint, err error) {
	var raw, lng []byte
	raw = []byte(str)
	if len(raw) > 256 {
		err = errors.New("Overflow: " + str)
		raw = raw[:256]
	}
	n = uint(len(raw))
	binary.BigEndian.PutUint64(lng, uint64(len(raw)))
	data = append(lng, raw...)
	return data, n, err
}

// BinaryToString is the reverse operation of StringToBinary.
//
// It will read a binary flow and convert it into a string str, report total length(including first length byte) of this string, err will always be nil.
func BinaryToString(data []byte) (str string, n uint, err error) {
	n = uint(binary.BigEndian.Uint64(data[:8]))
	str = string(data[8 : 8*n])
	return str, n, nil
}

