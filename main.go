package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

var bucket = os.Getenv("BUCKET")
var object = os.Getenv("OBJECT")
var serviceAccount = os.Getenv("SERVICE_ACCOUNT")

func main() {
	u, _ := generateV4GetObjectSignedURL(bucket, object, serviceAccount)
	fmt.Println(u)
}

// generateV4GetObjectSignedURL generates object signed URL with GET method.
func generateV4GetObjectSignedURL(bucket, object, serviceAccount string) (string, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccount := "service_account.json"
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}
	u, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}

	log.Printf("Generated GET signed URL:")
	log.Printf("%q\n", u)
	log.Printf("You can use this URL with any user agent, for example:")
	log.Printf("curl %q\n", u)
	return u, nil
}
