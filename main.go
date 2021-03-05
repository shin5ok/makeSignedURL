package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
)

var bucket = os.Getenv("BUCKET")
var object = os.Getenv("OBJECT")
var serviceAccount = os.Getenv("SERVICE_ACCOUNT")

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		// storage.googleapis.com/projects/_/buckets/kawanos-dists/objects/foo/bar/services
		subject := c.Request.Header.Get("ce-subject")
		t := strings.Split(subject, "/")
		bucket := t[4]
		object := strings.Join(t[6:], "/")
		u, _ := generateV4GetObjectSignedURL(bucket, object, serviceAccount)
		path := fmt.Sprintf("gs://%s/%s", bucket, object)
		log.Println(u)
		c.JSON(http.StatusOK, gin.H{"Path": path})
	})

	r.Run(":8080")
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
