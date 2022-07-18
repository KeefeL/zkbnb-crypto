package zecrey_legend

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

func KeccakHash(value []byte) []byte {
	hashVal := crypto.Keccak256Hash(value)
	return hashVal[:]
}

func ComputeAccountNameHash(accountName string) (res string, err error) {
	words := strings.Split(accountName, ".")
	if len(words) != 2 {
		return "", errors.New("[AccountNameHash] invalid account name")
	}
	buf := make([]byte, 32)
	label := KeccakHash([]byte(words[0]))
	res = common.Bytes2Hex(
		KeccakHash(append(
			KeccakHash(append(buf,
				KeccakHash([]byte(words[1]))...)), label...)))
	return res, nil
}

func GetAccountNameHash(accountName string) (hashVal string, err error) {
	// xxx.legend
	nameHash, err := ComputeAccountNameHash(accountName)
	if err != nil {
		return "", err
	}
	return nameHash, nil
}