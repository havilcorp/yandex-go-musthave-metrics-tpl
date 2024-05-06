// Package cryptorsa пакет для шифрования и расшифрования трафика
package cryptorsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func LoadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	pemData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(pemData)
	priv, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	return priv, err
}

func LoadPublicKey(filepath string) (*rsa.PublicKey, error) {
	pemData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(pemData)
	priv, err := x509.ParsePKCS1PublicKey(pemBlock.Bytes)
	return priv, err
}

func EncryptOAEP(public *rsa.PublicKey, msg []byte) ([]byte, error) {
	msgLen := len(msg)
	hash := sha512.New()
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, rand.Reader, public, msg[start:finish], nil)
		if err != nil {
			return nil, err
		}
		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}
	return encryptedBytes, nil
}

func DecryptOAEP(private *rsa.PrivateKey, msg []byte) ([]byte, error) {
	msgLen := len(msg)
	hash := sha512.New()
	step := private.PublicKey.Size()
	var decryptedBytes []byte
	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}
		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, rand.Reader, private, msg[start:finish], nil)
		if err != nil {
			return nil, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}
	return decryptedBytes, nil
}
