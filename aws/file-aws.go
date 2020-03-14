package aws

import (
	"bytes"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//UploadToAWS uploads file to aws
func UploadToAWS(sess *session.Session, dir string) (bool, error) {

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
		Key:                  aws.String(dir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return true, nil
}
