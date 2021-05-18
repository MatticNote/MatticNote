package misc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateRSAKeypair(length int) ([]byte, []byte) {
	rsaKeyRaw, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(rsaKeyRaw),
	})
	rsaPublicKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(rsaKeyRaw.Public().(*rsa.PublicKey)),
	})

	return rsaPrivateKey, rsaPublicKey
}
