package str

import (
	"bytes"
	"math/rand"
	"time"
)

var (
	base62Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	r           = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomBase62Str generates a random base-62 string of the specified length.
func RandomBase62Str(l int) string {
	return RandomStr(l, base62Chars)
}

// RandomBase62Str generates a random string from characters in the pool of the
// specified length.
func RandomStr(l int, pool string) string {
	s := &bytes.Buffer{}
	pLen := len(pool)
	for i := 0; i < l; i++ {
		s.WriteByte(pool[r.Intn(pLen)])
	}
	return s.String()
}
