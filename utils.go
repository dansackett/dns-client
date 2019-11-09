package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

const (
	// mostly so algorithms can make sense when read
	octetMaxIdx = 7

	// See RFC 1035 section 2.3.4
	maxDomainNameWireOctets = 255

	// This is the maximum number of compression pointers that should occur in a
	// semantically valid message. Each label in a domain name must be at least one
	// octet and is separated by a period. The root label won't be represented by a
	// compression pointer to a compression pointer, hence the -2 to exclude the
	// smallest valid root label.
	//
	// It is possible to construct a valid message that has more compression pointers
	// than this, and still doesn't loop, by pointing to a previous pointer. This is
	// not something a well written implementation should ever do, so we leave them
	// to trip the maximum compression pointer check.
	// borrowed from: https://github.com/miekg/dns/blob/master/msg.go
	maxCompressionPointers = (maxDomainNameWireOctets+1)/2 - 2
)

// makeOctetMask creates a bitmask. This bitmask can be used to extract a value
// from an octet by ANDing the bytes. For example, say we wanted to read just
// the first 4 bits of the octet. We would want to unset the other bits so only
// those 4 are set. We create a mask to do that.
//
// mask = (1 << 4) - 1 = 00001111
// val = 10011101
// val & mask = 00001101
func makeOctetMask(size uint) byte {
	return (1 << size) - 1
}

// In order to set the correct fields for the header, we want to use bit
// shifting in the octet. For example:
//
// Setting AA in the 3rd octet in the header:
// AA = 00000001 (is authoritative)
// octetMaxIdx = 7
// bitIdx = 5
// size = 1
// 00000001 << 7 - (5 + (1 - 1)) === 00000001 << 2 === 00000100
//
// This would set the AA field based on the RFC since the octet scheme is:
///
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
func setBitsAtIdx(b byte, bitIdx, size uint) byte {
	return b << (octetMaxIdx - (bitIdx + (size - 1)))
}

// getBitAtIdx is similiar to setting bits except we want to shift them to the
// right and then use a bit mask to unset all of the additional values to the
// left of our bits. For example, let's get the 3 bit sequence starting at the
// fourth bit in the octet.
//
// b = 10010010
// b = 100 | 100 | 10
// bitIdx = 3
// size = 3
// b >> (7 - (3 + (3 - 1))) === b >> (7 - 5) === 00100100
// 00100100 & 00000111 === 00000100
func getBitsAtIdx(b byte, bitIdx, size uint) byte {
	return (b >> (octetMaxIdx - (bitIdx + (size - 1)))) & makeOctetMask(size)
}

// decodeUint16 helps in unpacking a BigEndian Uint16 value based on the
// current offset of bytes read.
func decodeUint16(data []byte, offset int) (uint16, int, error) {
	if offset+2 > len(data) {
		return 0, len(data), errors.New("Error unpacking Uint16: Overflow")
	}
	return binary.BigEndian.Uint16(data[offset:]), offset + 2, nil
}

// decodeUint32 helps in unpacking a BigEndian Uint32 value based on the
// current offset of bytes read.
func decodeUint32(data []byte, offset int) (uint32, int, error) {
	if offset+4 > len(data) {
		return 0, len(data), errors.New("Error unpacking Uint32: Overflow")
	}
	return binary.BigEndian.Uint32(data[offset:]), offset + 4, nil
}

// extractDomainNameLabels parses data based on how domains names are stored in
// DNS messages. It takes into account name compression and recursively returns
// a slice of labels for the domain name.
func extractDomainNameLabels(data []byte, bytesRead int) ([]string, int, error) {
	var labels []string

	ptrsFollowed := 0
	domainSpaceLeft := maxDomainNameWireOctets

	for {
		if ptrsFollowed >= maxCompressionPointers {
			return labels, bytesRead, errors.New("Too many compression pointers in domain name")
		}

		currentByte := data[bytesRead]

		switch currentByte & 0xC0 {

		// we have a pointer
		case 0xC0:
			offset := (data[bytesRead] & makeOctetMask(6)) | data[bytesRead+1]
			lbls, _, err := extractDomainNameLabels(data, int(offset))
			ptrsFollowed++

			if err != nil {
				return labels, bytesRead, err
			}

			labels = append(labels, lbls...)
			bytesRead += 2
			return labels, bytesRead, nil

		// we have a label
		case 0x00:
			labelLen := int(currentByte)
			label := string(data[bytesRead+1 : bytesRead+labelLen+1])
			labels = append(labels, label)
			bytesRead += labelLen + 1

			// +1 for the label separator
			domainSpaceLeft -= labelLen + 1
			if domainSpaceLeft <= 0 {
				return labels, bytesRead, errors.New("Domains exceed max size for field")
			}

			if nextByte := data[bytesRead]; nextByte == 0x00 {
				bytesRead++
				return labels, bytesRead, nil
			}

		default:
			return labels, bytesRead, errors.New("Invalid RData found: could not parse labels for domain")
		}
	}
}

// Get the printable string for a domain record
func getPrintableDomainStr(data []byte, offset int) (string, int, error) {
	var err error

	labels, bytesRead, err := extractDomainNameLabels(data, offset)

	if err != nil {
		return "", bytesRead, err
	}

	return fmt.Sprintf("%s.", strings.Join(labels, ".")), bytesRead, err
}
