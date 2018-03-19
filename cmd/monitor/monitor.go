package main

import (
	"encoding/json"
	"errors"
	"flag"
	"net"
	"os"
	"strconv"
	"time"
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
	flag.StringVar(&configPath, "config", "config.json", "system temperature file")
}

func main() {
	flag.Parse()
	cfg := ParseJSON(configPath)
	lIP := IFNameToIPv4(cfg.GetNetworkInterfaceName())
	lUDPAddr, err := net.ResolveUDPAddr("udp", lIP.String()+":0")
	HandleErr(err)
	rAddr := cfg.GetRemoteHost() + ":" + strconv.Itoa(int(cfg.GetRemotePort()))
	rUDPAddr, err := net.ResolveUDPAddr("udp", rAddr)
	HandleErr(err)
	conn, err := net.DialUDP("udp", lUDPAddr, rUDPAddr)
	HandleErr(err)
	msg := make([]byte, 64)
	for {
		_, err = conn.Write([]byte(cfg.GetToken()))
		HandleErr(err)
		len, err := conn.Read(msg)
		HandleErr(err)
		os.Stdout.Write(msg[:len])
		time.Sleep(time.Second / time.Duration(cfg.GetSampleRate()))
	}
}

// HandleErr panic none nil error
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

//IFNameToIPv4 read a network interface name and return a net.IP
func IFNameToIPv4(ifname string) (ip net.IP) {
	iface, err := net.InterfaceByName(ifname)
	HandleErr(err)
	addrs, err := iface.Addrs()
	HandleErr(err)
	for _, addr := range addrs {
		ip = addr.(*net.IPNet).IP
		if len(ip.DefaultMask()) == net.IPv4len {
			return ip
		}
	}
	HandleErr(errors.New("No IPv4 address"))
	return nil
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
