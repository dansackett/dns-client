package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const maxHeaderSize = 12

// Header -- The header contains the following fields:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      ID                       |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    QDCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ANCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    NSCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ARCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Header struct {
	// A 16 bit identifier assigned by the program that generates any kind of
	// query. This identifier is copied the corresponding reply and can be used
	// by the requester to match up replies to outstanding queries.
	ID uint16

	// A one bit field that specifies whether this message is a query (0), or a
	// response (1).
	QR QRType

	// A four bit field that specifies kind of query in this message. This
	// value is set by the originator of a query and copied into the response.
	OPCODE Opcode

	// Authoritative Answer - this bit is valid in responses, and specifies
	// that the responding name server is an authority for the domain name in
	// question section.
	//
	// Note that the contents of the answer section may have multiple owner
	// names because of aliases.  The AA bit corresponds to the name which
	// matches the query name, or the first owner name in the answer section.
	AA byte

	// TrunCation - specifies that this message was truncated due to length
	// greater than that permitted on the transmission channel.
	TC byte

	// Recursion Desired - this bit may be set in a query and is copied into
	// the response. If RD is set, it directs the name server to pursue the
	// query recursively. Recursive query support is optional.
	RD byte

	// Recursion Available - this bit is set or cleared in a response, and
	// denotes whether recursive query support is available in the name server.
	RA byte

	// Reserved for future use. Must be zero in all queries and responses.
	Z byte

	// Response code - this 4 bit field is set as part of responses.
	RCODE ResponseCode

	// an unsigned 16 bit integer specifying the number of entries in the
	// question section.
	QDCOUNT uint16

	// an unsigned 16 bit integer specifying the number of resource records in
	// the answer section.
	ANCOUNT uint16

	// an unsigned 16 bit integer specifying the number of name server resource
	// records in the authority records section.
	NSCOUNT uint16

	// an unsigned 16 bit integer specifying the number of resource records in
	// the additional records section.
	ARCOUNT uint16
}

// Encode formats a header accetable to be sent to a DNS server
func (h Header) Encode() ([]byte, error) {
	var buf bytes.Buffer

	// ID does not need bit shifting
	binary.Write(&buf, binary.BigEndian, h.ID)

	// Create the third octet of the header
	resBytes1 := setBitsAtIdx(byte(h.QR), 0, 1)
	resBytes1 |= setBitsAtIdx(byte(h.OPCODE), 1, 4)
	resBytes1 |= setBitsAtIdx(h.AA, 5, 1)
	resBytes1 |= setBitsAtIdx(h.TC, 6, 1)
	resBytes1 |= setBitsAtIdx(h.RD, 7, 1)

	buf.WriteByte(resBytes1)

	// Create the fourth octet of the header
	resBytes2 := setBitsAtIdx(h.RA, 0, 1)
	// Z should always be 0 for future use
	resBytes2 |= setBitsAtIdx(h.Z, 1, 3)
	resBytes2 |= setBitsAtIdx(byte(h.RCODE), 4, 4)

	buf.WriteByte(resBytes2)

	// Create the final pieces of the header
	binary.Write(&buf, binary.BigEndian, h.QDCOUNT)
	binary.Write(&buf, binary.BigEndian, h.ANCOUNT)
	binary.Write(&buf, binary.BigEndian, h.NSCOUNT)
	binary.Write(&buf, binary.BigEndian, h.ARCOUNT)

	return buf.Bytes(), nil
}

// DecodeHeader translates a bytes buffer into a Header instance understandable
// by a client
func DecodeHeader(data []byte, bytesRead int, h *Header) (int, error) {
	if h == nil {
		return 0, errors.New("Cannot decode bytes to nil Header")
	}

	if len(data) < maxHeaderSize {
		return 0, fmt.Errorf("Header bytes should be %d bytes, found %d", maxHeaderSize, len(data))
	}

	var err error

	h.ID, bytesRead, err = decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	currentByte := data[bytesRead]
	h.QR = QRType(getBitsAtIdx(currentByte, 0, 1))
	h.OPCODE = Opcode(getBitsAtIdx(currentByte, 1, 4))
	h.AA = getBitsAtIdx(currentByte, 5, 1)
	h.TC = getBitsAtIdx(currentByte, 6, 1)
	h.RD = getBitsAtIdx(currentByte, 7, 1)
	bytesRead++

	currentByte = data[bytesRead]
	h.RA = getBitsAtIdx(currentByte, 0, 1)
	h.Z = getBitsAtIdx(currentByte, 1, 3)
	h.RCODE = ResponseCode(getBitsAtIdx(currentByte, 4, 4))
	bytesRead++

	// Set the remaining data
	h.QDCOUNT, bytesRead, err = decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	h.ANCOUNT, bytesRead, err = decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	h.NSCOUNT, bytesRead, err = decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	h.ARCOUNT, bytesRead, err = decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, err
}
