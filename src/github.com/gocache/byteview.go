package gocache

type ByteView struct {
	B []byte
}

func (b ByteView) Len()  int{
	return len(b.B)
}

// ByteSlice returns a copy of the data as a byte slice.
func (v ByteView) ByteSlice() []byte {
	return CloneBytes(v.B)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
