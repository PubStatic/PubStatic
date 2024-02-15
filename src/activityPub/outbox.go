package activityPub

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/PubStatic/PubStatic/models"
	"github.com/PubStatic/PubStatic/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func SentActivity(activity Activity, inbox url.URL, ownHost string, mongoDbConnectionString string) error {
	logger.Trace("Sending activity")

	// Create digest
	jsonBody, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	digest := computeHash(string(jsonBody))

	// Load private key
	privateKey, err := LoadPrivateKey(mongoDbConnectionString)
	if err != nil {
		logger.Warn("Error loading private key:", err)

		return err
	}

	// Create signature
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

	signedString := fmt.Sprintf("(request-target): post %s\nhost: %s\ndate: %s\ndigest: sha-256=%s", inbox.Path, inbox.Host, date, digest)

	sha := sha256.New()
	sha.Write([]byte(signedString))
	hashedBytes := sha.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashedBytes)
	if err != nil {

		return err
	}

	signatureString := base64.StdEncoding.EncodeToString(signature)

	// Create request
	req, err := http.NewRequest("POST", inbox.String(), bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Host", inbox.Host)
	req.Header.Set("Date", date)
	req.Header.Set("Digest", "sha-256="+digest)
	req.Header.Set("Signature", fmt.Sprintf("keyId=\"%s\",headers=\"(request-target) host date digest\",signature=\"%s\"",
		fmt.Sprintf("https://%s#main-key", ownHost), signatureString))
	req.Header.Add("Accept", "application/json")

	// Handle response
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {

		defer response.Body.Close()

		// Read the response body
		body, err := io.ReadAll(io.Reader(response.Body))
		if err != nil {
			logger.Warn("Reading body failed after failed request")

			return err
		}

		logger.Debugf("Body of failed request: %s", body)

		return fmt.Errorf("sending activity failed with error code: %d", response.StatusCode)
	}

	logger.Trace("Successfully send activity")

	return nil
}

func LoadPrivateKey(mongoConnectionString string) (*rsa.PrivateKey, error) {
	fileContentPrivateKey, err := repository.ReadMongo[models.KeyValue]("Actor", "Key", bson.D{{"key", "privateKey"}}, mongoConnectionString)

	if err != nil {
		return nil, err
	}

	privateKey, err := ImportPrivateKeyPem(fileContentPrivateKey.Value)

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func computeHash(jsonData string) string {
	sha := sha256.New()

	sha.Write([]byte(jsonData))
	hashedBytes := sha.Sum(nil)

	hashedString := base64.StdEncoding.EncodeToString(hashedBytes)

	return hashedString
}
