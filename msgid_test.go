package msgid

import (
	"math"
	"testing"
	"time"
)

// Set our node or spawnerId for the tests.
var spawnerId int32 = 3

func loop(id int, loops int, gen *MsgIdGenerator) int {
	for i := 0; i < loops; i++ {
		_ = gen.Next()
	}

	return id

}

// Benchmark creating and converting the msg id to a token.
func BenchmarkGenNextKey(b *testing.B) {
	gen := New(spawnerId)
	for n := 0; n < b.N; n++ {
		gen.NextToken()
	}

}

// Benchmark just creating a msg id N times.
func BenchmarkGenNextId(b *testing.B) {

	gen := New(spawnerId)

	for n := 0; n < b.N; n++ {
		gen.Next()
	}

}

func BenchmarkEncodeB64(b *testing.B) {

	gen := New(spawnerId)
	msgId := gen.Next()
	for n := 0; n < b.N; n++ {
		msgId.EncodeB64()
	}
}

func TestMaxSpawnerId(t *testing.T) {
	sId := int32(math.MaxInt32)
	gen := New(sId)
	msgId := gen.Next()
	if msgId.SpawnerId() != sId {
		t.Error("Max SpawnerId failed")
	}

}

func TestEncodeB64(t *testing.T) {
	gen := New(spawnerId)
	msgId := gen.Next()
	str := msgId.EncodeB64()
	t.Log(str)

	newId, err := DecodeB64(str)
	if err != nil {
		t.Error("Decode of new msgId base64 string failed.")
	}

	if newId.Id != msgId.Id {
		t.Error("MsgId base 64 decode error. Id not equal.")
	}

	if newId.Millis != msgId.Millis {
		t.Error("MsgId base 64 decode error. Millis not equal.")
	}

}

func TestMsgIdGenerator(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestMsgIdGenerator")
	}
	gen := New(spawnerId)
	test_loops := 8
	loopCount := 1000000
	c := make(chan int, test_loops)

	for i := 0; i < test_loops; i++ {
		id := i
		go func() {
			c <- loop(id+1, loopCount, gen)
		}()
	}

	total := 0
	for i := 0; i < test_loops; i++ {
		id := <-c
		total += id
	}

	t.Log("Sum of ids returned from goroutines (Ignore): ", total)
	/*
		nextVal := gen.Next()
		t.Log(nextVal.Millis, nextVal.Id, nextVal.SpawnerId(), nextVal.SequenceVal())
		if nextVal.SequenceVal() != int32(test_loops*loopCount+1) {
			t.Error("Next sequence should have been:", test_loops*loopCount+1, "was ", nextVal.SequenceVal())
		}
	*/
	time.Sleep(1 * time.Millisecond)
	nextVal := gen.Next()
	t.Log(nextVal.Millis, nextVal.Id, nextVal.SpawnerId(), nextVal.SequenceVal())
}

func TestSpawnerId(t *testing.T) {
	genner := New(spawnerId)
	msgId := genner.Next()

	if msgId.SpawnerId() != spawnerId {
		t.Error("SpawnerId was not correct")
	}
}

func TestTime(t *testing.T) {
	genner := New(spawnerId)
	now := time.Now()
	msgId := genner.Next()

	if msgId.Time().Before(now.Add(-5 * time.Millisecond)) {
		t.Error("The MsgId shouldn't be less than the time.")
	}

	if msgId.Time().After(now.Add(5 * time.Millisecond)) {
		t.Error("The MsgId shouldn't be that far ahead of the time .")
	}
}
