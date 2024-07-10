package main

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.Background()

	accessKey := "<access-Key>"
	secretKey := "<secret-Key>"
	cred := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	endpoint := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "http://localhost:9000",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(cred), config.WithEndpointResolver(endpoint))
	if err != nil {
		log.Fatalln(err)
	}

	// change object address style
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	// get buckets
	lbo, err := client.ListBuckets(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	buckets := make(map[string]struct{}, len(lbo.Buckets))
	for _, b := range lbo.Buckets {
		buckets[*b.Name] = struct{}{}
	}

	// create 'develop' bucket if not exist
	bucketName := "develop"
	if _, ok := buckets[bucketName]; !ok {
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: &bucketName,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	// put object
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    aws.String("hogehoge"),
		Body:   bytes.NewReader([]byte("Hello, MinIO!")),
	})
	if err != nil {
		log.Fatalln(err)
	}
}
