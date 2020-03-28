package aws

import (
	"GoDrive/config"
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//DownloadFromAWS takes hash from db and returns file downloaded from aws bucket
func DownloadFromAWS(hash string, fileName string) (bool, error) {
	path := config.WholeFileStoreLocation + fileName
	file, err := os.Create(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	sess := GetSession()
	downloader := s3manager.NewDownloader(sess)

	fileBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(AWSS3Bucket),
		Key:    aws.String(hash),
	})
	if err != nil {
		return false, err
	}
	fmt.Println("Downloaded: ", file.Name(), fileBytes, " bytes")
	return true, nil
}

//UploadToAWS uploads file to aws
func UploadToAWS(dir string, hash string) (bool, error) {
	sess := GetSession()
	file, err := os.Open(dir)
	if err != nil {
		return false, err
	}

	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings for object
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(AWSS3Bucket),
		Key:                  aws.String(hash),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return true, nil
}

// InitAWSMpUpload : init multipart uploading to S3
func InitAWSMpUpload(filehash string) string {
	sess := GetSession()
	svc := s3.New(sess)
	input := &s3.CreateMultipartUploadInput{
		Bucket:            aws.String(AWSS3Bucket),
		Key:               aws.String(filehash),
		SSECustomerKeyMD5: aws.String(filehash),
	}

	result, err := svc.CreateMultipartUpload(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
				panic(aerr.Error())
			}
		} else {
			panic(err.Error())
		}
	}

	return aws.StringValue(result.UploadId)
}
