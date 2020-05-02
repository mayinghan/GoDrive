package storage

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var cephSess *session.Session
var cephSvc *s3.S3

func init() {

}
