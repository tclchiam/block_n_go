package identity

import (
	"bytes"
	"crypto/sha256"
	"testing"
	"golang.org/x/crypto/ripemd160"
	"github.com/tclchiam/block_n_go/crypto"
	"github.com/mr-tron/base58/base58"
)

func TestAddress_Base58(t *testing.T) {
	privateKey := crypto.NewP256PrivateKey()
	address := NewAddress(privateKey)

	input := [][]byte{
		{version},
		address.PublicKeyHash(),
		address.Checksum(),
	}
	expectedBase58 := base58.Encode(bytes.Join(input, []byte{}))

	if len(expectedBase58) != 34 {
		t.Errorf("Expected len did not equal actual. Got: %d, wanted: %d", len(expectedBase58), 34)
	}
	if expectedBase58 != address.Base58() {
		t.Errorf("Expected base58 did not equal actual. Got: '%s', wanted: '%s'", address.Base58(), expectedBase58)
	}
}

func TestAddress_Version(t *testing.T) {
	privateKey := crypto.NewP256PrivateKey()
	address := NewAddress(privateKey)

	expectedVersion := byte(0x00)

	if expectedVersion != address.Version() {
		t.Errorf("Expected version did not equal actual. Got: '%s', wanted: '%s'", address, expectedVersion)
	}
}

func TestAddress_PublicKeyHash(t *testing.T) {
	privateKey := crypto.NewP256PrivateKey()
	publicKey := privateKey.PubKey()
	address := NewAddress(privateKey)

	publicSHA256 := sha256.Sum256(publicKey.Serialize())

	hashImpl := ripemd160.New()
	hashImpl.Write(publicSHA256[:])
	expectedHash := hashImpl.Sum(nil)

	if len(expectedHash) != 20 {
		t.Errorf("Expected len did not equal actual. Got: %d, wanted: %d", len(expectedHash), 20)
	}
	if bytes.Compare(expectedHash, address.PublicKeyHash()) != 0 {
		t.Errorf("Expected hash did not equal actual. Got: '%s', wanted: '%s'", address, expectedHash)
	}
}

func TestAddress_Checksum(t *testing.T) {
	privateKey := crypto.NewP256PrivateKey()
	address := NewAddress(privateKey)

	payload := append([]byte{address.Version()}, address.PublicKeyHash()...)

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	expectedChecksum := secondSHA[:checksumLength]

	if len(expectedChecksum) != 4 {
		t.Errorf("Expected len did not equal actual. Got: %d, wanted: %d", len(expectedChecksum), 4)
	}
	if bytes.Compare(expectedChecksum, address.Checksum()) != 0 {
		t.Errorf("Expected checksum did not equal actual. Got: '%s', wanted: '%s'", address.Checksum(), expectedChecksum)
	}
}