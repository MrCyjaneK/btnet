package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var pgpbits = 4096

// https://gist.github.com/eliquious/9e96017f47d9bd43cdf9
// GenerateKeys - creates a paid of PGP keys
func GenerateKeys(outdir string, keyname string) {
	key, err := rsa.GenerateKey(rand.Reader, pgpbits)
	if err != nil {
		log.Fatal(err)
	}

	priv, err := os.Create(filepath.Join(outdir, keyname+".priv.pgp.asc"))
	if err != nil {
		log.Fatal(err)
	}
	defer priv.Close()

	pub, err := os.Create(filepath.Join(outdir, keyname+".pgp.asc"))
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Close()

	encodePrivateKey(priv, key)
	encodePublicKey(pub, key)

}

func encodePrivateKey(out io.Writer, key *rsa.PrivateKey) {
	w, err := armor.Encode(out, openpgp.PrivateKeyType, make(map[string]string))
	if err != nil {
		log.Fatal(err)
	}

	pgpKey := packet.NewRSAPrivateKey(time.Now(), key)

	err = pgpKey.Serialize(w)
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func encodePublicKey(out io.Writer, key *rsa.PrivateKey) {
	w, err := armor.Encode(out, openpgp.PublicKeyType, make(map[string]string))
	if err != nil {
		log.Fatal(err)
	}
	pgpKey := packet.NewRSAPublicKey(time.Now(), &key.PublicKey)
	err = pgpKey.Serialize(w)
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}
