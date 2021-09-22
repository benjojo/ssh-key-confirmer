package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

func makeBogoKey() ssh.Signer {
	Pub, _, _ := ed25519.GenerateKey(rand.Reader)
	a := ed25519PublicKey(Pub)
	return bogoSigner{
		hiddenKey: ssh.PublicKey(a),
	}
}

type ed25519PublicKey ed25519.PublicKey

func (k ed25519PublicKey) Type() string {
	return ssh.KeyAlgoED25519
}

func (k ed25519PublicKey) Marshal() []byte {
	w := struct {
		Name     string
		KeyBytes []byte
	}{
		ssh.KeyAlgoED25519,
		[]byte(k),
	}
	return ssh.Marshal(&w)
}

func (k ed25519PublicKey) Verify(b []byte, sig *ssh.Signature) error {
	if sig.Format != k.Type() {
		return fmt.Errorf("ssh: signature type %s for key type %s", sig.Format, k.Type())
	}
	if l := len(k); l != ed25519.PublicKeySize {
		return fmt.Errorf("ssh: invalid size %d for Ed25519 public key", l)
	}

	if ok := ed25519.Verify(ed25519.PublicKey(k), b, sig.Blob); !ok {
		return errors.New("ssh: signature did not verify")
	}

	return nil
}

type bogoSigner struct {
	hiddenKey ssh.PublicKey
}

func (f bogoSigner) PublicKey() ssh.PublicKey {
	return f.hiddenKey
}

func (f bogoSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	fmt.Printf("Server is attemping any key\n")
	os.Exit(10)
	return nil, nil
}
