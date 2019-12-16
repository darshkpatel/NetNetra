package main

import (
	"fmt"

	"github.com/darshkpatel/NetNetra/model"
)

// takes DecodedPacket and returns a deserializedPacket with the necessary fields to index on elasticsearch
func deserializePacket(decodedpacket *DecodedPacket) model.DeserializedPacket {

	deserializedpacket := model.DeserializedPacket{
		NetFlow:        decodedpacket.ipv4.SrcIP.String() + "--->" + decodedpacket.ipv4.DstIP.String(),
		Payload:        decodedpacket.Payload.String(),
		TCPContentType: decodedpacket.TCP.SrcPort.LayerType().String(),
		FlowHash:       fmt.Sprintf("%v", decodedpacket.FlowHash),
		IPv4Length:     decodedpacket.ipv4.Length,
		SrcPort:        uint32(decodedpacket.TCP.SrcPort),
		DstPort:        uint32(decodedpacket.TCP.DstPort),
		SrcMAC:         decodedpacket.eth.SrcMAC.String(),
		DstMAC:         decodedpacket.eth.DstMAC.String(),
		TTL:            decodedpacket.ipv4.TTL,
		SrcIP:          decodedpacket.ipv4.SrcIP.String(),
		DstIP:          decodedpacket.ipv4.DstIP.String(),
		DnsRespCode:    decodedpacket.DNS.ResponseCode.String(),
		DnsContent:     string(decodedpacket.DNS.Payload()),
	}
	return deserializedpacket
}
