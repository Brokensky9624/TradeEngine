package auth

import (
	"crypto/tls"
	_ "embed"
	"tradeengine/utils/logger"

	"golang.org/x/crypto/pkcs12"
)

var (
	//go:embed ift.p12
	keyFile []byte
)

func GetTLSConfig() (*tls.Config, error) {
	priv := keyFile
	parsedPrivateKey, parsedCert, err := pkcs12.Decode(priv, "swserver")
	if err != nil {
		logger.SERVER.Warn("Private key decode error " + err.Error())
		return nil, err
	}

	// Gather secure cipher suites, provide by Google
	var cipherSuites []uint16
	for _, cipher := range tls.CipherSuites() {
		cipherSuites = append(cipherSuites, cipher.ID)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{parsedCert.Raw},
		PrivateKey:  parsedPrivateKey,
		Leaf:        parsedCert,
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12, CipherSuites: cipherSuites}, nil
}
