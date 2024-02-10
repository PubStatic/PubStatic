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
	"strings"
)

func validateSignature(header map[string][]string, publicKey PublicKey) (bool, error) {
	logger.Trace("Validating signature")

	signature := header["Signature"][0]

	parts := strings.Split(signature, ",")
	for _, part := range parts {
		fmt.Println("Signature Header Part=", part)
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
		comparisonString = fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s", currentPath, header["Host"][0], header["Date"][0], header["Digest"][0])
	case "(request-target) host date digest content-type":
		comparisonString = fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: %s\ncontent-type: %s", currentPath, header["Host"][0], header["Date"][0], header["Digest"][0], header["Content-Type"][0])
	default:
		fmt.Println("No header configuration found for", headers)
		return false, nil
	}

	fmt.Println("ComparisonString=", comparisonString)

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
