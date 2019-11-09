package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var domainFlagVal = flag.String("domain", "", "The domain to run DNS queries on. This is required.")
var recordTypeFlagVal = flag.String("type", "A", "Record type to lookup. Defaults to \"A\"")
var dnsServerAddrFlagVal = flag.String("server-addr", "8.8.8.8:53", "IP and Port for the DNS server to query. Defaults to \"8.8.8.8:53\".")

func main() {
	//-------------------------------------------------------------------------
	// 1. Initialize client and parse flags
	//-------------------------------------------------------------------------
	client, err := InitClient()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	flag.Parse()

	// Validate flags
	if *domainFlagVal == "" {
		log.Fatalf("error: %v", "'domain' is required")
	}

	if _, ok := RecordTypeStrToRecordTypeMap[*recordTypeFlagVal]; !ok {
		log.Fatalf("error: Type '%s' not implemented\n", *recordTypeFlagVal)
	}

	//-------------------------------------------------------------------------
	// 2. Create message for the DNS server
	//-------------------------------------------------------------------------
	var recursionDesired byte = 1

	questions := []Question{
		Question{
			QNAME:  *domainFlagVal,
			QTYPE:  RecordTypeStrToRecordTypeMap[*recordTypeFlagVal],
			QCLASS: RecordClassIN,
		},
	}

	header := Header{
		ID:      GenerateRandID(),
		QR:      QRTypeQuery,
		OPCODE:  OpcodeQuery,
		QDCOUNT: uint16(len(questions)),
		RD:      recursionDesired,
	}

	m := &Message{
		Header:    header,
		Questions: questions,
	}

	//-------------------------------------------------------------------------
	// 3. Encode the message bytes and write them to the server
	//-------------------------------------------------------------------------
	msgBytes, err := m.Encode()

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	startQueryTime := time.Now()
	_, err = client.conn.Write(msgBytes)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//-------------------------------------------------------------------------
	// 4. Listen for a response from the server and write it to a buffer
	//-------------------------------------------------------------------------
	_, err = client.conn.Read(client.respBuf)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	elapsedQueryTime := time.Since(startQueryTime)

	//-------------------------------------------------------------------------
	// 5. Decode the response into a Message object for parsing
	//-------------------------------------------------------------------------
	msg := new(Message)
	_, err = DecodeMessage(client.respBuf, msg, elapsedQueryTime)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//-------------------------------------------------------------------------
	// 6. Print the parsed response Message object into a dig-esque output
	//-------------------------------------------------------------------------
	fmt.Println(msg)
}
