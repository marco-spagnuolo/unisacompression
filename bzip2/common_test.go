
package bzip2

import (
	"strconv"
	"testing"

	"github.com/marco-spagnuolo/unisacompression/internal/testutil"
)

func TestCRC(t *testing.T) {
	vectors := []struct {
		crc uint32
		str string
	}{
		{0x00000000, ""},
		{0x19939b6b, "a"},
		{0xe993fdcd, "ab"},
		{0x648cbb73, "abc"},
		{0x3d4c334b, "abcd"},
		{0xa35b4df4, "abcde"},
		{0xa0f54fb9, "abcdef"},
		{0x077539d7, "abcdefg"},
		{0x5024ec61, "abcdefgh"},
		{0x63e0bcd4, "abcdefghi"},
		{0x73826444, "abcdefghij"},
		
	}

	var crc crc
	for i, v := range vectors {
		splits := []int{
			0 * (len(v.str) / 1),
			1 * (len(v.str) / 4),
			2 * (len(v.str) / 4),
			3 * (len(v.str) / 4),
			1 * (len(v.str) / 1),
		}
		for _, j := range splits {
			str1, str2 := []byte(v.str[:j]), []byte(v.str[j:])
			crc.val = 0
			crc.update(str1)
			if crc.update(str2); crc.val != v.crc {
				t.Errorf("test %d, crc.update(crc1, str2): got 0x%08x, want 0x%08x", i, crc.val, v.crc)
			}
		}
	}
}

func BenchmarkCRC(b *testing.B) {
	var c crc
	d := testutil.ResizeData([]byte("the quick brown fox jumped over the lazy dog"), 1<<16)
	for i := 1; i <= len(d); i <<= 4 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetBytes(int64(i))
			for j := 0; j < b.N; j++ {
				c.update(d[:i])
			}
		})
	}
}
