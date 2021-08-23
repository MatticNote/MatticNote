package misc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

const (
	PrivateKeyType = "PRIVATE KEY"
	PublicKeyType  = "PUBLIC KEY"
)

func GenerateRSAKeypair(length int) ([]byte, []byte) {
	rsaKeyRaw, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  PrivateKeyType,
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKeyRaw),
	})
	rsaPublicKey := pem.EncodeToMemory(&pem.Block{
		Type:  PublicKeyType,
		Bytes: x509.MarshalPKCS1PublicKey(rsaKeyRaw.Public().(*rsa.PublicKey)),
	})

	return rsaPrivateKey, rsaPublicKey
}
