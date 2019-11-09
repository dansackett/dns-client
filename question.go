package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// based on a domain name minimum number of labels
const minimumQnameLen = 2

// Question -- The question section is used to carry the "question" in most
// queries, i.e., the parameters that define what is being asked.  The section
// contains QDCOUNT (usually 1) entries, each of the following format:
//
//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                                               |
//     /                     QNAME                     /
//     /                                               /
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QTYPE                     |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                     QCLASS                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
type Question struct {
	// a domain name represented as a sequence of labels, where each label
	// consists of a length octet followed by that number of octets. The
	// domain name terminates with the zero length octet for the null label of
	// the root. Note that this field may be an odd number of octets; no
	// padding is used.
	QNAME string

	// a two octet code which specifies the type of the query. The values for
	// this field include all codes valid for a TYPE field, together with some
	// more general codes which can match more than one type of RR.
	QTYPE RecordType

	// a two octet code that specifies the class of the query. For example,
	// the QCLASS field is IN for the Internet.
	QCLASS RecordClass
}

// Encode translates a Question into a byte array suitable for a DNS server
func (q Question) Encode() ([]byte, error) {
	var err error
	var buf bytes.Buffer

	labels := strings.Split(q.QNAME, ".")

	if len(labels) < minimumQnameLen {
		return buf.Bytes(), errors.New("Malformed QName field")
	}

	for _, label := range labels {
		if len(label) == 0 {
			return buf.Bytes(), errors.New("Malformed label found, must not be 0 length")
		}

		buf.WriteByte(uint8(len(label)))
		buf.Write([]byte(label))
	}

	buf.WriteByte(0x00)

	binary.Write(&buf, binary.BigEndian, q.QTYPE)
	binary.Write(&buf, binary.BigEndian, q.QCLASS)

	return buf.Bytes(), err
}

// DecodeQuestion translates a byte slice to a Question object
func DecodeQuestion(data []byte, bytesRead int, q *Question) (int, error) {
	var err error

	if q == nil {
		return 0, errors.New("Cannot decode bytes to nil Question")
	}

	q.QNAME, bytesRead, err = getPrintableDomainStr(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	typ, bytesRead, err := decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	q.QTYPE = RecordType(typ)

	cls, bytesRead, err := decodeUint16(data, bytesRead)
	if err != nil {
		return bytesRead, err
	}

	q.QCLASS = RecordClass(cls)

	return bytesRead, err
}

func (q *Question) String() string {
	name := q.QNAME
	class := RecordClassToStrMap[q.QCLASS]
	typ := RecordTypeToStrMap[q.QTYPE]

	return fmt.Sprintf("%s\t\t\t%s\t%s", name, class, typ)
}
