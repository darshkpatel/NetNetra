package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"sync"

	"github.com/darshkpatel/NetNetra/model"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var wg sync.WaitGroup

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
	FlowHash []byte //Used to map each packet to a Flow
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
	wg.Add(1)
	defer wg.Done()

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
			flowHash := computeFlowHash(decodedPacket.NetFlow)
			decodedPacket.FlowHash = flowHash
		}

	}
	//TODO Deserialization from decoded Packet and indexing
	d := deserializePacket(decodedPacket)
	Elerr := model.Index(&d)
	//handling errors in stream of packets
	if Elerr != nil {
		fmt.Print(Elerr)
	}
	fmt.Printf("%+v", d)

}

// Since FastHash cannot be used as a key, this function computes hash of Flow object
// which can be used as a key.
// Ref: https://godoc.org/github.com/google/gopacket#Flow
func computeFlowHash(flow gopacket.Flow) []byte {
	h := sha256.New()
	byteTemp := fmt.Sprintf("%v", flow)
	flowHash := h.Sum([]byte(byteTemp))
	return flowHash
}
