package storage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//aws configs
const (
	AWSS3Region string = "us-east-2"
	AWSS3Bucket string = "godrive"
)

var sess *session.Session
var svc *s3.S3

func init() {
	awss, err := session.NewSession(
		&aws.Config{
			Region:          aws.String(AWSS3Region),
			Credentials:     credentials.NewSharedCredentials("", "yinghan"),
			Endpoint:        aws.String("s3-accelerate.amazonaws.com"),
			S3UseAccelerate: aws.Bool(true),
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to connect to aws bucket")
		os.Exit(1)
	}
	sess = awss
	svc = s3.New(sess)
}

// //GetSession returns a aws s3 session
// func GetSession() *session.Session {
// 	return sess
// }
