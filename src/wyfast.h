/*
* wyfast is a small 64-bit non-cryptographic hash function.
* It is a simpler and faster variation of wyhash v5  https://github.com/wangyi-fudan/wyhash
*/

#include <stdint.h>
#include <string.h>

const uint64_t wyp0 = 0xa0761d6478bd642full;
const uint64_t wyp1 = 0xe7037ed1a0b428dbull;
const uint64_t wyp2 = 0x8ebc6af09c88c6e3ull;
const uint64_t wyp3 = 0x589965cc75374cc3ull;
const uint64_t wyp4 = 0x1d8e4e27c47d124full;

#if defined(__BIG_ENDIAN__) || (defined(__BYTE_ORDER__) && __BYTE_ORDER__ == __ORDER_BIG_ENDIAN__)
	#if defined(__GNUC__) || defined(__INTEL_COMPILER) || defined(__clang__)
		#define in_order64(x) __builtin_bswap64(x)
		#define in_order32(x) __builtin_bswap32(x)
	#elif defined(_MSC_VER)
		#define in_order64(x) _byteswap_uint64(x)
		#define in_order32(x) _byteswap_ulong(x)
	#else
		#error "__BIG_ENDIAN__ not supported with current compiler)"
	#endif
#else
	#define in_order64(x) x
	#define in_order32(x) x
#endif

static inline uint64_t mulm64(uint64_t x, uint64_t y) {
#ifdef __SIZEOF_INT128__

	__uint128_t r = x; r *= y;
	return (r>>64)^r;

#elif defined(_MSC_VER)

	x = _umul128(x, y, &y);
	return x^y;

#else

	uint64_t hi, lo;
	const uint64_t mask32 = ((uint64_t)1<<(uint64_t)32) - (uint64_t)1;

	uint64_t x0 = x & mask32;
	uint64_t x1 = x >> 32;
	uint64_t y0 = y & mask32;
	uint64_t y1 = y >> 32;
	uint64_t w0 = x0 * y0;
	uint64_t t = x1*y0 + (w0>>32);
	uint64_t w1 = t & mask32;
	uint64_t w2 = t >> 32;
	w1 += x0 * y1;
	hi = x1*y1 + w2 + (w1>>32);
	lo = x * y;

	return hi^lo;

#endif
}

static inline uint64_t read8(const uint8_t *p) {
	uint64_t v;
	memcpy(&v, p, 8);
	return in_order64(v);
}

static inline uint64_t read4(const uint8_t *p) {
	uint64_t v;
	memcpy(&v, p, 4);
	return in_order32(v);
}

static inline uint64_t read3(const uint8_t *p, unsigned k) {
	return (((uint64_t)p[0])<<16)|(((uint64_t)p[k>>1])<<8)|p[k-1];
}

static inline uint64_t wymix(uint64_t x, uint64_t y) {
	return x^y^mulm64(x,y);
}

static inline uint64_t wyfast(const void* key, uint64_t len, uint64_t seed){
	const uint8_t *p = (const uint8_t*)key;
	uint64_t i = len;
	seed ^= wyp4;

	if (i == 0) {
		return wymix(mulm64(wyp0, seed), wyp4);
	}

	for(; i>=256; i-=256, p+=256){
		seed = mulm64(read8(p)^wyp0, read8(p+8)^seed)^
			mulm64(read8(p+16)^wyp1, read8(p+24)^seed)^
			mulm64(read8(p+32)^wyp2, read8(p+40)^seed)^
			mulm64(read8(p+48)^wyp3, read8(p+56)^seed);

		seed = mulm64(read8(p+64)^wyp0, read8(p+72)^seed)^
			mulm64(read8(p+80)^wyp1, read8(p+88)^seed)^
			mulm64(read8(p+96)^wyp2, read8(p+104)^seed)^
			mulm64(read8(p+112)^wyp3, read8(p+120)^seed);

		seed = mulm64(read8(p+128)^wyp0, read8(p+136)^seed)^
			mulm64(read8(p+144)^wyp1, read8(p+152)^seed)^
			mulm64(read8(p+160)^wyp2, read8(p+168)^seed)^
			mulm64(read8(p+176)^wyp3, read8(p+184)^seed);

		seed = mulm64(read8(p+192)^wyp0, read8(p+200)^seed)^
			mulm64(read8(p+208)^wyp1, read8(p+216)^seed)^
			mulm64(read8(p+224)^wyp2, read8(p+232)^seed)^
			mulm64(read8(p+240)^wyp3, read8(p+248)^seed);
	}

	for(; i>=64; i-=64, p+=64){
		seed = mulm64(read8(p)^wyp0, read8(p+8)^seed)^
			mulm64(read8(p+16)^wyp1, read8(p+24)^seed)^
			mulm64(read8(p+32)^wyp2, read8(p+40)^seed)^
			mulm64(read8(p+48)^wyp3, read8(p+56)^seed);
	}

	if (i > 0) {
		if (i < 4) {
			seed = mulm64(read3(p, i)^wyp0, seed^wyp1);
			return wymix(seed^len, wyp4);
		}

		if (i <= 8) {
			seed = mulm64(read4(p)^wyp0, read4(p+i-4)^seed);
			return wymix(seed^len, wyp4);
		}

		if (i <= 16) {
			seed = mulm64(read8(p)^wyp0, read8(p+i-8)^seed);
			return wymix(seed^len, wyp4);
		}

		if (i <= 32) {
			seed = mulm64(read8(p)^wyp0, read8(p+8)^seed)^
				mulm64(read8(p+i-16)^wyp1, read8(p+i-8)^seed);
			return wymix(seed^len, wyp4);
		}

		if (i <= 64) {
			seed = mulm64(read8(p)^wyp0, read8(p+8)^seed)^
				mulm64(read8(p+16)^wyp1,read8(p+24)^seed)^
				mulm64(read8(p+i-32)^wyp2,read8(p+i-24)^seed)^
				mulm64(read8(p+i-16)^wyp3,read8(p+i-8)^seed);
			return wymix(seed^len, wyp4);
		}
	}

	return wymix(seed^len, wyp4);
}


static inline uint64_t wyfast_rng(uint64_t *seed) {
	*seed+=wyp0;
	return mulm64(*seed^wyp1,*seed);
}
