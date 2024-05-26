package s3service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Service provides methods for uploading files to Amazon S3
type S3Service struct {
	BucketName string
	Region     string
	AccessKey  string
	SecretKey  string
	// SessionToken string
}

// NewS3Service creates a new instance of S3Service
func NewS3Service(bucketName, region, accessKey, secretKey string) *S3Service {
	return &S3Service{
		BucketName: bucketName,
		Region:     region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		// SessionToken: sessionToken,
	}
}

// UploadFileToS3 uploads a file to Amazon S3 and returns the key
func (s3Service *S3Service) UploadFileToS3(file io.Reader, folderPath, filename string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Service.Region),
		Credentials: credentials.NewStaticCredentials(s3Service.AccessKey, s3Service.SecretKey, ""),
	})
	if err != nil {
		return "", fmt.Errorf("error creating AWS session: %s", err)
	}

	svc := s3.New(sess)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("unable to read file: %s", err)
	}

	key := filepath.Join(folderPath, filename)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s3Service.BucketName),
		Key:           aws.String(key),
		ACL:           aws.String("private"),
		Body:          bytes.NewReader(fileBytes),
		ContentLength: aws.Int64(int64(len(fileBytes))),
		ContentType:   aws.String(http.DetectContentType(fileBytes)),
	})
	if err != nil {
		return "", fmt.Errorf("unable to upload file to S3: %s", err)
	}

	// fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3Service.BucketName, s3Service.Region, key)

	return key, nil
}

func (s3Service *S3Service) GetPreSignedURL(key string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Service.Region),
		Credentials: credentials.NewStaticCredentials(s3Service.AccessKey, s3Service.SecretKey, ""),
	})

	if err != nil {
		return "", fmt.Errorf("error creating AWS session: %s", err)
	}

	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s3Service.BucketName),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to sign request: %s", err)
	}

	return urlStr, err
}
func (s3Service *S3Service) DeleteKey(key string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Service.Region),
		Credentials: credentials.NewStaticCredentials(s3Service.AccessKey, s3Service.SecretKey, ""),
	})

	if err != nil {
		return fmt.Errorf("error creating AWS session: %s", err)
	}

	// Create S3 service client
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3Service.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to sign request: %s", err)
	}

	return err
}
