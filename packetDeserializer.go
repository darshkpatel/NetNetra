package main

// deserializePacket Will contain parameters to index on elasticsearch
type deserializedPacket struct {
	FlowHash      uint64
	NetFlow       string
	Payload       string
	IPv4Length    uint32
	TCPpacketType string
	SrcPort       uint32
	DstPort       uint32
	DstMAC        string
	SrcMAC        string
	ttl           uint32
	srcIP         string
	dstIP         string
	dnsType       string
	dnsRespCode   uint64
	dnsContent    string
}

// takes DecodedPacket and returns a deserializedPacket with the necessary fields to index on elasticsearch
func deserializePacket(decodedpacket DecodedPacket) deserializedPacket {
	deserializedpacket := deserializedPacket{}

	return deserializedpacket
}
