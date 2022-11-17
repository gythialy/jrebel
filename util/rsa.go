package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
)

const PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAND3cI/pKMSd4OLM
IXU/8xoEZ/nza+g00Vy7ygyGB1Nn83qpro7tckOvUVILJoN0pKw8J3E8rtjhSyr9
849qzaQKBhxFL+J5uu08QVn/tMt+Tf0cu5MSPOjT8I2+NWyBZ6H0FjOcVrEUMvHt
8sqoJDrDU4pJyex2rCOlpfBeqK6XAgMBAAECgYBM5C+8FIxWxM1CRuCs1yop0aM8
2vBC0mSTXdo7/3lknGSAJz2/A+o+s50Vtlqmll4drkjJJw4jacsR974OcLtXzQrZ
0G1ohCM55lC3kehNEbgQdBpagOHbsFa4miKnlYys537Wp+Q61mhGM1weXzosgCH/
7e/FjJ5uS6DhQc0Y+QJBAP43hlSSEo1BbuanFfp55yK2Y503ti3Rgf1SbE+JbUvI
IRsvB24xHa1/IZ+ttkAuIbOUomLN7fyyEYLWphIy9kUCQQDSbqmxZaJNRa1o4ozG
RORxR2KBqVn3EVISXqNcUH3gAP52U9LcnmA3NMSZs8tzXhUhYkWQ75Q6umXvvDm4
XZ0rAkBoymyWGeyJy8oyS/fUW0G63mIroZZ4Rp+F098P3j9ueJ2k/frbImXwabJr
hwjUZe/Afel+PxL2ElUDkQW+BMHdAkEAk/U7W4Aanjpfs1+Xm9DUztFicciheRa0
njXspvvxhY8tXAWUPYseG7L+iRPh+Twtn0t5nm7VynVFN0shSoCIAQJALjo7A6bz
svfnJpV+lQiOqD/WCw3A2yPwe+1d0X/13fQkgzcbB3K0K81Euo/fkKKiBv0A7yR7
wvrNjzefE9sKUw==
-----END RSA PRIVATE KEY-----`
const SERVER_RANDOM_NESS = "H2ulzLlh7E0="

type RSA struct {
}

// GetSign 对消息的散列值进行数字签名
func (r RSA) GetSign(msg string) string {
	//计算散列值
	hash := sha1.New()
	hash.Write([]byte(msg))
	bytes := hash.Sum(nil)
	//SignPKCS1v15使用RSA PKCS#1 v1.5规定的RSASSA-PKCS1-V1_5-SIGN签名方案计算签名
	sign, err := rsa.SignPKCS1v15(rand.Reader, r.GetRSAPrivateKey(), crypto.SHA1, bytes)
	if err != nil {
		panic(sign)
	}
	return base64.StdEncoding.EncodeToString(sign)
}
func (RSA) GetRSAPrivateKey() *rsa.PrivateKey {
	//pem解码
	block, _ := pem.Decode([]byte(PRIVATE_KEY))
	//X509解码
	pk, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	return pk.(*rsa.PrivateKey)
}
func (r RSA) Sign(clientRandomness string, guid string, offline bool, validFrom string, validUntil string) string {
	var s2 = ""
	if offline {
		s2 = strings.Join([]string{clientRandomness, SERVER_RANDOM_NESS, guid, "true", validFrom, validUntil}, ";")
	} else {
		s2 = strings.Join([]string{clientRandomness, SERVER_RANDOM_NESS, guid, "false"}, ";")
	}
	return r.GetSign(s2)
}
