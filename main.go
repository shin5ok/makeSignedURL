package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
)

var apiUrl = "https://api.uname.link/slack"

// var serviceAccount = os.Getenv("SERVICE_ACCOUNT")

// var slackChannel = os.Getenv("SLACK_CHANNEL")
// var slackURL = os.Getenv("SLACK_URL")

func main() {
	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		// storage.googleapis.com/projects/_/buckets/kawanos-dists/objects/foo/bar/services
		subject := c.Request.Header.Get("ce-subject")
		t := strings.Split(subject, "/")
		bucket := t[4]
		object := strings.Join(t[6:], "/")
		// not matching pattern then just return empty
		log.Printf("bucket:%s, object:%s, subject:%s\n", bucket, object, subject)
		r, err := generateV4GetObjectSignedURL(bucket, object)
		if err != nil {
			fmt.Println(err)
		}
		path := fmt.Sprintf("gs://%s/%s", bucket, object)
		result := fmt.Sprintf("SignURL: %q\nExpire: %s", r.SignedURL, r.Expire)
		notifySlack(result)
		// log.Println(u)
		c.JSON(http.StatusOK, gin.H{"Path": path, "result": result})
	})

	r.Run(":8080")
}

func notifySlack(message string) error {
	params := url.Values{}
	params.Add("message", message)
	http.PostForm(apiUrl, params)
	return nil
}

type ResultSignedURL struct {
	SignedURL string
	Expire    time.Time
}

// generateV4GetObjectSignedURL generates object signed URL with GET method.
func generateV4GetObjectSignedURL(bucket, object string) (ResultSignedURL, error) {
	log.Println("start generate")
	ctx := context.Background()
	// bucket := "bucket-name"
	// object := "object-name"
	// serviceAccount := "service_account.json"

	// jsonKey, err := ioutil.ReadFile(serviceAccount)
	// if err != nil {
	// 	return "", fmt.Errorf("ioutil.ReadFile: %v", err)
	// }

	creds, err := google.FindDefaultCredentials(ctx, storage.ScopeReadOnly)
	if err != nil {
		log.Printf("------\n1.%s--------\n", err)
	}
	// ComputeTokenSource
	log.Println(creds.JSON)
	conf, err := google.JWTConfigFromJSON(creds.JSON, storage.ScopeReadOnly)
	if err != nil {
		log.Printf("------\n2.%s--------\n", err)
		return ResultSignedURL{}, fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	expire := time.Now().Add(24 * 7 * time.Hour)
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        expire,
	}
	fmt.Println(opts)
	u, err := storage.SignedURL(bucket, object, opts)
	if err != nil {
		return ResultSignedURL{}, fmt.Errorf("storage.SignedURL: %v", err)
	}

	log.Printf("Generated GET signed URL:")
	log.Printf("%q\n", u)
	log.Printf("You can use this URL with any user agent, for example:")
	log.Printf("curl %q\n", u)
	log.Printf("Expire: %s", expire)
	return ResultSignedURL{u, expire}, nil
}
