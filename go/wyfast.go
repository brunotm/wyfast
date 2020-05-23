package wyfast

import (
	"encoding/binary"
	"math/bits"
	"sync/atomic"
)

const (
	wyp0 = 0xa0761d6478bd642f
	wyp1 = 0xe7037ed1a0b428db
	wyp2 = 0x8ebc6af09c88c6e3
	wyp3 = 0x589965cc75374cc3
	wyp4 = 0x1d8e4e27c47d124f
)

// Sum64 returns the 64-bit wyfast hash value for the given key and seed.
func Sum64(key []byte, seed uint64) uint64 {
	seed ^= wyp4
	p := key

	if len(p) == 0 {
		return wymix(mulm64(wyp0, seed), wyp4)
	}

	if len(p) > 64 {

		for ; len(p) > 256; p = p[256:] {
			seed = mulm64(binary.LittleEndian.Uint64(p[0:8])^wyp0, binary.LittleEndian.Uint64(p[8:16])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[16:24])^wyp1, binary.LittleEndian.Uint64(p[24:32])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[32:40])^wyp2, binary.LittleEndian.Uint64(p[40:48])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[48:56])^wyp3, binary.LittleEndian.Uint64(p[56:64])^seed)

			seed = mulm64(binary.LittleEndian.Uint64(p[64:72])^wyp0, binary.LittleEndian.Uint64(p[72:80])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[80:88])^wyp1, binary.LittleEndian.Uint64(p[88:96])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[96:104])^wyp2, binary.LittleEndian.Uint64(p[104:112])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[112:120])^wyp3, binary.LittleEndian.Uint64(p[120:128])^seed)

			seed = mulm64(binary.LittleEndian.Uint64(p[128:136])^wyp0, binary.LittleEndian.Uint64(p[136:144])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[144:152])^wyp1, binary.LittleEndian.Uint64(p[152:160])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[160:168])^wyp2, binary.LittleEndian.Uint64(p[168:176])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[176:184])^wyp3, binary.LittleEndian.Uint64(p[184:192])^seed)

			seed = mulm64(binary.LittleEndian.Uint64(p[192:200])^wyp0, binary.LittleEndian.Uint64(p[200:208])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[208:216])^wyp1, binary.LittleEndian.Uint64(p[216:224])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[224:232])^wyp2, binary.LittleEndian.Uint64(p[232:240])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[240:248])^wyp3, binary.LittleEndian.Uint64(p[248:256])^seed)
		}

		for ; len(p) > 64; p = p[64:] {
			seed = mulm64(binary.LittleEndian.Uint64(p[0:8])^wyp0, binary.LittleEndian.Uint64(p[8:16])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[16:24])^wyp1, binary.LittleEndian.Uint64(p[24:32])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[32:40])^wyp2, binary.LittleEndian.Uint64(p[40:48])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[48:56])^wyp3, binary.LittleEndian.Uint64(p[56:64])^seed)

		}
	}

	if len(p) > 0 {
		switch {
		case len(p) < 4:
			seed = mulm64(read3(p)^wyp0, seed^wyp1)
		case len(p) <= 8:
			seed = mulm64(uint64(binary.LittleEndian.Uint32(p[0:4]))^wyp0, uint64(binary.LittleEndian.Uint32(p[len(p)-4:len(p)]))^seed)
		case len(p) <= 16:
			seed = mulm64(binary.LittleEndian.Uint64(p[0:8])^wyp0, binary.LittleEndian.Uint64(p[len(p)-8:len(p)])^seed)
		case len(p) <= 32:
			seed = mulm64(binary.LittleEndian.Uint64(p[0:8])^wyp0, binary.LittleEndian.Uint64(p[8:16])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[len(p)-16:len(p)])^wyp1, binary.LittleEndian.Uint64(p[len(p)-8:len(p)])^seed)
		case len(p) <= 64:
			seed = mulm64(binary.LittleEndian.Uint64(p[0:8])^wyp0, binary.LittleEndian.Uint64(p[8:16])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[16:24])^wyp1, binary.LittleEndian.Uint64(p[24:32])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[len(p)-32:len(p)])^wyp2, binary.LittleEndian.Uint64(p[len(p)-24:len(p)])^seed) ^
				mulm64(binary.LittleEndian.Uint64(p[len(p)-16:len(p)])^wyp3, binary.LittleEndian.Uint64(p[len(p)-8:len(p)])^seed)
		}
	}

	return wymix(seed^uint64(len(key)), wyp4)
}

func read3(p []byte) uint64 {
	return (uint64(p[0]) << 16) | (uint64(p[len(p)>>1]) << 8) | uint64(p[len(p)-1])
}

func mulm64(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

func wymix(a, b uint64) uint64 {
	return a ^ b ^ mulm64(a, b)
}

// Rng implements a pseudo-random 64-bit number generator for wyfast
type Rng uint64

// Uint64 returns a pseudo-random 64-bit value as a uint64 from the Rng.
// It is safe to call Uint64 concurrently.
func (r *Rng) Uint64() uint64 {
	x := atomic.AddUint64((*uint64)(r), wyp0)
	return mulm64(x^wyp1, x)
}
