package main

import (
	"github.com/google/gopacket"
)

// Struct to be indexed in ElasticSearch
type deserializedFlow struct {
	Src          gopacket.Endpoint
	Dest         gopacket.Endpoint
	FlowHash     []byte
	EndpointType gopacket.EndpointType
}

// Utility function which returns deserializedFlow object from the input packet
func flowDeserializer(packet gopacket.Packet) deserializedFlow {
	netFlow := packet.NetworkLayer().NetworkFlow()
	src, dst := netFlow.Endpoints()
	endType := netFlow.EndpointType()
	flowHash := computeFlowHash(netFlow)
	sFlow := deserializedFlow{src, dst, flowHash, endType}
	return sFlow
}
