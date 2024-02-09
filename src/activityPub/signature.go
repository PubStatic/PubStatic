package activityPub

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/url"
	"strings"
)

func validateSignature(header map[string][]string, publicKey PublicKey) (bool, error) {

	signature := header["Signature"][0]

	parts := strings.Split(signature, ",")
	for _, part := range parts {
		fmt.Println("Signature Header Part=", part)
	}

	var keyIDString string
	for _, part := range parts {
		if strings.HasPrefix(part, "keyId") {
			keyIDString = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(part, "keyId=", ""), "\"", ""), "#main-key", "")
		}
	}

	if keyIDString == "" {
		fmt.Println("keyIdString is NullOrEmpty")
		return false, nil
	}

	keyID, err := url.Parse(keyIDString)
	if err != nil {
		fmt.Println("Error parsing keyId:", err)
		return false, err
	}

	signatureHash := ""
	headers := ""
	for _, part := range parts {
		if strings.HasPrefix(part, "signature") {
			signatureHash = strings.ReplaceAll(strings.ReplaceAll(part, "signature=", ""), "\"", "")
		} else if strings.HasPrefix(part, "headers") {
			headers = strings.ReplaceAll(strings.ReplaceAll(part, "headers=", ""), "\"", "")
		}
	}

	fmt.Println("KeyId=", keyID)

	var comparisionString string

	currentPath := "/inbox"

	switch headers {
	case "(request-target) host date digest":
		comparisionString = fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s", currentPath, header["Host"][0], header["Date"][0], header["Digest"][0])
	case "(request-target) host date digest content-type":
		comparisionString = fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s\ncontent-type: %s", currentPath, header["Host"][0], header["Date"][0], header["Digest"][0], header["Content-Type"][0])
	default:
		fmt.Println("No header configuration found for", headers)
		return false, nil
	}

	fmt.Println("ComparisonString=", comparisionString)

	hashed := sha256.Sum256([]byte(comparisionString))

	rsaKey := importPem(publicKey.PublicKeyPem)

	rsa.VerifyPKCS1v15(&rsaKey, crypto.SHA256, hashed[:], []byte(signatureHash))

	return true, nil
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
