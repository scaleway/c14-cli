package sis

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createSession(s3AccessKey string, s3SecretKey string, s3Profile string) *s3.S3 {

	s3Config := &aws.Config{
		Endpoint: aws.String("https://s3.fr-par.scw.cloud"),
		Region:   aws.String("fr-par"),
	}

	if s3Profile == "" {
		s3Config.Credentials = credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, "")
	}

	newSession := session.New(s3Config)

	return s3.New(newSession)
}

// CheckAPI : veryfy that we can issue call to the S3 API
func CheckAPI(s3AccessKey string, s3SecretKey string, s3Profile string) (err error) {

	fmt.Println("Checking S3 API credentials...")

	s3Client := createSession(s3AccessKey, s3SecretKey, s3Profile)

	// Sample call to check API credenetials
	_, err = s3Client.ListBuckets(nil)
	return
}

// CheckBucket : verify if destination migration bucket exists
func CheckBucket(bucketName string, s3AccessKey string, s3SecretKey string, s3Profile string) (bucketExists bool, err error) {

	s3Client := createSession(s3AccessKey, s3SecretKey, s3Profile)

	// Sample call to check API credenetials
	result, err := s3Client.ListBuckets(nil)

	bucketExists = false
	for _, b := range result.Buckets {
		if aws.StringValue(b.Name) == bucketName {
			bucketExists = true
		}
	}

	return
}

// CreateBucket : create S3 bucket for migration
func CreateBucket(bucketName string, s3AccessKey string, s3SecretKey string, s3Profile string) (err error) {

	s3Client := createSession(s3AccessKey, s3SecretKey, s3Profile)

	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = s3Client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return
	}

	fmt.Printf("Bucket %q successfully created\n", bucketName)
	return
}
