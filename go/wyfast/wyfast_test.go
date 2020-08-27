package wyfast

import (
	"bytes"
	"hash/maphash"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

func TestRngUint64(t *testing.T) {
	rng := NewRng(wyp5)
	var prev uint64

	for x := 0; x < 10000; x++ {
		cur := rng.Uint64()
		if cur == prev {
			t.Fatalf("Rng collision, run=%d, cur=%d, prev=%d", x, cur, prev)
		}
		prev = cur
	}
}

func TestRngRead(t *testing.T) {
	rng := NewRng(wyp5)
	size := 16
	cur := make([]byte, size)
	prev := make([]byte, size)

	for x := 0; x < 10000; x++ {
		n, _ := rng.Read(cur)
		if n != size {
			t.Fatalf("Rng wrong read cap, run=%d, expected=%d,  got=%d", x, size, n)
		}

		if bytes.Equal(cur, prev) {
			t.Fatalf("Rng collision, run=%d, cur=%#v, prev=%#v", x, cur, prev)
		}
		copy(prev, cur)
	}
}

func BenchmarkWyFastRng(b *testing.B) {
	rng := NewRng(wyp5)
	for n := 0; n < b.N; n++ {
		rng.Uint64()
	}
}

func BenchmarkGoRandInt63(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rand.Int63()
	}
}

func BenchmarkHash(b *testing.B) {

	var hashFuncs = []struct {
		name string
		sum  func([]byte, uint64) uint64
	}{
		{name: "WyFast", sum: Sum64},
		{name: "GoMapHash", sum: goMapHash},
	}

	var benchSet = []struct {
		name    string
		keySize int64
		keys    [][]byte
	}{
		{name: "-32byte keys", keySize: 32},
		{name: "-64byte keys", keySize: 64},
		{name: "-128byte keys", keySize: 128},
		{name: "-256byte keys", keySize: 256},
		{name: "-512byte keys", keySize: 512},
	}

	for x := 0; x < len(benchSet); x++ {
		for y := 0; y < 1000; y++ {
			benchSet[x].keys = append(benchSet[x].keys, randStringBytes(int(benchSet[x].keySize)))
		}
	}

	for x := 0; x < len(benchSet); x++ {
		for h := 0; h < len(hashFuncs); h++ {
			b.Run(hashFuncs[h].name+benchSet[x].name, func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					hashFuncs[h].sum(benchSet[x].keys[n%1000], wyp5)
					b.SetBytes(benchSet[x].keySize)
				}
			})
		}
	}
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	// for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func goMapHash(key []byte, _ uint64) uint64 {
	h := maphash.Hash{}
	h.Write(key)
	return h.Sum64()
}
