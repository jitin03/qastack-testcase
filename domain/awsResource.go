package domain

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/gommon/log"

	"fmt"
)

const (
	AWS_S3_REGION = "us-west-2"
	AWS_S3_BUCKET = ""
)

func connectAWS() *session.Session {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(AWS_S3_REGION)})
	if err != nil {
		panic(err)
	}
	return sess
}

func CreateBucket(client *s3.S3, bucketName string) error {
	_, err := client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})

	return err
}

func CreateAWSBucket(sess *session.Session, bucketName string) error {
	// sess, err := session.NewSessionWithOptions(session.Options{
	// 	Profile: "default",
	// 	Config: aws.Config{
	// 		Region: aws.String("us-east-1"),
	// 	},
	// })

	// if err != nil {
	// 	fmt.Printf("Failed to initialize new session: %v", err)
	// 	return
	// }

	var err error
	s3Client := s3.New(sess)
	err = CreateBucket(s3Client, "qastack-test-results")
	if err != nil {
		fmt.Printf("Couldn't create bucket: %v", err)
		return err
	}

	fmt.Println("Successfully created bucket")
	return nil
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, fileDir string, S3_BUCKET string) error {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("qastack-test-results"),
		Key:                  aws.String(S3_BUCKET + "/" + fileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func DownloadFile(s *session.Session, bucketName string, key string) (string, error) {
	file, err := os.Create(key)
	if err != nil {
		return "", err
	}

	defer file.Close()

	// Create S3 service client
	svc := s3.New(s)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("qastack-test-results"),
		Key:    aws.String(bucketName + "/" + key),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Info("Failed to sign request", err)
	}

	log.Info("The URL is", urlStr)

	return urlStr, err
}
func ListItems(client *s3.S3, bucketName string, prefix string) (*s3.ListObjectsV2Output, error) {
	log.Info(bucketName)
	res, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String("qastack-test-results"),
		Prefix: aws.String(bucketName + "/" + prefix),
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
