package aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

//aws configs
const (
	AWSS3Region string = "us-east-1"
	AWSS3Bucket string = "godrive-bucket"
)

var sess *session.Session

func init() {
	awss, err := session.NewSession(
		&aws.Config{
			Region: aws.String(AWSS3Region),
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to connect to aws bucket")
		os.Exit(1)
	}
	sess = awss
}

//GetSession returns a aws s3 session
func GetSession() *session.Session {
	return sess
}
