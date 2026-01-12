package shortener

import (
	"errors"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Encode(number uint64) string {
	if number == 0 {
		return string(charset[0])
	}
	var res []byte
	for number > 0 {
		res = append(res, charset[number%62])
		number /= 62
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}

func Decode(code string) (uint64, error) {
	var id uint64
	for _, char := range code {
		index := strings.IndexRune(charset, char)
		if index == -1 {
			return 0, errors.New("invalid character")
		}
		id = id*62 + uint64(index)
	}
	return id, nil
}
