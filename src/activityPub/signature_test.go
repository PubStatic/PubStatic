package activityPub

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestValidateSignature(t *testing.T) {
	// Prepare test data
	host := "example.com"
	date := "2022-02-24"
	digest := "digest-data"

	signatureString := fmt.Sprintf("(request-target): post /inbox\nhost: %s\ndate: %s\ndigest: sha-256=%s\ncontent-type: application/json", host, date, digest)

	key, err := rsa.GenerateKey(rand.Reader, 2048)

	pem, err := GetPublicKeyPem(*key)

	if err != nil {
		t.Fail()
	}

	publicKey := PublicKey{
		Id:           "https://example.com/id",
		Owner:        "https://example.com/key",
		PublicKeyPem: *pem,
	}

	hash := sha256.Sum256([]byte(signatureString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		t.FailNow()
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	signatureHeader := "keyId=\"https://example.com/key\",headers=\"(request-target) host date digest content-type\"," +
		fmt.Sprintf("signature=\"%s\"", signatureBase64)

	header := map[string][]string{
		"Signature":    {signatureHeader},
		"Host":         {host},
		"Date":         {date},
		"Digest":       {"sha-256=" + digest},
		"Content-Type": {"application/json"},
	}

	// Call the function being tested
	result, err := validateSignature(header, publicKey)

	// Check the result
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Add more assertions based on the expected behavior of the function
	if !result {
		t.Errorf("Expected true, got false")
	}
}
