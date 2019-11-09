package main

import (
	"fmt"
	"net"
)

// ResourceDataField is an interface used to make RData easily readable in an RR
type ResourceDataField interface {
	String() string
}

//-----------------------------------------------------------------------------
// UNKNOWN Record RDATA
//-----------------------------------------------------------------------------

// RDataUnknown represents an A Record
type RDataUnknown struct {
	qType RecordType
}

// NewRDataUnknown creates a new RDataUnknown instance
func NewRDataUnknown(typ RecordType) (*RDataUnknown, error) {
	return &RDataUnknown{qType: typ}, nil
}

// String makes this record printable
func (r *RDataUnknown) String() string {
	return fmt.Sprintf("UNKNOWN RECORD TYPE: %d", r.qType)
}

//-----------------------------------------------------------------------------
// NOT IMPLEMENTED Record RDATA
//-----------------------------------------------------------------------------

// RDataNotImplemented represents an A Record
type RDataNotImplemented struct{}

// NewRDataNotImplemented creates a new RDataNotImplemented instance
func NewRDataNotImplemented() (*RDataNotImplemented, error) {
	return &RDataNotImplemented{}, nil
}

// String makes this record printable
func (r *RDataNotImplemented) String() string {
	return "Not Implemented"
}

//-----------------------------------------------------------------------------
// OBSOLETE Record RDATA
//-----------------------------------------------------------------------------

// RDataObsolete represents an A Record
type RDataObsolete struct{}

// NewRDataObsolete creates a new RDataObsolete instance
func NewRDataObsolete() (*RDataObsolete, error) {
	return &RDataObsolete{}, nil
}

// String makes this record printable
func (r *RDataObsolete) String() string {
	return "Not Implemented: Obsolete Record Type"
}

//-----------------------------------------------------------------------------
// A Record RDATA
//-----------------------------------------------------------------------------

// RDataA represents an A Record
type RDataA struct {
	ipAddr net.IP
}

// NewRDataA creates a new RDataA instance
func NewRDataA(data []byte) (*RDataA, error) {
	return &RDataA{
		ipAddr: net.IPv4(data[0], data[1], data[2], data[3]),
	}, nil
}

// String makes this record printable
func (r *RDataA) String() string {
	return r.ipAddr.String()
}

//-----------------------------------------------------------------------------
// AAAA Record RDATA
//-----------------------------------------------------------------------------

// RDataAAAA represents an A Record
type RDataAAAA struct {
	ipAddr net.IP
}

// NewRDataAAAA creates a new RDataAAAA instance
func NewRDataAAAA(data []byte) (*RDataAAAA, error) {
	return &RDataAAAA{
		ipAddr: append(make(net.IP, 0, net.IPv6len), data...),
	}, nil
}

// String makes this record printable
func (r *RDataAAAA) String() string {
	return r.ipAddr.String()
}

//-----------------------------------------------------------------------------
// CNAME Record RDATA
//-----------------------------------------------------------------------------

// RDataCNAME represents an A Record
type RDataCNAME struct {
	domain string
}

// NewRDataCNAME creates a new RDataCNAME instance
func NewRDataCNAME(data []byte, offset int) (*RDataCNAME, error) {
	domain, _, err := getPrintableDomainStr(data, offset)

	if err != nil {
		return nil, err
	}

	return &RDataCNAME{domain: domain}, nil
}

// String makes this record printable
func (r *RDataCNAME) String() string {
	return r.domain
}

//-----------------------------------------------------------------------------
// NS Record RDATA
//-----------------------------------------------------------------------------

// RDataNS represents an A Record
type RDataNS struct {
	domain string
}

// NewRDataNS creates a new RDataNS instance
func NewRDataNS(data []byte, offset int) (*RDataNS, error) {
	domain, _, err := getPrintableDomainStr(data, offset)

	if err != nil {
		return nil, err
	}

	return &RDataNS{domain: domain}, nil
}

// String makes this record printable
func (r *RDataNS) String() string {
	return r.domain
}

//-----------------------------------------------------------------------------
// NS Record RDATA
//-----------------------------------------------------------------------------

// RDataTXT represents an A Record
type RDataTXT struct {
	txt string
}

// NewRDataTXT creates a new RDataTXT instance
func NewRDataTXT(data []byte) (*RDataTXT, error) {
	return &RDataTXT{txt: string(data)}, nil
}

// String makes this record printable
func (r *RDataTXT) String() string {
	return r.txt
}

//-----------------------------------------------------------------------------
// SOA Record RDATA
//-----------------------------------------------------------------------------

// RDataSOA represents an A Record
type RDataSOA struct {
	mname   string
	rname   string
	serial  uint32
	refresh uint32
	retry   uint32
	expire  uint32
	minimum uint32
}

// NewRDataSOA creates a new RDataSOA instance
func NewRDataSOA(data []byte, offset int) (*RDataSOA, error) {
	mname, offset, err := getPrintableDomainStr(data, offset)
	if err != nil {
		return nil, err
	}

	rname, offset, err := getPrintableDomainStr(data, offset)
	if err != nil {
		return nil, err
	}

	serial, offset, err := decodeUint32(data, offset)
	if err != nil {
		return nil, err
	}

	refresh, offset, err := decodeUint32(data, offset)
	if err != nil {
		return nil, err
	}

	retry, offset, err := decodeUint32(data, offset)
	if err != nil {
		return nil, err
	}

	expire, offset, err := decodeUint32(data, offset)
	if err != nil {
		return nil, err
	}

	minimum, offset, err := decodeUint32(data, offset)
	if err != nil {
		return nil, err
	}

	return &RDataSOA{
		mname:   mname,
		rname:   rname,
		serial:  serial,
		refresh: refresh,
		retry:   retry,
		expire:  expire,
		minimum: minimum,
	}, nil
}

// String makes this record printable
func (r *RDataSOA) String() string {
	return fmt.Sprintf("%s %s %d %d %d %d %d", r.mname, r.rname, r.serial, r.refresh, r.retry, r.expire, r.minimum)
}

//-----------------------------------------------------------------------------
// MX Record RDATA
//-----------------------------------------------------------------------------

// RDataMX represents an A Record
type RDataMX struct {
	preference uint16
	exchange   string
}

// NewRDataMX creates a new RDataMX instance
func NewRDataMX(data []byte, offset int) (*RDataMX, error) {
	preference, offset, err := decodeUint16(data, offset)

	if err != nil {
		return nil, err
	}

	exchange, _, err := getPrintableDomainStr(data, offset)

	if err != nil {
		return nil, err
	}

	return &RDataMX{
		preference: preference,
		exchange:   exchange,
	}, nil
}

// String makes this record printable
func (r *RDataMX) String() string {
	return fmt.Sprintf("%d %s", r.preference, r.exchange)
}

//-----------------------------------------------------------------------------
// PTR Record RDATA
//-----------------------------------------------------------------------------

// RDataPTR represents an A Record
type RDataPTR struct {
	domain string
}

// NewRDataPTR creates a new RDataPTR instance
func NewRDataPTR(data []byte, offset int) (*RDataPTR, error) {
	domain, _, err := getPrintableDomainStr(data, offset)

	if err != nil {
		return nil, err
	}

	return &RDataPTR{domain: domain}, nil
}

// String makes this record printable
func (r *RDataPTR) String() string {
	return r.domain
}
