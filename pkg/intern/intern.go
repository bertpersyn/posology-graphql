package intern

import "sync"

var (
	pool sync.Pool = sync.Pool{
		New: func() interface{} {
			return make(map[string]string)
		},
	}
)

// String returns s, interned.
func String(s string) string {
	m := pool.Get().(map[string]string)
	c, ok := m[s]
	if ok {
		pool.Put(m)
		return c
	}
	m[s] = s
	pool.Put(m)
	return s
}

// Bytes returns b converted to a string, interned.
func Bytes(b []byte) string {
	m := pool.Get().(map[string]string)
	c, ok := m[string(b)]
	if ok {
		pool.Put(m)
		return c
	}
	s := string(b)
	m[s] = s
	pool.Put(m)
	return s
}
