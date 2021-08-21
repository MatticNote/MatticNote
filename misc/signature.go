package misc

import "github.com/go-fed/httpsig"

func GetHttpSignatureMethod() (httpsig.Signer, error) {
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{
			httpsig.RSA_SHA256,
		},
		httpsig.DigestSha256,
		[]string{
			httpsig.RequestTarget,
			"date",
			"host",
			"digest",
		},
		httpsig.Signature,
		0,
	)
	return signer, err
}
