package pwmfan

import (
	"encoding/binary"
	"errors"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
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
func (tp TempPair) String() string {
	sr, err := StructProbe(tp, ":", "\t")
	HandleErr(err)
	return sr.String()
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
func (ns NetworkSettings) String() string {
	sr, err := StructProbe(ns, ":", "\t")
	HandleErr(err)
	return sr.String()
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

func (cfg Config) String() string {
	sr, err := StructProbe(cfg, ":", "\t")
	HandleErr(err)
	return sr.String()
}

// State represent Fan's current running state
type State uint8

const (
	// Stop represent Pwm Fan stay in a full stop state
	Stop State = iota
	// Start represent Pwm Fan enter in s Start state
	Start
	// Run represent Pwm Fan stay in a running state
	Run
)

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

// Fan represent a Fan
//
// State, Cycle, Temp should be updated 'in realtime'
// Pin better not be modified after call NewFan()
type Fan struct {
	Pin     uint8
	Current TempPair
	StateRecord
	Cfg     Config
	UDPAddr *net.UDPAddr
}

func (fan Fan) String() string {
	sr, err := StructProbe(fan, ":", "\t")
	HandleErr(err)
	return sr.String()
}

// StateRecord represents a Fan's state and state's counter
type StateRecord struct {
	State
	StopCounter  uint16
	StartCounter uint16
}

func (sRec StateRecord) String() string {
	sr, err := StructProbe(sRec, ":", "\t")
	HandleErr(err)
	return sr.String()
}

// RemapFunc is a function that read at least one float64 input and output a float64 data
//
// More specific,this function should take 1st argument as an input data(s),then read more input, and output result data(s)
//
// Example: LinearRemap and LinearClampRemap
type RemapFunc func([]float64, ...float64) []float64

// ValueToString will read a uint16/float64/string/StructRepresent type value, convert them to string.
// Return result string str and any error.
// If value's type is not any above, will return a "Can't convert error".
func ValueToString(value interface{}) (str string, err error) {
	switch v := value.(type) {
	case StructRepresent:
		str = v.String()
	case uint16:
		str = strconv.FormatUint(uint64(v), 10)
	case float32:
		str = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case string:
		str = v
	default:
		err = errors.New("Can't convert type " + reflect.TypeOf(value).String() + " to string")
		return str, err
	}
	return str, err
}

// ValueToBinary will read a uint16/float64/string type value, convert them to BigEndian binary flow. With following rules:
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
func ValueToBinary(value interface{}) (data []byte, n uint, err error) {
	var (
		raw, lng []byte
	)
	switch v := value.(type) {
	case StructRepresent:
		raw, err = v.MarshalBinary()
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

// ValueFromBinary need a uint16/float32/string type pointer, and exact []byte type data, convert data the put it into value.
func ValueFromBinary(value interface{}, data []byte) (n uint, err error) {
	switch value.(type) {
	case *uint16:
		value = binary.BigEndian.Uint16(data)
	case *float32:
		value = math.Float32frombits(binary.BigEndian.Uint32(data))
	case *string:
		value = string(data)
	default:
		err = errors.New("Unsupport type: " + reflect.TypeOf(value).String())
	}
	return 0, err
}

// FieldPair represents a field name & field value pairs of a struct field.
// Seperator use to seperate Name and Value when string() called
type FieldPair struct {
	Name      string
	Value     interface{}
	Seperator string
}

func (fp FieldPair) String() (str string) {
	valueStr, err := ValueToString(fp.Value)
	// TODO: define a error type for compare and handle error
	HandleErr(err)
	name := fp.Name
	sep := fp.Seperator
	str = name + sep + valueStr
	return str
}

// MarshalBinary implements binary.BinaryMarshaler
func (fp FieldPair) MarshalBinary() (data []byte, err error) {
	data, _, err = ValueToBinary(fp.Value)
	return data, err
}

// StructRepresent store the struct type and all fields' name/value pairs of a struct.
// Delimeter is the string bwtween two FildPairs when string() called.
type StructRepresent struct {
	Type      reflect.Type
	Pair      []FieldPair
	Delimeter string
}

// StructProbe read a struct value, output a StructRepresent of the structure.
// Sep sets FieldPair.Seperator. Deli sets StructRepresent.Delimeter.
// If structure is not a struct value, then it will return an "Not a struct" error.
func StructProbe(structure interface{}, sep string, deli string) (sr StructRepresent, err error) {
	rv := reflect.ValueOf(structure)
	rt := reflect.TypeOf(structure)
	if rt.Kind() != reflect.Struct {
		err = errors.New("Value: " + rt.String() + " is not a struct")
		return sr, err
	}
	nf := rt.NumField()
	sr.Type = rt
	sr.Pair = make([]FieldPair, nf, nf)
	sr.Delimeter = deli
	var fp FieldPair
	for i := 0; i < nf; i++ {
		fp.Name = rt.Field(i).Name
		value := rv.Field(i).Interface()
		if reflect.TypeOf(value).Kind() == reflect.Struct {
			value, _ = StructProbe(value, sep, deli)
		}
		fp.Value = value
		fp.Seperator = sep
		sr.Pair[i] = fp
	}
	return sr, nil
}

// indentCount is a global counter for counting indent table of a nest struct
var indentCount int

func (rs StructRepresent) String() (str string) {
	lng := len(rs.Pair)
	for i, fp := range rs.Pair {
		if i == 0 {
			str += rs.Type.String() + "{" + rs.Delimeter
			indentCount++
		}
		str += strings.Repeat("\t", indentCount) + fp.String()
		if i < lng-1 {
			str += rs.Delimeter
		}
		if i == lng-1 {
			indentCount--
			str += rs.Delimeter + strings.Repeat("\t", indentCount) + "}"
		}
	}
	return str
} 

// MarshalBinary implements binary.BinaryMarshaler
func (rs StructRepresent) MarshalBinary() (data []byte, err error) {
	var (
		curData []byte
		curErr  error
	)
	for i, fp := range rs.Pair {
		curData, curErr = fp.MarshalBinary()
		if err != nil {
			return data, curErr
		}
		data = append(data, curData...)
	}
	return data, err
}

func GetValueFromStruct(structure interface{}, name string, value interface{}) (err error) {
	sr, err := StructProbe(structure, "", "")
	if err != nil {
		return err
	}
	for _, item := range sr.Pair {
		if name == item.Name {
			value = item.Value
			return nil
		}
	}
	err = errors.New("No Field Found")
	return err
}
 
func SetValueToStruct(structure interface{}, name string, value interface{}) (err error) {
	rt := reflect.TypeOf(structure)
	rv := reflect.ValueOf(structure)
	if rt.Kind() != reflect.Struct {
		err = errors.New("No a structure")
		return err
	}
	numField := rt.NumField()
	for i := 0; i < numField; i++ {
		if rt.Field(i).Name == name {
			vrv := reflect.ValueOf(value)
			rv.Field(i).Set(vrv)
			return nil
		}
	}
	err = errors.New("No Field Found")
	return err
}
