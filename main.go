package main

import (
	"errors"
	"log"
	"net"

	"github.com/darshkpatel/NetNetra/model"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// Borrowed from examples https://github.com/google/gopacket/blob/0ad7f2610e344e58c1c95e2adda5c3258da8e97b/examples/arpscan/arpscan.go#L58
// From all the interfaces, checks for sane interfaces
func isSane(iface net.Interface) bool {
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return false
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}

	if addr == nil || addr.IP[0] == 127 || addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return false
	}
	return true
}

// Currently returns first sane interface's name
// Later on, probably return a list of names for the user to choose from in case of multiple sane interfaces found
func getInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Panic(err)
	}
	for _, iface := range ifaces {
		if isSane(iface) {
			return iface.Name, nil
		}
	}
	return "", errors.New("no interface found")

}

func main() {

	_, err := model.Connect()
	if err != nil {
		log.Panic(err)
	}

	iface, err := getInterface()
	if err != nil {
		log.Panic(err)
	}
	handle, err := pcap.OpenLive(iface, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Panic(err)
	}
	defer handle.Close()

	src := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range src.Packets() {
		go handlePacket(packet)
	}
}
