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

func printInstanceName(svc *ec2.EC2, instanceId *string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{instanceId},
	}
	resp, err := svc.DescribeInstances(input)
	if err != nil {
		return "", err
	}
	tags := resp.Reservations[0].Instances[0].Tags
	for _, elem := range tags {
		if aws.StringValue(elem.Key) == "Name" {
			return aws.StringValue(elem.Value), nil
		}
	}
	return "", nil
}

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
		name, err := printInstanceName(svc, instance.InstanceId)
		if err != nil {
			log.Fatal(err)
		}
		for _, event := range instance.Events {
			fmt.Printf("%v", *instance.InstanceId)
			if name != "" {
				fmt.Printf(" (%s)", name)
			}
			fmt.Printf(": %v %v\n", *event.Code, *event.Description)
		}
	}
}
