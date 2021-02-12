/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awstasks

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"k8s.io/klog"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
	"k8s.io/kops/upup/pkg/fi/cloudup/cloudformation"
	"k8s.io/kops/upup/pkg/fi/cloudup/terraform"
)

// +kops:fitask
type SQS struct {
	Name      *string
	Lifecycle *fi.Lifecycle

	Tags map[string]string
}

var _ fi.CompareWithID = &SQS{}

func (s *SQS) CompareWithID() *string {
	return s.Name
}

func (s *SQS) Find(c *fi.Context) (*SQS, error) {
	klog.Warning("JASON in sqs Find")
	//
	//cloud := c.Cloud.(awsup.AWSCloud)
	//
	//request := &sqs.ListQueuesInput{
	//	MaxResults:      aws.Int64(2),
	//	NextToken:       nil,
	//	QueueNamePrefix: s.Name,
	//}
	//
	//response, err := cloud.SQS().ListQueues(request)
	//if err != nil {
	//	return nil, fmt.Errorf("error listing SQS queues: %v", err)
	//}
	//if response == nil || len(response.QueueUrls) == 0 {
	//	return nil, nil
	//}
	//
	//if len(response.QueueUrls) != 1 {
	//	return nil, fmt.Errorf("found multiple SQS queues matching queue name")
	//}
	
	//q := response.QueueUrls[0]

	//{
	//  "MessageRetentionPeriod": "300",
	//  "Policy": "$(echo $QUEUE_POLICY | sed 's/\"/\\"/g')"
	//}
	//
	//    * MessageRetentionPeriod – Returns the length of time, in seconds, for
	//    which Amazon SQS retains a message.
	//
	//    * Policy – Returns the policy of the queue.

	//all := "ALL"
	//t, err := cloud.SQS().GetQueueAttributes(&sqs.GetQueueAttributesInput{
	//	AttributeNames: []*string{&all},
	//	QueueUrl:       q,
	//})
	//
	//actual := &SQS{
	//	Name:      t.,
	//	Lifecycle: nil,
	//	Tags:      nil,
	//}

	return nil, nil
}

func (s *SQS) Run(c *fi.Context) error {
	// DefaultDeltaRunMethod implements the standard change-based run procedure:
	// find the existing item; compare properties; call render with (actual, expected, changes)
	return fi.DefaultDeltaRunMethod(s, c)
}

func (_ *SQS) CheckChanges(a, e, changes *SQS) error {
	return nil
}

func (s *SQS) RenderAWS(t *awsup.AWSAPITarget, a, e, changes *SQS) error {
	klog.Warning("JASON in sqs RenderAWS")
	// thing is new, hasn't been created yet
	if a == nil {
		cloud := t.Cloud
		accountId, _, _ := cloud.AccountInfo()
		region := cloud.Region()
		// build the policy
		policy, err := s.buildQueuePolicy(accountId, region)
		if err != nil {
			return err
		}
		// turn struct object into json
		p, err := json.Marshal(policy)
		if err != nil {
			return err
		}

		request := &sqs.CreateQueueInput{
			Attributes: map[string]*string{
				"DelaySeconds":           aws.String("60"),
				"MessageRetentionPeriod": aws.String("86400"),
				"Policy": aws.String(string(p)),
			},
			QueueName:  s.Name,
			Tags:       nil,
		}
		klog.Warningf("The policy: %v", string(p))
		_, err = t.Cloud.SQS().CreateQueue(request)
		if err != nil {
			return fmt.Errorf("error creating SQS queue: %v", err)
		}

	}

	// We don't tag the zone - we expect it to be shared
	return nil
}

func (_ *SQS) RenderTerraform(t *terraform.TerraformTarget, a, e, changes *SQS) error {
	return nil
}

func (_ *SQS) TerraformLink() *terraform.Literal {
	return nil
}

func (_ *SQS) RenderCloudformation(t *cloudformation.CloudformationTarget, a, e, changes *SQS) error {
	return nil
}

func (_ *SQS) CloudformationLink() *cloudformation.Literal {
	return nil
}

type Principal struct {
	Service []string
}

type Statement struct {
	Effect string
	Principal Principal
	Action string
	Resource string
}

type Policy struct {
	Version string
	Statement Statement
}

func (q *SQS) buildQueuePolicy(accountId, region string) (Policy, error) {
	klog.Warning("JASON In buildQueuePolicy")
	//{
	//    "Version": "2012-10-17",
	//    "Id": "MyQueuePolicy",
	//    "Statement": [{
	//        "Effect": "Allow",
	//        "Principal": {
	//            "Service": ["events.amazonaws.com", "sqs.amazonaws.com"]
	//        },
	//        "Action": "sqs:SendMessage",
	//        "Resource": [
	//            "arn:aws:sqs:${AWS_REGION}:${ACCOUNT_ID}:${SQS_QUEUE_NAME}"
	//        ]
	//    }]
	//}
	statement := Statement{
		Effect:    "Allow",
		Principal: Principal{
			Service:   []string{"events.amazonaws.com", "sqs.amazonaws.com"},
		},
		Action: "sqs:SendMessage",
		Resource: "arn:aws:sqs:" + region + ":" + accountId + ":" + *q.Name,
	}

	policy := Policy{
		Version: "2012-10-17",
		Statement: statement,
	}

	return policy, nil
}

//"Policy": "{\"Version\":\"2012-10-17\",\"Id\":\"MyQueuePolicy\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"Service\":[\"sqs.amazonaws.com\",\"events.amazonaws.com\"]},\"Action\":\"sqs:SendMessage\",\"Resource\":\"arn:aws:sqs:us-east-2:259072402577:testQueueName\"}]}",
//"Policy": "{\"Version\":\"2012-10-17\",\"Statement\":{\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"[\\\"events.amazonaws.com\\\", \\\"sqs.amazonaws.com\\\"]\"},\"Action\":\"sqs:SendMessage\",\"Resource\":\"arn:aws:sqs:us-east-2:953943258687:myFirstQueuePlsWork\"}}"