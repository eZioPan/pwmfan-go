//TODO: move all network relative code into a new repo

package pwmfan

import (
	"errors"
	"net"
)

// ResolveUDPAddr resolve the listening udp address from fan config
func (fan *Fan) ResolveUDPAddr() {
	cfg := fan.GetCfg()
	ip := IFNameToIPv4(cfg.NetworkInterfaceName)
	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: cfg.ListenPort,
		Zone: "",
	}
	fan.SetUDPAddr(udpAddr)
}

// IFNameToIPv4 reads an network interface name, return an net.IP, with the first address of the network interface
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

// GetUDPAddr get fan's udp listening address
func (fan Fan) GetUDPAddr() (UDPAddr *net.UDPAddr) {
	return fan.UDPAddr
}

// SetUDPAddr set fan's udp listening address
func (fan *Fan) SetUDPAddr(udpAddr *net.UDPAddr) {
	fan.UDPAddr = udpAddr
}

// CreateServer use fan's configuration, return a *net.UDPConn object
func (fan *Fan) CreateServer() (udpConn *net.UDPConn) {
	fan.ResolveUDPAddr()
	udpConn, err := net.ListenUDP("udp", fan.GetUDPAddr())
	HandleErr(err)
	return udpConn
}

// HandleRequest handle request from network
func (fan *Fan) HandleRequest(udpConn *net.UDPConn) {
	//use *Fan for up-to-date data
	msg := make([]byte, 16)
	for {
		cnt, rAddr, err := udpConn.ReadFromUDP(msg)
		HandleErr(err)
		if string(msg[:cnt]) == fan.GetCfg().Token {
			udpConn.WriteToUDP([]byte(fan.String()), rAddr)
		}
	}
}
