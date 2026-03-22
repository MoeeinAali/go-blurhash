package blurhash

import "fmt"

const digitCharacters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

var base83Index [256]int16

func init() {
	for i := range base83Index {
		base83Index[i] = -1
	}
	for i := 0; i < len(digitCharacters); i++ {
		base83Index[digitCharacters[i]] = int16(i)
	}
}

func decode83(s string) (int, error) {
	value := 0
	for i := 0; i < len(s); i++ {
		idx := base83Index[s[i]]
		if idx < 0 {
			return 0, fmt.Errorf("%w: %q", ErrInvalidCharacter, s[i])
		}
		value = value*83 + int(idx)
	}
	return value, nil
}

func encode83(v int, length int) string {
	if length <= 0 {
		return ""
	}
	out := make([]byte, length)
	for i := 1; i <= length; i++ {
		digit := (v / pow83(length-i)) % 83
		out[i-1] = digitCharacters[digit]
	}
	return string(out)
}

func pow83(exp int) int {
	p := 1
	for i := 0; i < exp; i++ {
		p *= 83
	}
	return p
}
