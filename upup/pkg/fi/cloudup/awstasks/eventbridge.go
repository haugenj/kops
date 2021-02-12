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
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
	"k8s.io/kops/upup/pkg/fi/cloudup/cloudformation"
	"k8s.io/kops/upup/pkg/fi/cloudup/terraform"
)

// +kops:fitask
type EventBridge struct {
	ID        *string
	Name      *string
	Lifecycle *fi.Lifecycle

	Tags map[string]string
}

var _ fi.CompareWithID = &EventBridge{}

func (e *EventBridge) CompareWithID() *string {
	return e.Name
}

func (e *EventBridge) Find(c *fi.Context) (*DNSZone, error) {
	return nil, nil
}

func (e *EventBridge) Run(c *fi.Context) error {
	return nil
}

func (_ *EventBridge) CheckChanges(a, e, changes *EventBridge) error {
	return nil
}

func (_*EventBridge) RenderAWS(t *awsup.AWSAPITarget, a, e, changes *EventBridge) error {
	return nil
}

func (_ *EventBridge) RenderTerraform(t *terraform.TerraformTarget, a, e, changes *EventBridge) error {
	return nil
}

func (_ *EventBridge) TerraformLink() *terraform.Literal {
	return nil
}

func (_ *EventBridge) RenderCloudformation(t *cloudformation.CloudformationTarget, a, e, changes *EventBridge) error {
	return nil
}

func (_ *EventBridge) CloudformationLink() *cloudformation.Literal {
	return nil
}