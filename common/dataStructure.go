package common

import (
	"errors"
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
	Temp  float64
	Cycle uint32
	Count uint32
}

// String implements fmt.Stringer interface
func (tp TempPair) String() string {
	sr, err := StructProbe(tp, ":", "\t")
	HandleErr(err)
	return sr.String()
}

// NetworkSettings represents a basic network settings
type NetworkSettings struct {
	InterfaceName string
	ListenPort    int
	Token         string
}

// String implements fmt.Stringer interface
func (ns NetworkSettings) String() string {
	sr, err := StructProbe(ns, ":", "\t")
	HandleErr(err)
	return sr.String()
}

// Config contain all information that define a PWM fan's user configuration
type Config struct {
	Pin             uint8
	CPUTempPath     string
	SampleRate      uint32
	PwmFreq         int
	FullCycle       uint32
	StopCycle       uint32
	Start           TempPair
	High            TempPair
	Low             TempPair
	NetworkSettings NetworkSettings
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
	Trigger Trigger
	Cfg     Config
	UDPAddr *net.UDPAddr
}

func (fan Fan) String() string {
	/*
		sr, err := StructProbe(fan, ":", "\t")
		HandleErr(err)
		return sr.String()
	*/
	str1, _ := ValueToString(fan.Pin)
	str1 = "Pin: " + str1
	str2, _ := ValueToString(fan.Current.Temp)
	str2 = "Temp: " + str2
	str3, _ := ValueToString(fan.Trigger.State)
	str3 = "State: " + str3
	str4, _ := ValueToString(fan.Current.Cycle)
	str4 = "Cycle: " + str4
	str5, _ := ValueToString(fan.Current.Count)
	str5 = "CrCnt: " + str5
	str6, _ := ValueToString(fan.Trigger.Count)
	str6 = "TgCnt: " + str6
	var lng int
	for _, item := range []*string{&str1, &str2, &str3, &str4, &str5, &str6} {
		lng = len(*item)
		*item += strings.Repeat(" ", 16-lng)
	}
	return str1 + str2 + str3 + str4 + str5 + str6
}

// Trigger represents a Fan's state and state's counter
type Trigger struct {
	State State
	Count uint32
}

func (tg Trigger) String() string {
	sr, err := StructProbe(tg, ":", "\t")
	HandleErr(err)
	return sr.String()
}

// RemapFunc is a function that read at least one float64 input and output a list of float64 data
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
	case int:
		str = strconv.FormatInt(int64(v), 10)
	case uint8:
		str = strconv.FormatUint(uint64(v), 10)
	case uint32:
		str = strconv.FormatUint(uint64(v), 10)
	case float64:
		str = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		str = v
	case State:
		str = v.String()
	case *net.UDPAddr:
		str = (*v).String()
	default:
		err = errors.New("Can't convert type " + reflect.TypeOf(value).String() + " to string")
		return str, err
	}
	return str, err
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
		str += strings.Repeat(rs.Delimeter, indentCount)
		if i == 0 {
			str += rs.Type.String() + "{" + "\n"
			indentCount++
			str += strings.Repeat(rs.Delimeter, indentCount)
		}
		str += fp.String()
		str += "\n"
		if i == lng-1 {
			indentCount--
			str += "}" + "\n"
		}
	}
	return str
}
