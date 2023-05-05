package signutil

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// 算法类型. 安全性递增.
type Algo int

const (
	MD5 Algo = iota
	Sha256
	HMACMD5
	HMACSha256
)

func Sign(data, sk string, algo Algo) (signature []byte) {
	switch algo {
	case MD5:
		temp := md5.Sum([]byte(data + sk))
		signature = temp[:]
	case Sha256:
		temp := sha256.Sum256([]byte(data + sk))
		signature = temp[:]
	case HMACMD5:
		signature = hmac.New(md5.New, []byte(sk)).Sum([]byte(data))
	case HMACSha256:
		signature = hmac.New(sha256.New, []byte(sk)).Sum([]byte(data))
	}
	return signature
}

func SignHex(data, sk string, algo Algo) (signatureHex string) {
	sign := Sign(data, sk, algo)
	return hex.EncodeToString(sign[:])
}
