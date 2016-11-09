package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	region = "ap-northeast-1"
)

func main() {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Fatal(err)
	}

	svc := ec2.New(sess)

	params := &ec2.DescribeInstanceStatusInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("event.code"),
				Values: []*string{
					aws.String("instance-reboot"),
					aws.String("instance-stop"),
					aws.String("instance-retirement"),
					aws.String("system-reboot"),
					aws.String("system-maintenance"),
				},
			},
		},
	}

	resp, err := svc.DescribeInstanceStatus(params)
	if err != nil {
		log.Fatal(err)
	}

	for _, instance := range resp.InstanceStatuses {
		for _, event := range instance.Events {
			fmt.Printf("%v: %v %v\n", *instance.InstanceId, *event.Code, *event.Description)
		}
	}
}
