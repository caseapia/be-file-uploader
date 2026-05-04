package generate

import (
	"crypto/rand"
	"math/big"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func FileID() (string, error) {
	codeLen := 12
	code := make([]byte, codeLen)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}

		code[i] = letters[n.Int64()]
	}

	return "file__" + string(code), nil
}
