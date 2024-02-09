package activityPub

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

	pem, err := getPem(*key)

	if err != nil {
		t.Fail()
	}

	publicKey := PublicKey{
		Id: "https://example.com/id",
		Owner: "https://example.com/key",
		PublicKeyPem: *pem,
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, []byte(signatureString))

	header := map[string][]string{
		"Signature":    {"keyId=\"https://example.com/key\""},
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

func getPem(privateKey rsa.PrivateKey) (*string, error) {
	// Extract public key from private key
	publicKey := &privateKey.PublicKey

	// Marshal the public key to DER format
	pubDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Failed to marshal public key:", err)
		return nil, nil
	}

	// Create PEM block for the public key
	pubBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	}

	// Encode the PEM block to string
	pubPEM := string(pem.EncodeToMemory(pubBlock))

	return &pubPEM, nil
}
