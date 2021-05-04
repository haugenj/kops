/*
Copyright 2020 The Kubernetes Authors.

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

package model

import (
	"crypto/x509"
	"testing"

	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/pki"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/util/pkg/vfs"
)

// fakeKeyStore mocks out some of fi.KeyStore, for our tests.
type fakeKeyStore struct {
	T *testing.T
}

var _ fi.Keystore = &fakeKeyStore{}

func (k fakeKeyStore) FindPrimaryKeypair(name string) (*pki.Certificate, *pki.PrivateKey, error) {
	panic("fakeKeyStore does not implement FindPrimaryKeypair")
}

func (k fakeKeyStore) FindKeyset(name string) (*fi.Keyset, error) {
	panic("fakeKeyStore does not implement FindKeyset")
}

func (k fakeKeyStore) CreateKeypair(signer string, name string, template *x509.Certificate, privateKey *pki.PrivateKey) (*pki.Certificate, error) {
	panic("fakeKeyStore does not implement CreateKeypair")
}

func (k fakeKeyStore) StoreKeypair(id string, keyset *fi.Keyset) error {
	panic("fakeKeyStore does not implement StoreKeypair")
}

func (k fakeKeyStore) MirrorTo(basedir vfs.Path) error {
	panic("fakeKeyStore does not implement MirrorTo")
}

// fakeCAStore mocks out some of fi.CAStore, for our tests.
// Although CAStore currently embeds KeyStore, we maintain the split here in the hope we can clean this up in future.
type fakeCAStore struct {
	fakeKeyStore

	privateKeysets map[string]*kops.Keyset
	certs          map[string]*pki.Certificate
}

var _ fi.CAStore = &fakeCAStore{}

func (k fakeCAStore) FindCertificatePool(name string) (*fi.CertificatePool, error) {
	panic("fakeCAStore does not implement FindCertificatePool")
}

func (k fakeCAStore) FindCertificateKeyset(name string) (*kops.Keyset, error) {
	panic("fakeCAStore does not implement FindCertificateKeyset")
}

func (k fakeCAStore) FindPrivateKey(name string) (*pki.PrivateKey, error) {
	primaryId := k.privateKeysets[name].Spec.PrimaryId
	for _, item := range k.privateKeysets[name].Spec.Keys {
		if item.Id == primaryId {
			return pki.ParsePEMPrivateKey(item.PrivateMaterial)
		}
	}
	return nil, nil
}

func (k fakeCAStore) FindPrivateKeyset(name string) (*kops.Keyset, error) {
	return k.privateKeysets[name], nil
}

func (k fakeCAStore) FindCert(name string) (*pki.Certificate, error) {
	return k.certs[name], nil
}

func (k fakeCAStore) ListKeysets() ([]*kops.Keyset, error) {
	panic("fakeCAStore does not implement ListKeysets")
}

func (k fakeCAStore) DeleteKeysetItem(item *kops.Keyset, id string) error {
	panic("fakeCAStore does not implement DeleteKeysetItem")
}
