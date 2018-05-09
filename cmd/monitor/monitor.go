package main

import (
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
	RemotePort       uint
	Token            string
	SampleRate       float64
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
func (cfg Config) GetRemotePort() (remotePort uint) {
	return cfg.RemotePort
}

// SetRemotePort set remote listening host to config
func (cfg *Config) SetRemotePort(remotePort uint) {
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
