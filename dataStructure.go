package pwmfan

import (
	"encoding/binary"
	"errors"
	"math"
	"net"
	"reflect"
	"strconv"
)

const (
	uint8Capacity = 1 << 8
)

// TempPair represent a temperature and correspond Cycle and Count.
//
// TempPair struct can be static, as a configuration. Also can be dynamic, as a representation of current fan state.
type TempPair struct {
	Temp  float32
	Cycle uint16
	Count uint16
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
		tmp      float32
		cyl, cnt uint16
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
	ListenPort    uint16
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
		lpt      uint16
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
	Pin         uint8
	CPUTempPath string
	PwmFreq     uint16
	FullCycle   uint16
	SampleRate  uint16
	Start       TempPair
	Stop        TempPair
	High        TempPair
	NetworkSettings
}

// State represent Fan's current running state
//
// TODO: try to implement String() method
type State uint8

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
	StopCounter  uint16
	StartCounter uint16
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

// ConvertToString will read a uint16/float64/string type value, convert them to string.
// Return result string str and its length n.
// If value's type is not any above, will return a "Can't convert error".
func ConvertToString(value interface{}) (str string, n uint, err error) {
	switch v := value.(type) {
	case uint16:
		str = strconv.FormatUint(uint64(v), 10)
	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case string:
		str = v
	default:
		err = errors.New("Can't convert type " + reflect.TypeOf(value).String() + " to string")
		return str, n, err
	}
	n = uint(len(str))
	return str, n, err
}

// ConvertToBinary will read a uint16/float64/string type value, convert them to BigEndian binary flow. With following rules:
//
// Binary flow always begins with one byte "total bytes length", including the first byte and any data length.
//
// The rest of length defined below:
//
// uint16: convert to 2 bytes data,
//
// float32: convert to 4 bytes data,
//
// string: convert dynamic length bytes data, if string type value byte length great than 256, will retrun a "Overflow error",
//
// If value's type is not any above, will return a "Can't convert error".
func ConvertToBinary(value interface{}) (data []byte, n uint, err error) {
	var (
		raw, lng []byte
	)
	switch v := value.(type) {
	case uint16:
		binary.BigEndian.PutUint16(raw, v)
	case float32:
		binary.BigEndian.PutUint32(raw, math.Float32bits(v))
	case string:
		raw = []byte(v)
		if len(raw) > uint8Capacity {
			err = errors.New("Overflow: " + v)
			raw = raw[:uint8Capacity]
		}
	default:
		err = errors.New("Can't convert type " + reflect.TypeOf(value).String() + "to binary flow")
		return data, n, err
	}
	n = uint(len(raw) + 1)
	lng = []byte{byte(n)}
	data = append(lng, raw...)
	return data, n, err
}

// ConvertFromBinary need a uint16/float32/string type pointer, and exact []byte type data, convert data the put it into value.
func ConvertFromBinary(value interface{}, data []byte) (err error) {
	switch v := value.(type) {
	case *uint16:
		value = binary.BigEndian.Uint16(data)
	case *float32:
		value = math.Float32frombits(binary.BigEndian.Uint32(data))
	case *string:
		value = string(data)
	default:
		err = errors.New("Unsupport type: " + reflect.TypeOf(value).String())
	}
	return err
}
