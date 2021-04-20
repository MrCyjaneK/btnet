package utils

import (
	"bufio"
	"crypto"
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

func decodePublicKey(filename string) *packet.PublicKey {

	// open ascii armored public key
	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	block, err := armor.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	if block.Type != openpgp.PublicKeyType {
		log.Fatal("Invalid private key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		log.Fatal(err)
	}
	key, ok := pkt.(*packet.PublicKey)
	if !ok {
		log.Fatal("Invalid public key")
	}
	return key
}

func decodePrivateKey(filename string) *packet.PrivateKey {

	// open ascii armored private key
	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	block, err := armor.Decode(in)
	if err != nil {
		log.Fatal(err)
	}
	if block.Type != openpgp.PrivateKeyType {
		log.Fatal("Invalid private key file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		log.Fatal(err)
	}

	key, ok := pkt.(*packet.PrivateKey)
	if !ok {
		log.Fatal("Invalid private key")
	}
	return key
}

func SignFile(PublicKey string, PrivateKey string, file string, target string) {
	pubKey := decodePublicKey(PublicKey)
	privKey := decodePrivateKey(PrivateKey)

	signer := createEntityFromKeys(pubKey, privKey)

	keyfile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	targetfile, err := os.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(targetfile)
	err = openpgp.ArmoredDetachSign(writer, signer, keyfile, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func VerifyFile(PublicKey string, SignatureFile string) bool {
	pubKey := decodePublicKey(PublicKey)
	sig := decodeSignature(SignatureFile)

	hash := sig.Hash.New()
	io.Copy(hash, os.Stdin)

	err := pubKey.VerifySignature(hash, sig)
	return err == nil
}

func decodeSignature(filename string) *packet.Signature {

	// open ascii armored public key
	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	block, err := armor.Decode(in)
	if err != nil {
		log.Fatal(err)
	}
	if block.Type != openpgp.SignatureType {
		log.Fatal("Invalid signature file")
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		log.Fatal(err)
	}
	sig, ok := pkt.(*packet.Signature)
	if !ok {
		log.Fatal("Invalid signature")
	}
	return sig
}

func createEntityFromKeys(pubKey *packet.PublicKey, privKey *packet.PrivateKey) *openpgp.Entity {
	config := packet.Config{
		DefaultHash:            crypto.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		CompressionConfig: &packet.CompressionConfig{
			Level: 9,
		},
		RSABits: pgpbits,
	}
	currentTime := config.Now()
	uid := packet.NewUserId("", "", "")

	e := openpgp.Entity{
		PrimaryKey: pubKey,
		PrivateKey: privKey,
		Identities: make(map[string]*openpgp.Identity),
	}
	isPrimaryId := false

	e.Identities[uid.Id] = &openpgp.Identity{
		Name:   uid.Name,
		UserId: uid,
		SelfSignature: &packet.Signature{
			CreationTime: currentTime,
			SigType:      packet.SigTypePositiveCert,
			PubKeyAlgo:   packet.PubKeyAlgoRSA,
			Hash:         config.Hash(),
			IsPrimaryId:  &isPrimaryId,
			FlagsValid:   true,
			FlagSign:     true,
			FlagCertify:  true,
			IssuerKeyId:  &e.PrimaryKey.KeyId,
		},
	}

	keyLifetimeSecs := uint32(86400 * 365)

	e.Subkeys = make([]openpgp.Subkey, 1)
	e.Subkeys[0] = openpgp.Subkey{
		PublicKey:  pubKey,
		PrivateKey: privKey,
		Sig: &packet.Signature{
			CreationTime:              currentTime,
			SigType:                   packet.SigTypeSubkeyBinding,
			PubKeyAlgo:                packet.PubKeyAlgoRSA,
			Hash:                      config.Hash(),
			PreferredHash:             []uint8{8}, // SHA-256
			FlagsValid:                true,
			FlagEncryptStorage:        true,
			FlagEncryptCommunications: true,
			IssuerKeyId:               &e.PrimaryKey.KeyId,
			KeyLifetimeSecs:           &keyLifetimeSecs,
		},
	}
	return &e
}
