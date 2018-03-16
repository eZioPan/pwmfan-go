package pwmfan

import (
	"net"
)

// ResolveUDPAddr resolve the listening udp address from fan config
func (fan *Fan) ResolveUDPAddr() {
	cfg := fan.GetCfg()
	ip := IFNameToIP(cfg.NetworkInterfaceName)
	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: cfg.ListenPort,
		Zone: "",
	}
	fan.SetUDPAddr(udpAddr)
}

// IFNameToIP reads an network interface name, return an net.IP, with the first address of the network interface
func IFNameToIP(ifname string) (ip net.IP) {
	iface, err := net.InterfaceByName(ifname)
	HandleErr(err)
	addrs, err := iface.Addrs()
	HandleErr(err)
	ip = addrs[0].(*net.IPAddr).IP
	return ip
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
func (fan Fan) HandleRequest(udpConn *net.UDPConn) {
	msg := make([]byte, 16)
	for {
		cnt, err := udpConn.Read(msg)
		HandleErr(err)
		if string(msg[:cnt]) == fan.GetCfg().Token {
			udpConn.Write([]byte(fan.String()))
		}
	}
}
