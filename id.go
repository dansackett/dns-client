package main

// NOTE: borrowed from: https://github.com/miekg/dns/blob/master/msg.go

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
)

var (
	idLock sync.Mutex
	idRand *rand.Rand
)

// GenerateRandID returns a 16 bits random number to be used as a
// message id. The random provided should be good enough.
func GenerateRandID() uint16 {
	idLock.Lock()

	if idRand == nil {
		// This (partially) works around
		// https://github.com/golang/go/issues/11833 by only
		// seeding idRand upon the first call to id.

		var seed int64
		var buf [8]byte

		if _, err := crand.Read(buf[:]); err == nil {
			seed = int64(binary.LittleEndian.Uint64(buf[:]))
		} else {
			seed = rand.Int63()
		}

		idRand = rand.New(rand.NewSource(seed))
	}

	// The call to idRand.Uint32 must be within the
	// mutex lock because *rand.Rand is not safe for
	// concurrent use.
	//
	// There is no added performance overhead to calling
	// idRand.Uint32 inside a mutex lock over just
	// calling rand.Uint32 as the global math/rand rng
	// is internally protected by a sync.Mutex.
	id := uint16(idRand.Uint32())

	idLock.Unlock()
	return id
}
