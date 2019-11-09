package main

import (
	"errors"
	"fmt"
)

// RR is a Resource Record and is the response given for a DNS Question
//
// All RRs have the same top level format shown below:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                                               /
//     /                      NAME                     /
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     CLASS                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      TTL                      |
//     |                                               |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                   RDLENGTH                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//     /                     RDATA                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type RR struct {
	// an owner name, i.e., the name of the node to which this resource record
	// pertains.
	NAME string

	// two octets containing one of the RR TYPE codes.
	TYPE RecordType

	// two octets containing one of the RR CLASS codes.
	CLASS RecordClass

	// a 32 bit signed integer that specifies the time interval that the
	// resource record may be cached before the source of the information
	// should again be consulted. Zero values are interpreted to mean that the
	// RR can only be used for the transaction in progress, and should not be
	// cached.  For example, SOA records are always distributed with a zero TTL
	// to prohibit caching. Zero values can also be used for extremely
	// volatile data.
	TTL uint32

	// an unsigned 16 bit integer that specifies the length in octets of the
	// RDATA field.
	RDLENGTH uint16

	// a variable length string of octets that describes the resource. The
	// format of this information varies according to the TYPE and CLASS of the
	// resource record.
	RDATA ResourceDataField
}

// Encode translates an RR to a byte slice for sending as a DNS message.
// @TODO implement this method...
func (rr *RR) Encode() ([]byte, error) {
	return []byte{}, errors.New("RR encoding not implemented")
}

func (rr *RR) String() string {
	name := rr.NAME
	ttl := rr.TTL
	class := RecordClassToStrMap[rr.CLASS]
	typ := RecordTypeToStrMap[rr.TYPE]
	rData := rr.RDATA

	return fmt.Sprintf("%s\t\t%d\t%s\t%s\t%s", name, ttl, class, typ, rData)
}

// DecodeRR translates a byte slice to an RR object
func DecodeRR(data []byte, bytesRead int, rr *RR) (int, error) {
	var err error

	if rr == nil {
		return 0, errors.New("Cannot decode bytes to nil RR")
	}

	rr.NAME, bytesRead, err = getPrintableDomainStr(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	typ, bytesRead, err := decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	rr.TYPE = RecordType(typ)

	cls, bytesRead, err := decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	rr.CLASS = RecordClass(cls)

	rr.TTL, bytesRead, err = decodeUint32(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	rr.RDLENGTH, bytesRead, err = decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	rr.RDATA, err = getResourceDataFieldForResourceType(rr.TYPE, data, bytesRead, rr.RDLENGTH)
	bytesRead += int(rr.RDLENGTH)
	if err != nil {
		return bytesRead, err
	}

	return bytesRead, err
}

// getResourceDataFieldForResourceType returns a struct which implements the
// ResourceDataField interface. Each of these types make understanding RR data easier.
func getResourceDataFieldForResourceType(rrType RecordType, data []byte, bytesRead int, dataLen uint16) (ResourceDataField, error) {

	switch rrType {

	case RecordTypeA:
		return NewRDataA(data[bytesRead : uint16(bytesRead)+dataLen])

	case RecordTypeAAAA:
		return NewRDataAAAA(data[bytesRead : uint16(bytesRead)+dataLen])

	case RecordTypeTXT:
		return NewRDataTXT(data[bytesRead+1 : uint16(bytesRead)+dataLen])

	case RecordTypeCNAME:
		return NewRDataCNAME(data, bytesRead)

	case RecordTypeNS:
		return NewRDataNS(data, bytesRead)

	case RecordTypeSOA:
		return NewRDataSOA(data, bytesRead)

	case RecordTypeMX:
		return NewRDataMX(data, bytesRead)

	case RecordTypePTR:
		return NewRDataPTR(data, bytesRead)

	default:
		if IsRecordTypeObsolete(rrType) {
			return NewRDataObsolete()
		}

		if IsRecordTypeNotImplemented(rrType) {
			return NewRDataNotImplemented()
		}

		return NewRDataUnknown(rrType)

	}
}
