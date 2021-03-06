package stream

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kevinburke/nacl"
)

var xSalsa20TestData = []struct {
	in    []byte
	nonce nacl.Nonce
	key   nacl.Key
	out   []byte
}{
	{
		[]byte("Hello world!"),
		&[24]byte{'2', '4', '-', 'b', 'y', 't', 'e', ' ', 'n', 'o', 'n', 'c', 'e', ' ', 'f', 'o', 'r', ' ', 'x', 's', 'a', 'l', 's', 'a'},
		&[32]byte{'t', 'h', 'i', 's', ' ', 'i', 's', ' ', '3', '2', '-', 'b', 'y', 't', 'e', ' ', 'k', 'e', 'y', ' ', 'f', 'o', 'r', ' ', 'x', 's', 'a', 'l', 's', 'a', '2', '0'},
		[]byte{0x00, 0x2d, 0x45, 0x13, 0x84, 0x3f, 0xc2, 0x40, 0xc4, 0x01, 0xe5, 0x41},
	},
	{
		make([]byte, 64),
		&[24]byte{'2', '4', '-', 'b', 'y', 't', 'e', ' ', 'n', 'o', 'n', 'c', 'e', ' ', 'f', 'o', 'r', ' ', 'x', 's', 'a', 'l', 's', 'a'},
		&[32]byte{'t', 'h', 'i', 's', ' ', 'i', 's', ' ', '3', '2', '-', 'b', 'y', 't', 'e', ' ', 'k', 'e', 'y', ' ', 'f', 'o', 'r', ' ', 'x', 's', 'a', 'l', 's', 'a', '2', '0'},
		[]byte{0x48, 0x48, 0x29, 0x7f, 0xeb, 0x1f, 0xb5, 0x2f, 0xb6,
			0x6d, 0x81, 0x60, 0x9b, 0xd5, 0x47, 0xfa, 0xbc, 0xbe, 0x70,
			0x26, 0xed, 0xc8, 0xb5, 0xe5, 0xe4, 0x49, 0xd0, 0x88, 0xbf,
			0xa6, 0x9c, 0x08, 0x8f, 0x5d, 0x8d, 0xa1, 0xd7, 0x91, 0x26,
			0x7c, 0x2c, 0x19, 0x5a, 0x7f, 0x8c, 0xae, 0x9c, 0x4b, 0x40,
			0x50, 0xd0, 0x8c, 0xe6, 0xd3, 0xa1, 0x51, 0xec, 0x26, 0x5f,
			0x3a, 0x58, 0xe4, 0x76, 0x48},
	},
}

func TestXOR(t *testing.T) {
	for i, test := range xSalsa20TestData {
		out := XOR(test.in, test.nonce, test.key)
		if !bytes.Equal(out, test.out) {
			t.Errorf("%d: expected %x, got %x", i, test.out, out)
		}
	}
}

var (
	keyArray [32]byte
	key      = &keyArray
	msg      = make([]byte, 1<<10)
)

var pkgOut []byte

func BenchmarkXOR1K(b *testing.B) {
	b.StopTimer()
	nonce := nacl.NewNonce()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		pkgOut = XOR(msg[:1024], nonce, key)
	}
	b.SetBytes(1024)
}

var firstkey = &[32]byte{
	0x1b, 0x27, 0x55, 0x64, 0x73, 0xe9, 0x85, 0xd4,
	0x62, 0xcd, 0x51, 0x19, 0x7a, 0x9a, 0x46, 0xc7,
	0x60, 0x09, 0x54, 0x9e, 0xac, 0x64, 0x74, 0xf2,
	0x06, 0xc4, 0xee, 0x08, 0x44, 0xf6, 0x83, 0x89,
}

var firstnonce = &[24]byte{
	0x69, 0x69, 0x6e, 0xe9, 0x55, 0xb6, 0x2b, 0x73,
	0xcd, 0x62, 0xbd, 0xa8, 0x75, 0xfc, 0x73, 0xd6,
	0x82, 0x19, 0xe0, 0x03, 0x6b, 0x7a, 0x0b, 0x37,
}

var firstsum = "662b9d0e3463029156069b12f918691a98f7dfb2ca0393c96bbfc6b1fbd630a2"

func TestStream1(t *testing.T) {
	out := Stream(4194304, firstnonce, firstkey)
	result := sha256.Sum256(out)
	if fmt.Sprintf("%x", result) != firstsum {
		t.Errorf("Stream: want %s, got %x", firstsum, result)
	}
}

var thirdkey = &[32]byte{
	0x1b, 0x27, 0x55, 0x64, 0x73, 0xe9, 0x85, 0xd4,
	0x62, 0xcd, 0x51, 0x19, 0x7a, 0x9a, 0x46, 0xc7,
	0x60, 0x09, 0x54, 0x9e, 0xac, 0x64, 0x74, 0xf2,
	0x06, 0xc4, 0xee, 0x08, 0x44, 0xf6, 0x83, 0x89,
}

var thirdnonce = &[24]byte{
	0x69, 0x69, 0x6e, 0xe9, 0x55, 0xb6, 0x2b, 0x73,
	0xcd, 0x62, 0xbd, 0xa8, 0x75, 0xfc, 0x73, 0xd6,
	0x82, 0x19, 0xe0, 0x03, 0x6b, 0x7a, 0x0b, 0x37,
}

var thirdexpected = []byte{
	0xee, 0xa6, 0xa7, 0x25, 0x1c, 0x1e, 0x72, 0x91,
	0x6d, 0x11, 0xc2, 0xcb, 0x21, 0x4d, 0x3c, 0x25,
	0x25, 0x39, 0x12, 0x1d, 0x8e, 0x23, 0x4e, 0x65,
	0x2d, 0x65, 0x1f, 0xa4, 0xc8, 0xcf, 0xf8, 0x80,
}

func TestStream3(t *testing.T) {
	out := Stream(32, thirdnonce, thirdkey)
	if !cmp.Equal(out, thirdexpected) {
		t.Errorf("Stream: want %x, got %x", thirdexpected, out)
	}
}
