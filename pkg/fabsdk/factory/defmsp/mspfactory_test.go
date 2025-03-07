/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package defmsp

import (
	"path/filepath"
	"testing"

	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/test/mockmsp"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/core/config"
	cryptosuiteImpl "github.com/VRamakrishna/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/fab"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/fab/mocks"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/fabsdk/factory/defcore"
	mspimpl "github.com/VRamakrishna/fabric-sdk-go/pkg/msp"
	"github.com/VRamakrishna/fabric-sdk-go/test/metadata"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserStore(t *testing.T) {
	factory := NewProviderFactory()

	config := mocks.NewMockIdentityConfig()

	userStore, err := factory.CreateUserStore(config)
	if err != nil {
		t.Fatalf("Unexpected error creating state store %s", err)
	}

	_, ok := userStore.(*mspimpl.CertFileUserStore)
	if !ok {
		t.Fatal("Unexpected state store created")
	}
}

func newMockUserStore(t *testing.T) msp.UserStore {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mockmsp.NewMockIdentityConfig(mockCtrl)
	mockConfig.EXPECT().CredentialStorePath().Return("/tmp/fabsdkgo_test/store")

	factory := NewProviderFactory()
	userStore, err := factory.CreateUserStore(mockConfig)
	if err != nil {
		t.Fatalf("Unexpected error creating user store %s", err)
	}
	return userStore
}

func TestCreateUserStoreByConfig(t *testing.T) {
	userStore := newMockUserStore(t)
	_, ok := userStore.(*mspimpl.CertFileUserStore)
	if !ok {
		t.Fatal("Unexpected user store created")
	}
}

func TestCreateUserStoreEmptyConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mockmsp.NewMockIdentityConfig(mockCtrl)
	mockConfig.EXPECT().CredentialStorePath().Return("")

	factory := NewProviderFactory()
	_, err := factory.CreateUserStore(mockConfig)
	if err != nil {
		t.Fatal("Expected user store created")
	}
}

func TestCreateUserStoreFailConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mockmsp.NewMockIdentityConfig(mockCtrl)
	mockConfig.EXPECT().CredentialStorePath().Return("")

	factory := NewProviderFactory()
	_, err := factory.CreateUserStore(mockConfig)
	if err != nil {
		t.Fatal("Expected user store created")
	}
}

func TestCreateIdentityManager(t *testing.T) {

	coreFactory := defcore.NewProviderFactory()

	configPath := filepath.Join(metadata.GetProjectPath(), metadata.SDKConfigPath, "config_test.yaml")
	configBackend, err := config.FromFile(configPath)()
	if err != nil {
		t.Fatal(err)
	}

	cryptoCfg := cryptosuiteImpl.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatal(err)
	}

	endpointCfg, err := fab.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatal(err)
	}

	identityCfg, err := mspimpl.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatal(err)
	}

	cryptosuite, err := coreFactory.CreateCryptoSuiteProvider(cryptoCfg)
	if err != nil {
		t.Fatalf("Unexpected error creating cryptosuite provider %s", err)
	}

	factory := NewProviderFactory()
	userStore, err := factory.CreateUserStore(identityCfg)
	if err != nil {
		t.Fatalf("Unexpected error creating user store %s", err)
	}

	provider, err := factory.CreateIdentityManagerProvider(endpointCfg, cryptosuite, userStore)
	if err != nil {
		t.Fatalf("Unexpected error creating provider %s", err)
	}

	mgr, ok := provider.IdentityManager("Org1")
	if !ok {
		t.Fatalf("Unexpected error creating identity manager %s", err)
	}

	_, ok = mgr.(msp.IdentityManager)
	if !ok {
		t.Fatal("Unexpected signing manager created")
	}
}

func TestCreateUserStoreWithoutCredentialStorePath(t *testing.T) {

	configPath := filepath.Join(metadata.GetProjectPath(), "pkg", "core", "config", "testdata", "config_test_embedded_pems.yaml")
	configBackend, err := config.FromFile(configPath)()
	if err != nil {
		t.Fatal(err)
	}

	identityCfg, err := mspimpl.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, identityCfg.CredentialStorePath())

	factory := NewProviderFactory()
	userStore, err := factory.CreateUserStore(identityCfg)
	if err != nil {
		t.Fatalf("Unexpected error creating user store %s", err)
	}

	_, err = userStore.Load(msp.IdentityIdentifier{MSPID: "abc", ID: "ef"})
	assert.Equal(t, msp.ErrUserNotFound, err)

	assert.Equal(t, reflect.TypeOf(mspimpl.NewMemoryUserStore()), reflect.TypeOf(userStore))
}
