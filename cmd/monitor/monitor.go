package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eZioPan/pwmfan-go/common"
)

// Config represent a client config
type Config struct {
	NetworkInterface string
	RemoteHost       string
	RemotePort       int
	Token            string
	SampleRate       float64
	Message          []string
}

var (
	configPath = ""
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "monitor configuration file path")
}

func main() {
	flag.Parse()
	cfg := &Config{}
	common.ParseJSON(configPath, cfg)
	lIP := common.IFNameToIPv4(cfg.GetNetworkInterfaceName())
	lUDPAddr, err := net.ResolveUDPAddr("udp", lIP.String()+":0")
	common.HandleErr(err)
	rAddr := cfg.GetRemoteHost() + ":" + strconv.Itoa(int(cfg.GetRemotePort()))
	rUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	common.HandleErr(err)
	conn, err := net.DialUDP("udp", lUDPAddr, rUDPAddr)
	common.HandleErr(err)
	msg := make([]byte, 1024)
	var lastMsgLen int
	for {
		_, err = conn.Write([]byte(cfg.GetToken()))
		common.HandleErr(err)
		lng, err := conn.Read(msg)
		common.HandleErr(err)
		cls := strings.Repeat(" ", lastMsgLen)
		lastMsgLen = lng
		os.Stdout.Write(append([]byte("\r"), []byte(cls)...))
		os.Stdout.Write(append([]byte("\r"), msg[:lng]...))
		time.Sleep(time.Second / time.Duration(cfg.GetSampleRate()))
	}
}

// GetRemoteHost get remote host from config
func (cfg Config) GetRemoteHost() (remoteHost string) {
	return cfg.RemoteHost
}

// SetRemoteHost set remote host to config
func (cfg *Config) SetRemoteHost(remoteHost string) {
	cfg.RemoteHost = remoteHost
}

// GetRemotePort get remote listening port from config
func (cfg Config) GetRemotePort() (remotePort int) {
	return cfg.RemotePort
}

// SetRemotePort set remote listening host to config
func (cfg *Config) SetRemotePort(remotePort int) {
	cfg.RemotePort = remotePort
}

// GetNetworkInterfaceName get local network interface name from config
func (cfg Config) GetNetworkInterfaceName() (ifn string) {
	return cfg.NetworkInterface
}

// SetNetworkInterfaceName set local network interface name to config
func (cfg *Config) SetNetworkInterfaceName(ifn string) {
	cfg.NetworkInterface = ifn
}

// GetToken get token from config
func (cfg Config) GetToken() (token string) {
	return cfg.Token
}

// SetToken set token to config
func (cfg *Config) SetToken(token string) {
	cfg.Token = token
}

// GetSampleRate get sample rate from config
func (cfg Config) GetSampleRate() (sampleRate float64) {
	return cfg.SampleRate
}

// SetSampleRate set sample rate to config
func (cfg *Config) SetSampleRate(sampleRate float64) {
	cfg.SampleRate = sampleRate
}

// GetMessage get message value
func (cfg Config) GetMessage() (message []string) {
	return cfg.Message
}

// SetMessage set message value
func (cfg *Config) SetMessage(message []string) {
	cfg.Message = message
}

// ChkSum caculate md5 ChkSum of network port and token
func (cfg Config) ChkSum() [md5.Size]byte {
	b := make([]byte, 0, 512)
	binary.PutVarint(b, int64(cfg.GetRemotePort()))
	b = append(b, []byte(cfg.GetToken())...)
	return md5.Sum([]byte(cfg.GetToken()))
}

// EncodingMessage encoding request message
func (cfg Config) EncodingMessage() (rawStream []byte, err error) {
	rawStream = make([]byte, 0, 512)
	buf := bytes.NewBuffer(rawStream)
	ge := gob.NewEncoder(buf)
	err = ge.Encode(cfg.ChkSum())
	if err != nil {
		return nil, err
	}
	err = ge.Encode(cfg.Message)
	if err != nil {
		return nil, err
	}
	return nil, nil
}


