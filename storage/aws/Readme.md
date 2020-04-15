# AWS S3 Setup
Guide to setting up AWS S3 for GoDrive

## Step 1: Download AWS CLI
Use Amazon's user guide to install AWS CLI on [Windows](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-windows.html), [Mac](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-mac.html), or [Linux](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-linux.html).


## Step 2: Configure AWS and set up credentials
```bash
aws configure
```
Put the following:

AWS Access Key ID: AKIAISGDUVPZ4M6K7I5Q
AWS Secret Access Key: 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
Default region name: us-east-1
Default output format: json


## Step 3: Install AWS SDK for Go:
```bash
go get github.com/aws/aws-sdk-go
```
To use AWS in GoDrive, import the packages: 

```
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)
```

## Step 4: Using AWS S3
### Uploading to AWS S3 
#### Inputs:
**dir** - The location of the file allows AWS S3 to read the file's content into a buffer. 
**hash** - The file hash acts as a key that allows AWS S3 to retrieve the file.

#### Outputs: 
**bool** - Success indicator
**error** - Error message

```
func UploadToAWS(dir string, hash string) (bool, error)
```
### Downloading from AWS S3
### Inputs:
**fileName** - Downloads the file to the given directory 
**hash** - The file hash acts as a key that allows AWS S3 to retrieve the file.

### Outputs:
**bool** - Success indicator
**error** - Error message

```
func DownloadFromAWS(hash string, fileName string) (bool, error)
```

