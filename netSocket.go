package pwmfan

import (
	"net"
)

// ResolveUDPAddr resolve the listening udp address from fan config
func (fan *Fan) ResolveUDPAddr() (udpAddr *net.UDPAddr) {
	cfg := fan.GetCfg()
	ip := IFNameToIP(cfg.NetworkInterfaceName)
	udpAddr = &net.UDPAddr{
		IP:   ip,
		Port: cfg.ListenPort,
		Zone: "",
	}
	return udpAddr
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
