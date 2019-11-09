package main

type (
	// QRType is used in the header denoting whether we have a Query or Response
	QRType byte

	// RecordType indicates the type of DNS record
	RecordType uint16

	// RecordClass indicates the class type requested and received
	RecordClass uint16

	// Opcode indicates the type of query being done in the header
	Opcode byte

	// ResponseCode indicates the type of response returned from the server
	ResponseCode byte
)

// These are all of the different constants used in the application
const (
	QRTypeQuery    QRType = 0
	QRTypeResponse QRType = 1

	RecordTypeA          RecordType = 1
	RecordTypeNS         RecordType = 2
	RecordTypeMD         RecordType = 3 // obsolete
	RecordTypeMF         RecordType = 4 // obsolete
	RecordTypeCNAME      RecordType = 5
	RecordTypeSOA        RecordType = 6
	RecordTypeMB         RecordType = 7  // obsolete
	RecordTypeMG         RecordType = 8  // obsolete
	RecordTypeMR         RecordType = 9  // obsolete
	RecordTypeNULL       RecordType = 10 // obsolete
	RecordTypeWKS        RecordType = 11 // obsolete
	RecordTypePTR        RecordType = 12
	RecordTypeHINFO      RecordType = 13 // obsolete
	RecordTypeMINFO      RecordType = 14 // obsolete
	RecordTypeMX         RecordType = 15
	RecordTypeTXT        RecordType = 16
	RecordTypeRP         RecordType = 17 // obsolete
	RecordTypeAFSDB      RecordType = 18
	RecordTypeX25        RecordType = 19 // obsolete
	RecordTypeISDN       RecordType = 20 // obsolete
	RecordTypeRT         RecordType = 21 // obsolete
	RecordTypeNSAP       RecordType = 22 // obsolete
	RecordTypeNSAPPTR    RecordType = 23 // obsolete
	RecordTypeSIG        RecordType = 24 // obsolete
	RecordTypeKEY        RecordType = 25 // obsolete
	RecordTypePX         RecordType = 26 // obsolete
	RecordTypeGPOS       RecordType = 27 // obsolete
	RecordTypeAAAA       RecordType = 28
	RecordTypeLOC        RecordType = 29
	RecordTypeNXT        RecordType = 30 // obsolete
	RecordTypeEID        RecordType = 31 // obsolete
	RecordTypeNIMLOC     RecordType = 32 // obsolete
	RecordTypeSRV        RecordType = 33
	RecordTypeATMA       RecordType = 34 // obsolete
	RecordTypeNAPTR      RecordType = 35
	RecordTypeKX         RecordType = 36
	RecordTypeCERT       RecordType = 37
	RecordTypeA6         RecordType = 38 // obsolete
	RecordTypeDNAME      RecordType = 39
	RecordTypeSINK       RecordType = 40 // obsolete
	RecordTypeOPT        RecordType = 41 // pseudo resource record
	RecordTypeAPL        RecordType = 42 // obsolete
	RecordTypeDS         RecordType = 43
	RecordTypeSSHFP      RecordType = 44
	RecordTypeIPSECKEY   RecordType = 45
	RecordTypeRRSIG      RecordType = 46
	RecordTypeNSEC       RecordType = 47
	RecordTypeDNSKEY     RecordType = 48
	RecordTypeDHCID      RecordType = 49
	RecordTypeNSEC3      RecordType = 50
	RecordTypeNSEC3PARAM RecordType = 51
	RecordTypeTLSA       RecordType = 52
	RecordTypeSMIMEA     RecordType = 53
	RecordTypeHIP        RecordType = 55
	RecordTypeNINFO      RecordType = 56 // obsolete
	RecordTypeRKEY       RecordType = 57 // obsolete
	RecordTypeTALINK     RecordType = 58 // obsolete
	RecordTypeCDS        RecordType = 59
	RecordTypeCDNSKEY    RecordType = 60
	RecordTypeOPENPGPKEY RecordType = 61
	RecordTypeCSYNC      RecordType = 62
	RecordTypeSPF        RecordType = 99  // obsolete
	RecordTypeUINFO      RecordType = 100 // obsolete
	RecordTypeUID        RecordType = 101 // obsolete
	RecordTypeGID        RecordType = 102 // obsolete
	RecordTypeUNSPEC     RecordType = 103 // obsolete
	RecordTypeNID        RecordType = 104 // obsolete
	RecordTypeL32        RecordType = 105 // obsolete
	RecordTypeL64        RecordType = 106 // obsolete
	RecordTypeLP         RecordType = 107 // obsolete
	RecordTypeEUI48      RecordType = 108 // obsolete
	RecordTypeEUI64      RecordType = 109 // obsolete
	RecordTypeTKEY       RecordType = 249
	RecordTypeTSIG       RecordType = 250
	RecordTypeIXFR       RecordType = 251 // pseudo resource record
	RecordTypeAXFR       RecordType = 252 // pseudo resource record
	RecordTypeMAILB      RecordType = 253 // obsolete
	RecordTypeMAILA      RecordType = 254 // obsolete
	RecordTypeWildcard   RecordType = 255 // pseudo resource record
	RecordTypeURI        RecordType = 256
	RecordTypeCAA        RecordType = 257
	RecordTypeDOA        RecordType = 259 // obsolete
	RecordTypeTA         RecordType = 32768
	RecordTypeDLV        RecordType = 32769

	RecordClassUnknown  RecordClass = 0
	RecordClassIN       RecordClass = 1
	RecordClassCS       RecordClass = 2
	RecordClassCH       RecordClass = 3
	RecordClassHS       RecordClass = 4
	RecordClassWildcard RecordClass = 255

	OpcodeQuery  Opcode = 0
	OpcodeIQuery Opcode = 1
	OpcodeStatus Opcode = 2

	ResponseCodeNoError        ResponseCode = 0
	ResponseCodeFormatError    ResponseCode = 1
	ResponseCodeServerFailure  ResponseCode = 2
	ResponseCodeNameError      ResponseCode = 3
	ResponseCodeNotImplemented ResponseCode = 4
	ResponseCodeRefused        ResponseCode = 5
)
