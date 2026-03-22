package blurhash

import "testing"

func TestBase83RoundTrip(t *testing.T) {
	cases := []struct {
		value  int
		length int
	}{
		{0, 1},
		{82, 1},
		{83, 2},
		{12345, 3},
		{16777215, 4},
	}
	for _, tc := range cases {
		enc := encode83(tc.value, tc.length)
		got, err := decode83(enc)
		if err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if got != tc.value {
			t.Fatalf("roundtrip mismatch: got=%d want=%d", got, tc.value)
		}
	}
}
