package wyfast

import (
	"hash/maphash"
	"math/rand"
	"testing"
	"time"
)

const (
	fibSeed = 0x16069317E428CA9
)

func init() {
	rand.Seed(time.Now().UnixNano())

}

func TestBasicRng(t *testing.T) {
	rng := NewRng(fibSeed)
	var prev uint64

	for x := 0; x < 10000; x++ {
		cur := rng.Uint64()
		if cur == prev {
			t.Fatalf("Rng collision, run=%d, cur=%d, prev=%d", x, cur, prev)
		}
		prev = cur
	}
}

func BenchmarkWyFastRng(b *testing.B) {
	rng := NewRng(fibSeed)
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
					hashFuncs[h].sum(benchSet[x].keys[n%1000], fibSeed)
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
