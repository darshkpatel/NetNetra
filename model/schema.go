package model

import (
	"context"
	"errors"
)

type DeserializedPacket struct {
	FlowHash      uint64
	NetFlow       string
	Payload       string
	IPv4Length    uint16
	TCPpacketType string
	SrcPort       uint32
	DstPort       uint32
	DstMAC        string
	SrcMAC        string
	ttl           uint8
	srcIP         string
	dstIP         string
	dnsType       string
	dnsRespCode   string
	dnsContent    string
}

//TODO Mapping subject to Change
const Mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"Packet":{
			"properties":{
				
				"FlowHash":{
					"type":"long"
				},
				"NetFlow":{
					"type":"text",
				},
				"Payload":{
					"type":"text"
				},
				"IPv4Length":{
					"type":"long"
				},
				"TCPpacketType":{
					"type":"keyword"
				},
				"SrcPort":{
					"type":"keyword"
				},
				"DstPort":{
					"type":"keyword"
				},
				"DstMAC":{
					"type":"text"
				},
				"SrcMAC":{
					"type":"text"
				},
				"ttl":{
					"type":"long"
				},
				"srcIP":{
					"type":"IP"
				},
				"dstIP":{
					 "type":"IP"

				},
				"dnsType":{
					"type":"text"
				},
				"dnsRespCode":{
					"type":"long"
				},
				"dnsContent":{
					"type":"text"
				}
			}
		}
	}	

	
}`

// Index Function indexes the deserialized packet into packet index
func Index(ds deserializedPacket) error {

	ctx := context.Background()
	exists, err := client.IndexExists("Packet").Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		createIndex, err := client.CreateIndex("Packet").BodyString(Mapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("Indexing error into ElasticSearch")
		}
	}
	_, err = client.Index().Index("Packet").BodyJson(ds).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
