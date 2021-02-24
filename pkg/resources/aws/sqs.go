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

package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"k8s.io/klog/v2"
	"strings"

	"k8s.io/kops/pkg/resources"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
)


func DumpSQSQueue(op *resources.DumpOperation, r *resources.Resource) error {
	klog.Warningf("JASON DUMPING SQS QUEUES")
	data := make(map[string]interface{})
	data["id"] = r.ID
	data["name"] = r.Name
	data["type"] = r.Type
	data["raw"] = r.Obj
	op.Dump.Resources = append(op.Dump.Resources, data)

	return nil
}

func DeleteSQSQueue(cloud fi.Cloud, r *resources.Resource) error {
	klog.Warningf("JASON DELETING SQS QUEUES")
	c := cloud.(awsup.AWSCloud)

	url := r.ID

	klog.V(2).Infof("Deleting SQS Queue %q", url)
	request := &sqs.DeleteQueueInput{
		QueueUrl: &url,
	}
	_, err := c.SQS().DeleteQueue(request)
	if err != nil {
		return fmt.Errorf("error deleting SQS Queue %q: %v", url, err)
	}
	return nil
}


func ListSQSQueues(cloud fi.Cloud, clusterName string) ([]*resources.Resource, error) {
	klog.Warningf("JASON LISTING SQS QUEUES")
	c := cloud.(awsup.AWSCloud)

	klog.V(2).Infof("Listing SQS queues")
	queueName := strings.Replace(clusterName, ".", "_", -1)

	request := &sqs.ListQueuesInput{
		NextToken:       nil,
		QueueNamePrefix: &queueName,
	}
	klog.Warningf("the request: %v", request)
	response, err := c.SQS().ListQueues(request)
	if err != nil {
		return nil, fmt.Errorf("error listing SQS queues: %v", err)
	}
	if response == nil || len(response.QueueUrls) == 0 {
		return nil, nil
	}

	if len(response.QueueUrls) != 1 {
		return nil, fmt.Errorf("found multiple SQS queues matching queue name")
	}
	queue := response.QueueUrls[0]

	var resourceTrackers []*resources.Resource
	resourceTracker := &resources.Resource{
			Name: 	 queueName,
			ID:      *queue,
			Type:    "sqs",
			Deleter: DeleteSQSQueue,
			Dumper:  DumpSQSQueue,
			Obj:     queue,
		}

	resourceTrackers = append(resourceTrackers, resourceTracker)

	return resourceTrackers, nil
}
