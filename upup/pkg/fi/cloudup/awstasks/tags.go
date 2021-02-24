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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func mapEC2TagsToMap(tags []*ec2.Tag) map[string]string {
	if tags == nil {
		return nil
	}
	m := make(map[string]string)
	for _, t := range tags {
		if strings.HasPrefix(aws.StringValue(t.Key), "aws:cloudformation:") {
			continue
		}
		m[aws.StringValue(t.Key)] = aws.StringValue(t.Value)
	}
	return m
}

func findNameTag(tags []*ec2.Tag) *string {
	for _, tag := range tags {
		if aws.StringValue(tag.Key) == "Name" {
			return tag.Value
		}
	}
	return nil
}

// intersectTags returns the tags of interest from a specified list of AWS tags;
// because we only add tags, this set of tags of interest is the tags that occur in the desired set.
func intersectTags(tags []*ec2.Tag, desired map[string]string) map[string]string {
	if tags == nil {
		return nil
	}
	actual := make(map[string]string)
	for _, t := range tags {
		k := aws.StringValue(t.Key)
		v := aws.StringValue(t.Value)

		if _, found := desired[k]; found {
			actual[k] = v
		}
	}
	if len(actual) == 0 && desired == nil {
		// Avoid problems with comparison between nil & {}
		return nil
	}
	return actual
}

// intersectSQSTags does the same thing as intersectTags, but takes different input because SQS tags are listed differently
func intersectSQSTags(tags map[string]*string, desired map[string]string) map[string]string {
	if tags == nil {
		return nil
	}
	actual := make(map[string]string)
	for k, v := range tags {
		vv := aws.StringValue(v)

		if _, found := desired[k]; found {
			actual[k] = vv
		}
	}
	if len(actual) == 0 && desired == nil {
		// Avoid problems with comparison between nil & {}
		return nil
	}
	return actual
}