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
	URL 		*string
	Lifecycle *fi.Lifecycle

	Tags map[string]string
}

var _ fi.CompareWithID = &SQS{}

func (q *SQS) CompareWithID() *string {
	return q.URL
}

func (q *SQS) Find(c *fi.Context) (*SQS, error) {
	klog.Warning("JASON in sqs Find")

	cloud := c.Cloud.(awsup.AWSCloud)

	request := &sqs.ListQueuesInput{
		MaxResults:      aws.Int64(2),
		NextToken:       nil,
		QueueNamePrefix: q.Name,
	}
	klog.Warningf("the request: %v", request)
	response, err := cloud.SQS().ListQueues(request)
	if err != nil {
		return nil, fmt.Errorf("error listing SQS queues: %v", err)
	}
	klog.Warningf("JASON the find response: %v", response)
	if response == nil || len(response.QueueUrls) == 0 {
		return nil, nil
	}

	if len(response.QueueUrls) != 1 {
		return nil, fmt.Errorf("found multiple SQS queues matching queue name")
	}
	q.URL = response.QueueUrls[0]

	tags, err := cloud.SQS().ListQueueTags(&sqs.ListQueueTagsInput{
		QueueUrl: q.URL,
	})
	klog.Warningf("the found tags are: %v", tags)

	actual := &SQS{
		Name:   q.Name,
		URL: 	q.URL,
		Tags:   intersectSQSTags(tags.Tags, q.Tags),
		Lifecycle: q.Lifecycle,
	}
	klog.Warningf("JASON this would be the actual: %v", actual)
	return actual, nil
}

func (q *SQS) Run(c *fi.Context) error {
	// DefaultDeltaRunMethod implements the standard change-based run procedure:
	// find the existing item; compare properties; call render with (actual, expected, changes)
	return fi.DefaultDeltaRunMethod(q, c)
}

// TODO
func (q *SQS) CheckChanges(a, e, changes *SQS) error {
	return nil
}

func (q *SQS) RenderAWS(t *awsup.AWSAPITarget, a, e, changes *SQS) error {
	klog.Warning("JASON in sqs RenderAWS")
	// queue is new, hasn't been created yet
	if a == nil {
		cloud := t.Cloud
		accountId, _, _ := cloud.AccountInfo()
		region := cloud.Region()
		// build the policy
		policy, err := q.buildQueuePolicy(accountId, region)
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
				"MessageRetentionPeriod": aws.String("300"),
				"Policy": aws.String(string(p)),
			},
			QueueName: q.Name,
			Tags:      convertTagsToPointers(q.Tags),
		}
		klog.Warningf("The policy: %v", string(p))
		response, err := t.Cloud.SQS().CreateQueue(request)
		if err != nil {
			return fmt.Errorf("error creating SQS queue: %v", err)
		}
		q.URL = response.QueueUrl
	}

	return nil
}

// change tags to format required by CreateQueue
func convertTagsToPointers(tags map[string]string) map[string]*string {
	newTags := map[string]*string{}
	for k, v := range tags {
		vv := v
		newTags[k] = &vv
	}

	return newTags
}
// TODO
func (q *SQS) RenderTerraform(t *terraform.TerraformTarget, a, e, changes *SQS) error {
	return nil
}
// TODO
func (q *SQS) TerraformLink() *terraform.Literal {
	return nil
}
// TODO
func (q *SQS) RenderCloudformation(t *cloudformation.CloudformationTarget, a, e, changes *SQS) error {
	return nil
}
// TODO
func (q *SQS) CloudformationLink() *cloudformation.Literal {
	return nil
}

// Todo see if we can leverage existing IAM stuff instead of making these structs here.
// when I tried initially to do so I got a cyclic import :/
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