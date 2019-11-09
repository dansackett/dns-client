package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// Message is a request or response to / from a DNS Server
type Message struct {
	queryTime  time.Duration
	bytesRead  int
	Header     Header
	Questions  []Question
	Answers    []RR
	Authority  []RR
	Additional []RR
}

// Encode converts a message object into a DNS-safe message for a server
func (m Message) Encode() ([]byte, error) {
	var err error
	var data bytes.Buffer

	hBytes, err := m.Header.Encode()

	if err != nil {
		return hBytes, err
	}

	data.Write(hBytes)

	for _, question := range m.Questions {
		qBytes, err := question.Encode()

		if err != nil {
			return append(data.Bytes(), qBytes...), err
		}

		data.Write(qBytes)
	}

	for _, answer := range m.Answers {
		anBytes, err := answer.Encode()

		if err != nil {
			return append(data.Bytes(), anBytes...), err
		}

		data.Write(anBytes)
	}

	for _, authority := range m.Authority {
		authBytes, err := authority.Encode()

		if err != nil {
			return append(data.Bytes(), authBytes...), err
		}

		data.Write(authBytes)
	}

	for _, additional := range m.Additional {
		addBytes, err := additional.Encode()

		if err != nil {
			return append(data.Bytes(), addBytes...), err
		}

		data.Write(addBytes)
	}

	return data.Bytes(), err
}

// DecodeMessage decodes a message returned from the DNS server
func DecodeMessage(data []byte, m *Message, queryTime time.Duration) (int, error) {
	var err error
	var bytesRead int

	m.queryTime = queryTime

	//-------------------------------------------------------------------------
	// 1. Decode the header and set it
	//-------------------------------------------------------------------------

	respHeader := new(Header)
	bytesRead, err = DecodeHeader(data, bytesRead, respHeader)

	m.Header = *respHeader

	if err != nil {
		return bytesRead, err
	}

	//-------------------------------------------------------------------------
	// 2. Decode the questions and set them
	//-------------------------------------------------------------------------

	for i := uint16(0); i < m.Header.QDCOUNT; i++ {
		respQuestion := new(Question)
		bytesRead, err = DecodeQuestion(data, bytesRead, respQuestion)

		if err != nil {
			return bytesRead, err
		}

		m.Questions = append(m.Questions, *respQuestion)
	}

	//-------------------------------------------------------------------------
	// 3. Decode the answer RRs and set them
	//-------------------------------------------------------------------------

	for i := uint16(0); i < respHeader.ANCOUNT; i++ {
		rr := new(RR)
		bytesRead, err = DecodeRR(data, bytesRead, rr)

		if err != nil {
			return bytesRead, err
		}

		m.Answers = append(m.Answers, *rr)
	}

	//-------------------------------------------------------------------------
	// 4. Decode the authority RRs and set them
	//-------------------------------------------------------------------------

	for i := uint16(0); i < respHeader.NSCOUNT; i++ {
		rr := new(RR)
		bytesRead, err = DecodeRR(data, bytesRead, rr)

		if err != nil {
			return bytesRead, err
		}

		m.Authority = append(m.Authority, *rr)
	}

	//-------------------------------------------------------------------------
	// 5. Decode the additional RRs and set them
	//-------------------------------------------------------------------------

	for i := uint16(0); i < respHeader.ARCOUNT; i++ {
		rr := new(RR)
		bytesRead, err = DecodeRR(data, bytesRead, rr)

		if err != nil {
			return bytesRead, err
		}

		m.Additional = append(m.Additional, *rr)
	}

	m.bytesRead = bytesRead

	return bytesRead, err
}

func (m *Message) String() string {
	var sb strings.Builder

	domain := *domainFlagVal
	recordType := *recordTypeFlagVal
	dnsServerAddr := *dnsServerAddrFlagVal
	queryTime := m.queryTime
	bytesRead := m.bytesRead
	currentTime := time.Now().Format(time.RFC1123)

	id := m.Header.ID
	opcode := OpcodeToStrMap[m.Header.OPCODE]
	responseCode := ResponseCodeToStrMap[m.Header.RCODE]

	numQuestions := len(m.Questions)
	numAnswers := len(m.Answers)
	numAuthority := len(m.Authority)
	numAdditional := len(m.Additional)

	sb.WriteString(fmt.Sprintf("\n> [ Simple DNS Client ] >>> %s %s", domain, recordType))
	sb.WriteString(fmt.Sprintf("\n> ID: %d, opcode: %s, status: %s", id, opcode, responseCode))
	sb.WriteString(fmt.Sprintf("\n> QUERY: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d,", numQuestions, numAnswers, numAuthority, numAdditional))

	if numQuestions > 0 {
		sb.WriteString(fmt.Sprintf("\n\n> QUESTION SECTION:\n"))

		for _, question := range m.Questions {
			sb.WriteString(fmt.Sprintf("%s\n", question.String()))
		}
	}

	if numAnswers > 0 {
		sb.WriteString(fmt.Sprintf("\n> ANSWER SECTION:\n"))

		for _, answer := range m.Answers {
			sb.WriteString(fmt.Sprintf("%s\n", answer.String()))
		}
	}

	if numAuthority > 0 {
		sb.WriteString(fmt.Sprintf("\n> AUTHORITY SECTION:\n"))

		for _, authority := range m.Authority {
			sb.WriteString(fmt.Sprintf("%s\n", authority.String()))
		}
	}

	if numAdditional > 0 {
		sb.WriteString(fmt.Sprintf("\n> ADDITIONAL SECTION:\n"))

		for _, additional := range m.Additional {
			sb.WriteString(fmt.Sprintf("%s\n", additional.String()))
		}
	}

	sb.WriteString(fmt.Sprintf("\n> Query time: %s", queryTime))
	sb.WriteString(fmt.Sprintf("\n> Server: %s", dnsServerAddr))
	sb.WriteString(fmt.Sprintf("\n> When: %s", currentTime))
	sb.WriteString(fmt.Sprintf("\n> Msg Size: rcvd %d", bytesRead))

	return sb.String()
}
