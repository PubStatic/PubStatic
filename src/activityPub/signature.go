package activityPub

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func validateSignature(header http.Header, publicKey PublicKey) (bool, error) {
	logger.Trace("Validating signature")

	signature := header.Get("Signature")

	parts := strings.Split(signature, ",")
	for _, part := range parts {
		logger.Debug("Signature Header Part=", part)
	}

	signatureHash := ""
	headers := ""
	for _, part := range parts {
		if strings.HasPrefix(part, "signature") {
			signatureHash = strings.ReplaceAll(strings.ReplaceAll(string(part), "signature=", ""), "\"", "")
		} else if strings.HasPrefix(part, "headers") {
			headers = strings.ReplaceAll(strings.ReplaceAll(part, "headers=", ""), "\"", "")
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(signatureHash)
	if err != nil {
		return false, err
	}

	var comparisonString string

	currentPath := "/inbox"

	switch headers {
	case "(request-target) host date digest":
		comparisonString = fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s", currentPath, header.Get("Host"), header.Get("Date"), header.Get("Digest"))
	case "(request-target) host date digest content-type":
		comparisonString = fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s\ncontent-type: %s", currentPath, header.Get("Host"), header.Get("Date"), header.Get("Digest"), header.Get("Content-Type"))
	default:
		logger.Warn("No header configuration found for", headers)
		return false, nil
	}

	logger.Debug("ComparisonString=", comparisonString)

	hashed := sha256.Sum256([]byte(comparisonString))

	rsaKey := importPem(publicKey.PublicKeyPem)

	rsaError := rsa.VerifyPKCS1v15(&rsaKey, crypto.SHA256, hashed[:], decoded)

	if rsaError != nil {
		return false, rsaError
	} else {
		return true, nil
	}
}

func importPem(pemPublicKey string) rsa.PublicKey {
	// Decode the PEM data
	block, _ := pem.Decode([]byte(pemPublicKey))
	if block == nil {
		log.Fatalf("Failed to decode PEM data")
	}

	// Parse the RSA public key
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse RSA public key: %v", err)
	}

	return *pubKey
}

func GetPublicKeyPem(privateKey rsa.PrivateKey) (*string, error) {
	// Extract public key from private key
	publicKey := &privateKey.PublicKey

	// Marshal the public key to DER format
	pubDER := x509.MarshalPKCS1PublicKey(publicKey)

	// Create PEM block for the public key
	pubBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	}

	// Encode the PEM block to string
	pubPEM := string(pem.EncodeToMemory(pubBlock))

	return &pubPEM, nil
}

func GetPrivateKeyPem(privateKey rsa.PrivateKey) (string, error) {
	// Convert RSA private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(&privateKey),
	}

	privateKeyPEMString := string(pem.EncodeToMemory(privateKeyPEM))

	return privateKeyPEMString, nil
}
