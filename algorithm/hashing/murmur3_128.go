package hashing

import (
	"encoding/binary"
)

const (
	CHUNK_SIZE = 16
)

var (
	C1 uint64 = 0x87c37b91114253d5
	C2 uint64 = 0x4cf5ad432745937f
)

type murmur3_x64_128Func struct {
	seed uint64
}

type murmur3_x64_128Hasher struct {
	h1       uint64
	h2       uint64
	bytes    []byte
	position uint8
	length   int
}

var murmur3_x64_128 = NewMurmur3_128(0)
var good_fast_hash_128 = NewMurmur3_128(good_fast_hash_seed)

func NewMurmur3_128(seed uint64) *murmur3_x64_128Func {
	return &murmur3_x64_128Func{
		seed: seed,
	}
}

func (m *murmur3_x64_128Func) NewHasher() hahser {
	return &murmur3_x64_128Hasher{
		h1: m.seed,
		h2: m.seed,
	}
}

func (m *murmur3_x64_128Hasher) process() {
	k1 := binary.LittleEndian.Uint64(m.bytes[0:4])
	k2 := binary.LittleEndian.Uint64(m.bytes[4:])
	m.bmix64(k1, k2)
	m.length += CHUNK_SIZE
}

func (m *murmur3_x64_128Hasher) bmix64(k1, k2 uint64) {
	m.h1 ^= mixK1(k1)
	m.h1 = rotateLeft(m.h1, 27)
	m.h1 += m.h2
	m.h1 = m.h1*5 + 0x52dce729
	m.h2 ^= mixK2(k2)
	m.h2 = rotateLeft(m.h2, 31)
	m.h2 += m.h1
	m.h2 = m.h2*5 + 0x38495ab5
}

func (m *murmur3_x64_128Hasher) processRemaining() {
	var k1 uint64
	var k2 uint64
	m.length += int(m.position)
	switch m.position {
	case 15:
		k2 ^= uint64(m.bytes[14]) << 48
		fallthrough
	case 14:
		k2 ^= uint64(m.bytes[13]) << 40
		fallthrough
	case 13:
		k2 ^= uint64(m.bytes[12]) << 32
		fallthrough
	case 12:
		k2 ^= uint64(m.bytes[11]) << 24
		fallthrough
	case 11:
		k2 ^= uint64(m.bytes[10]) << 16
		fallthrough
	case 10:
		k2 ^= uint64(m.bytes[9]) << 8
		fallthrough
	case 9:
		k2 ^= uint64(m.bytes[8])
		fallthrough
	case 8:
		k1 ^= uint64(m.bytes[7]) << 56
		fallthrough
	case 7:
		k1 ^= uint64(m.bytes[6]) << 48
		fallthrough
	case 6:
		k1 ^= uint64(m.bytes[5]) << 40
		fallthrough
	case 5:
		k1 ^= uint64(m.bytes[4]) << 32
		fallthrough
	case 4:
		k1 ^= uint64(m.bytes[3]) << 24
		fallthrough
	case 3:
		k1 ^= uint64(m.bytes[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint64(m.bytes[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint64(m.bytes[0])
	default:
		panic("Should never get here.")
	}
	m.h1 ^= mixK1(k1)
	m.h2 ^= mixK2(k2)
}

func (m *murmur3_x64_128Hasher) makeHash() []byte {
	m.h1 ^= uint64(m.length)
	m.h2 ^= uint64(m.length)

	m.h1 += m.h2
	m.h2 += m.h1

	m.h1 = fmix64(m.h1)
	m.h2 = fmix64(m.h2)

	m.h1 += m.h2
	m.h2 += m.h1
	bytes := make([]byte, 16, 16)
	binary.LittleEndian.PutUint64(bytes[0:8], m.h1)
	binary.LittleEndian.PutUint64(bytes[8:], m.h2)
	return bytes
}

func rotateLeft(i uint64, distance uint64) uint64 {
	return i<<distance | (i >> (64 - distance))
}

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}

func mixK1(k1 uint64) uint64 {
	k1 *= C1
	k1 = rotateLeft(k1, 31)
	k1 *= C2
	return k1
}

func mixK2(k2 uint64) uint64 {
	k2 *= C2
	k2 = rotateLeft(k2, 33)
	k2 *= C1
	return k2
}

func (m *murmur3_x64_128Hasher) Sum(bytes []byte) {
	blen := len(bytes)
	if blen > 0 {
		if m.bytes == nil {
			m.bytes = make([]byte, CHUNK_SIZE, CHUNK_SIZE)
		}
		for _, b := range bytes {
			m.bytes[m.position] = b
			m.position++
			if m.position == CHUNK_SIZE {
				m.position = 0
				m.process()
			}
		}
		if m.position > 0 {
			m.processRemaining()
		}
	}
}
