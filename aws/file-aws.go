package aws

import (
	"GoDrive/config"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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
func UploadToAWS(dir string, hash string, filename string) (bool, error) {
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
		ContentDisposition:   aws.String("attachment;filename=\"" + filename + "\""),
		ServerSideEncryption: aws.String("AES256"),
	})
	return true, nil
}

// InitAWSMpUpload : init multipart uploading to S3
func InitAWSMpUpload(filehash string, filename string) string {
	// sess := GetSession()
	// svc := s3.New(sess)
	input := &s3.CreateMultipartUploadInput{
		Bucket:             aws.String(AWSS3Bucket),
		Key:                aws.String(filehash),
		ContentDisposition: aws.String("attachment;filename=\"" + filename + "\""),
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

// UploadChunkToAws : upload file chunks to AWS
func UploadChunkToAws(content io.Reader, filehash string, idx int64, uploadId string) {
	// sess := GetSession()
	// svc := s3.New(sess)
	log.Printf("Uplaoding part %d to aws\n", idx)
	input := &s3.UploadPartInput{
		Body:       aws.ReadSeekCloser(content),
		Bucket:     aws.String(AWSS3Bucket),
		Key:        aws.String(filehash),
		PartNumber: aws.Int64(idx),
		UploadId:   aws.String(uploadId),
	}

	_, err := svc.UploadPart(input)
	log.Printf("Uploading part %d DONE\n", idx)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				panic(aerr.Error())
			}
		} else {
			panic(err.Error())
		}
	}
}

// CompleteAWSPartUpload : complete the upload
func CompleteAWSPartUpload(filehash string, uploadId string) {
	// sess := GetSession()
	// svc := s3.New(sess)
	listInput := &s3.ListPartsInput{
		Bucket:   aws.String(AWSS3Bucket),
		Key:      aws.String(filehash),
		UploadId: aws.String(uploadId),
	}

	listResult, err := svc.ListParts(listInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				panic(aerr.Error())
			}
		} else {
			panic(err.Error())
		}
	}
	parts := listResult.Parts
	var completeParts []*s3.CompletedPart
	for _, p := range parts {
		completeParts = append(completeParts, &s3.CompletedPart{
			ETag:       p.ETag,
			PartNumber: p.PartNumber,
		})
	}
	input := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(AWSS3Bucket),
		Key:      &filehash,
		UploadId: &uploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completeParts,
		},
	}

	result, err := svc.CompleteMultipartUpload(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				panic(aerr.Error())
			}
		} else {
			panic(err.Error())
		}
	}

	fmt.Printf("%v\n", result)
}

// GetPartList : get the part list
func GetPartList(filehash string, uploadId string) []int {
	// sess := GetSession()
	// svc := s3.New(sess)
	input := &s3.ListPartsInput{
		Bucket:   aws.String(AWSS3Bucket),
		Key:      aws.String(filehash),
		UploadId: aws.String(uploadId),
	}

	result, err := svc.ListParts(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				panic(aerr.Error())
			}
		} else {
			panic(err.Error())
		}
	}

	idxList := make([]int, 0)
	for _, part := range result.Parts {
		idxList = append(idxList, int(*part.PartNumber))
	}
	log.Printf("uploaded ID LIst: %v\n", idxList)
	return idxList
}

// GetDownloadURL : get a temporary download signed url by S3
func GetDownloadURL(filehash string, filename string) string {
	// sess := GetSession()
	// svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(AWSS3Bucket),
		Key:                        aws.String(filehash),
		ResponseContentDisposition: aws.String("attachment;filename=\"" + filename + "\""),
	})
	urlStr, err := req.Presign(24 * time.Hour)

	if err != nil {
		panic(err)
	}

	return urlStr
}

//DeleteFromAWS removes file from bucket
func DeleteFromAWS(filehash string) {

	fileToBeDeleted := &s3.DeleteObjectInput{
		Bucket: aws.String(AWSS3Bucket),
		Key:    aws.String(filehash),
	}

	result, err := svc.DeleteObject(fileToBeDeleted)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

}
