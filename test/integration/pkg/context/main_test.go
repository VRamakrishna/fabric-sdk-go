/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package context

import (
	"testing"

	"github.com/VRamakrishna/fabric-sdk-go/pkg/fabsdk"
	"github.com/VRamakrishna/fabric-sdk-go/test/integration"
	"github.com/VRamakrishna/fabric-sdk-go/test/integration/util/runner"
)

const (
	org1Name      = "Org1"
	org1AdminUser = "Admin"
)

var mainSDK *fabsdk.FabricSDK
var mainTestSetup *integration.BaseSetupImpl
var mainChaincodeID string

func TestMain(m *testing.M) {
	r := runner.NewWithExampleCC()
	r.Initialize()
	mainSDK = r.SDK()
	mainTestSetup = r.TestSetup()
	mainChaincodeID = r.ExampleChaincodeID()

	r.Run(m)
}
