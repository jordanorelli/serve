package main

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
)

// NewObjectId returns a new unique ObjectId.
// This function causes a runtime error if it fails to get the hostname
// of the current machine.
func newRequestId() RequestId {
	b := make([]byte, 12)
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b, uint32(time.Now().Unix()))
	b[4] = global.machineId[0]
	b[5] = global.machineId[1]
	b[6] = global.machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	b[7] = byte(global.pid >> 8)
	b[8] = byte(global.pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&global.requestCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return RequestId(b)
}

// RequestId is used for tagging each incoming http request for logging
// purposes.  The actual implementation is just the ObjectId implementation
// found in launchpad.net/mgo/bson.  This will most likely change and evolve
// into its own format.
type RequestId string

func (id RequestId) String() string {
	return fmt.Sprintf("%x", string(id))
}

// Time returns the timestamp part of the id.
// It's a runtime error to call this method with an invalid id.
func (id RequestId) Time() time.Time {
	secs := int64(binary.BigEndian.Uint32(id.byteSlice(0, 4)))
	return time.Unix(secs, 0)
}

// byteSlice returns byte slice of id from start to end.
// Calling this function with an invalid id will cause a runtime panic.
func (id RequestId) byteSlice(start, end int) []byte {
	if len(id) != 12 {
		panic(fmt.Sprintf("Invalid RequestId: %q", string(id)))
	}
	return []byte(string(id)[start:end])
}
