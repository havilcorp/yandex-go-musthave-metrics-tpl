package cryptorsa

import (
	"errors"
	"testing"
)

func TestAll(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		pkey, err := LoadPublicKey("../../tls/key.pub")
		if err != nil {
			t.Error(err)
		}
		cryptMess, err := EncryptOAEP(pkey, []byte("Hello world"))
		if err != nil {
			t.Error(err)
		}
		privKey, err := LoadPrivateKey("../../tls/server.key")
		if err != nil {
			t.Error(err)
		}
		decMess, err := DecryptOAEP(privKey, cryptMess)
		if err != nil {
			t.Error(err)
		}
		if string(decMess) != "Hello world" {
			t.Error(errors.New("decMess is not the same as encMess"))
		}
	})
	t.Run("Bad public key path", func(t *testing.T) {
		_, err := LoadPublicKey("../../tls/none")
		if err == nil {
			t.Error(errors.New("LoadPublicKey"))
		}
	})
	t.Run("Bad private key path", func(t *testing.T) {
		_, err := LoadPrivateKey("../../tls/none")
		if err == nil {
			t.Error(errors.New("LoadPrivateKey"))
		}
	})
}
