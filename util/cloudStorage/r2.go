package cloudStorage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"net/http"
)

type CloudflareR2Storage struct {
	AccountId       string
	BucketName      string
	AccessKeyId     string
	AccessKeySecret string
}

func NewCloudflareR2Storage(config CloudflareR2Storage) *CloudflareR2Storage {
	return &CloudflareR2Storage{
		AccountId:       config.AccountId,
		BucketName:      config.BucketName,
		AccessKeyId:     config.AccessKeyId,
		AccessKeySecret: config.AccessKeySecret,
	}
}

func (c *CloudflareR2Storage) R2Init() (*s3.Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", c.AccountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.AccessKeyId, c.AccessKeySecret, "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	return client, nil
}

func (c *CloudflareR2Storage) Upload(r2Client *s3.Client, fileName string, fileType string, fileData multipart.File, ctx *fiber.Ctx) (*string, error) {

	// TODO: a function to check if bucket exists

	_, err := r2Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.BucketName),
		Key:         aws.String(fileName),
		ContentType: aws.String(fileType),
		Body:        fileData,
	})
	if err != nil {
		return nil, ctx.Status(http.StatusInternalServerError).JSON(err)
	}

	// TODO: insert url variable on frontend when configuring the cloud stroage providers variable on the frontend
	fileUrl := aws.String(fmt.Sprintf("https://pub-2983768645604ff7b7f5947cea0a55a9.r2.dev/%s", fileName))

	return fileUrl, nil
}
