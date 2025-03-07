/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/fab"
)

const (
	cscc            = "cscc"
	csccConfigBlock = "GetConfigBlock"
)

func createConfigBlockInvokeRequest(channelID string) fab.ChaincodeInvokeRequest {
	cir := fab.ChaincodeInvokeRequest{
		ChaincodeID: cscc,
		Fcn:         csccConfigBlock,
		Args:        [][]byte{[]byte(channelID)},
	}
	return cir
}
