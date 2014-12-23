package msgid

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"strconv"
	"sync"
	"time"
)

// Reserve a byte out of the counter half of the Id.
// Maybe for subnodes or go routines if useful.
var max int64 = 1<<23 - 1

func read_int64(data []byte) (ret int64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

type MsgId struct {
	Millis int64
	Id     int64
}

// Get the SpawnId from the MsgId instance.
func (m MsgId) SpawnerId() int32 {
	return int32(m.Id >> 32)
}

// Get the sequence (counter) from the MsgId instance.  This will overlap
// after enough instances or on a restart.  The overall instance will be unique
// when combined with Millis.  This will rarely be useful by itself.
func (m MsgId) SequenceVal() int32 {
	return int32(m.Id)
}

// Get the milliseconds of the MsgId converted to a time.Time struct.
// The time has millisecond resolution per the name.
func (m MsgId) Time() time.Time {
	return time.Unix(0, m.Millis*1e6)

}

// Encode the MsgId instance as a base64 encoded string.
func (m MsgId) EncodeB64() string {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[0:8], uint64(m.Millis))
	binary.BigEndian.PutUint64(b[8:16], uint64(m.Id))
	ret := base64.StdEncoding.EncodeToString(b)
	return ret
}

type InvalidDataLengthError int64

func (e InvalidDataLengthError) Error() string {
	return "The token to be decoded must be 16 bytes after decoding.  The length was " + strconv.FormatInt(int64(e), 10)
}

// Decode a base64 encoded string into a MsgId instance.
func DecodeB64(input string) (*MsgId, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	if len(data) != 16 {
		return nil, InvalidDataLengthError(len(data))
	}
	ret := new(MsgId)
	ret.Millis = read_int64(data[0:8])
	ret.Id = read_int64(data[8:])
	return ret, nil
}

type counter struct {
	curVal int64
	sync.Mutex
}

func (ctr *counter) next() int64 {
	ctr.Lock()
	defer ctr.Unlock()
	if ctr.curVal == max {
		ctr.curVal = 0
	}
	ctr.curVal++
	return ctr.curVal
}

// MsgIdGenerator struct holding the spawnerId and internal counter.
type MsgIdGenerator struct {
	spawnerId int64
	ctr       *counter
}

// SpawnerId associated with this instance.
func (gen *MsgIdGenerator) SpawnerId() int32 {
	return int32(gen.spawnerId)
}

// Create a new MsgIdGenerator.
func New(spawnerId int32) *MsgIdGenerator {
	gen := &MsgIdGenerator{spawnerId: int64(spawnerId)}
	gen.ctr = new(counter)
	return gen
}

// Get the next MsgId.
func (gen *MsgIdGenerator) Next() MsgId {
	id := gen.spawnerId<<32 | gen.ctr.next()
	return MsgId{time.Now().UnixNano() / 1e6, id}
}

// Convenience method to get the next MsgId as a Base64 encoded string.
func (gen *MsgIdGenerator) NextToken() string {
	msgId := gen.Next()
	return msgId.EncodeB64()
}
