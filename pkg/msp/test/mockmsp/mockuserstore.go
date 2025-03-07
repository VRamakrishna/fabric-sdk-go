/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mockmsp

import (
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/msp"
)

// MockUserStore ...
type MockUserStore struct {
}

// Store ...
func (m *MockUserStore) Store(*msp.UserData) error {
	return nil
}

// Load ...
func (m *MockUserStore) Load(identifier msp.IdentityIdentifier) (*msp.UserData, error) {
	return &msp.UserData{}, nil
}
