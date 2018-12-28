package aws

// this is temporary and will be moved in to go-whosonfirst-aws when the dust settles

import (
	"github.com/aaronland/go-string/dsn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	"strings"
)

func NewSQSServiceFromString(str_dsn string) (*sqs.SQS, string, error) {

	dsn_map, err := dsn.StringToDSNWithKeys(str_dsn, "region", "credentials", "queue")

	if err != nil {
		return nil, "", err
	}

	sqs_creds, _ := dsn_map["credentials"]
	sqs_region, _ := dsn_map["region"]
	sqs_queue, _ := dsn_map["queue"]

	sess, err := session.NewSessionWithCredentials(sqs_creds, sqs_region)

	if err != nil {
		return nil, "", err
	}

	svc := sqs.New(sess)

	if !strings.HasPrefix(sqs_queue, "https://sqs") {

		rsp, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
			QueueName: aws.String(sqs_queue),
		})

		if err != nil {
			return nil, "", err
		}

		sqs_queue = *rsp.QueueUrl
	}

	return svc, sqs_queue, nil
}
