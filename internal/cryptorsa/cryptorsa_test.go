package cryptorsa

import (
	"errors"
	"testing"
)

func TestAll(t *testing.T) {
	pkey, err := LoadPublicKey("../../tls/key.pub")
	if err != nil {
		t.Error(err)
	}
	cryptMess, err := EncryptOAEP(pkey, []byte("Hello world"))
	if err != nil {
		t.Error(err)
	}
	privKey, err := LoadPrivateKey("../../tls/priv.pem")
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
}
