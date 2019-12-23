package model

import (
	"context"
	"errors"
)

type DeserializedPacket struct {
	FlowHash       string
	TCPContentType string
	NetFlow        string
	Payload        string
	IPv4Length     uint16
	SrcPort        uint32
	DstPort        uint32
	DstMAC         string
	SrcMAC         string
	TTL            uint8
	SrcIP          string
	DstIP          string
}

//TODO Mapping subject to Change
const Mapping = `{"mappings":{"packet":{"properties":{"TCPContentType":{"type":"keyword"},"FlowHash":{"type":"keyword"},"NetFlow":{"type":"text"},"Payload":{"type":"text"},"IPv4Length":{"type":"long"},"SrcPort":{"type":"keyword"},"DstPort":{"type":"keyword"},"DstMAC":{"type":"text"},"SrcMAC":{"type":"text"},"TTL":{"type":"long"},"SrcIP":{"type":"text"},"DstIP":{"type":"text"}}}}}`

// Index Function indexes the deserialized packet into packet index
func Index(ds *DeserializedPacket) error {

	ctx := context.Background()
	exists, err := client.IndexExists("packets").Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		createIndex, err := client.CreateIndex("packets").BodyString(Mapping).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("Indexing error into ElasticSearch")
		}

	}

	_, err = client.Index().Index("packets").Type("packet").BodyJson(*ds).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
