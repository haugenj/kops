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

package awsmodel

import (
	"k8s.io/klog"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awstasks"
	"strings"
)

//  todo SQSBuilder builds an SQS boiiiiiii
// from apply_cluster.go
type SQSBuilder struct {
	*AWSModelContext
	QueueName string

	Lifecycle         *fi.Lifecycle
}

var _ fi.ModelBuilder = &SQSBuilder{}

func (b *SQSBuilder) Build(c *fi.ModelBuilderContext) error {
	klog.Warning("JASON In the SQS build")
	clusterName := b.AWSModelContext.ClusterName()
	queueName := strings.Replace(clusterName, ".", "_", -1)
	// @step: now lets build the sqs task
	task := &awstasks.SQS{
		Name:           &queueName,
		Lifecycle:      b.Lifecycle,
		Tags:           b.CloudTags(b.QueueName, false),
	}

	c.AddTask(task)

	return nil
}
