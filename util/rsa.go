package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
)

const (
	privateKey = `-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAND3cI/pKMSd4OLMIXU/8xoEZ/nz
a+g00Vy7ygyGB1Nn83qpro7tckOvUVILJoN0pKw8J3E8rtjhSyr9849qzaQKBhxFL+J5uu08QVn/
tMt+Tf0cu5MSPOjT8I2+NWyBZ6H0FjOcVrEUMvHt8sqoJDrDU4pJyex2rCOlpfBeqK6XAgMBAAEC
gYBM5C+8FIxWxM1CRuCs1yop0aM82vBC0mSTXdo7/3lknGSAJz2/A+o+s50Vtlqmll4drkjJJw4j
acsR974OcLtXzQrZ0G1ohCM55lC3kehNEbgQdBpagOHbsFa4miKnlYys537Wp+Q61mhGM1weXzos
gCH/7e/FjJ5uS6DhQc0Y+QJBAP43hlSSEo1BbuanFfp55yK2Y503ti3Rgf1SbE+JbUvIIRsvB24x
Ha1/IZ+ttkAuIbOUomLN7fyyEYLWphIy9kUCQQDSbqmxZaJNRa1o4ozGRORxR2KBqVn3EVISXqNc
UH3gAP52U9LcnmA3NMSZs8tzXhUhYkWQ75Q6umXvvDm4XZ0rAkBoymyWGeyJy8oyS/fUW0G63mIr
oZZ4Rp+F098P3j9ueJ2k/frbImXwabJrhwjUZe/Afel+PxL2ElUDkQW+BMHdAkEAk/U7W4Aanjpf
s1+Xm9DUztFicciheRa0njXspvvxhY8tXAWUPYseG7L+iRPh+Twtn0t5nm7VynVFN0shSoCIAQJA
Ljo7A6bzsvfnJpV+lQiOqD/WCw3A2yPwe+1d0X/13fQkgzcbB3K0K81Euo/fkKKiBv0A7yR7wvrN
jzefE9sKUw==
-----END PRIVATE KEY-----`
)

type Signer struct {
	key *rsa.PrivateKey
}

func NewSigner() *Signer {
	block, _ := pem.Decode([]byte(privateKey))
	pk, _ := x509.ParsePKCS8PrivateKey(block.Bytes)

	return &Signer{key: pk.(*rsa.PrivateKey)}
}

// Sign 对消息的散列值进行数字签名
func (s *Signer) Sign(msg string) (string, error) {
	// 计算散列值
	h := crypto.Hash.New(crypto.SHA1)
	h.Write([]byte(msg))
	hashed := h.Sum(nil)
	// SignPKCS1v15使用RSA PKCS#1 v1.5规定的RSASSA-PKCS1-V1_5-SIGN签名方案计算签名
	if sign, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA1, hashed); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(sign), nil
	}
}

func (s *Signer) SignLease(clientRandomness, serverRandomness, guid string, offline bool, validFrom string, validUntil string) (string, error) {
	data := ""
	if offline {
		data = strings.Join([]string{clientRandomness, serverRandomness, guid, "true", validFrom, validUntil}, ";")
	} else {
		data = strings.Join([]string{clientRandomness, serverRandomness, guid, "false"}, ";")
	}
	return s.Sign(data)
}
