package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// DecodedPacket holds necessary information extracted from packet
type DecodedPacket struct {
	decodedLayers []gopacket.LayerType

	TCP layers.TCP
	UDP layers.UDP
	ARP layers.ARP
	DNS layers.DNS

	eth  layers.Ethernet
	ipv4 layers.IPv4
	ipv6 layers.IPv6

	Payload  gopacket.Payload
	FlowHash uint64 //Can be used to track Flows/Sessions in future
	NetFlow  gopacket.Flow

	Parser *gopacket.DecodingLayerParser
}

func createDecodedPacket() *DecodedPacket {
	decodedPacket := &DecodedPacket{}

	decodedPacket.Parser = gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.TCP)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.UDP)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.ARP)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.DNS)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.eth)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.ipv4)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.ipv6)
	decodedPacket.Parser.AddDecodingLayer(&decodedPacket.Payload)
	return decodedPacket
}

// This function should decode packets and create a DecodedPacket Object with necessary decoded layers
func handlePacket(packet gopacket.Packet) {
	decodedPacket := createDecodedPacket()

	err := decodedPacket.Parser.DecodeLayers(packet.Data(), &decodedPacket.decodedLayers)
	if err != nil {
		log.Printf("Could not decode layer: %v\n", err)
	}

	for _, layerType := range decodedPacket.decodedLayers {
		switch layerType {
		case layers.LayerTypeIPv6:
			decodedPacket.NetFlow = decodedPacket.ipv6.NetworkFlow()
			log.Println("IPv6: ", decodedPacket.ipv6.SrcIP, "->", decodedPacket.ipv6.DstIP)

		case layers.LayerTypeIPv4:
			decodedPacket.NetFlow = decodedPacket.ipv4.NetworkFlow()
			log.Println("IPv4: ", decodedPacket.ipv4.SrcIP, "->", decodedPacket.ipv4.DstIP)

		case layers.LayerTypeTCP:
			decodedPacket.FlowHash = hashCombine(decodedPacket.NetFlow.FastHash(), decodedPacket.TCP.TransportFlow().FastHash())
		}
	}
}

// based on boost::hash_combine
// http://www.boost.org/doc/libs/1_63_0/boost/functional/hash/hash.hpp
func hashCombine(h, k uint64) uint64 {
	m := uint64(0xc6a4a7935bd1e995)
	r := uint64(47)

	k *= m
	k ^= k >> r
	k *= m

	h ^= k
	h *= m

	return h + 0xe6546b64
}
