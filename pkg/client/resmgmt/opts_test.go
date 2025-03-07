/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package resmgmt

import (
	"testing"

	"time"

	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/fab"
	fcmocks "github.com/VRamakrishna/fabric-sdk-go/pkg/fab/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWithTargetURLsInvalid(t *testing.T) {
	ctx := setupTestContext("test", "Org1MSP")
	opt := WithTargetEndpoints("invalid")

	mockConfig := &fcmocks.MockConfig{}

	oConfig := &fab.PeerConfig{
		URL: "127.0.0.1:7050",
	}

	mockConfig.SetCustomPeerCfg(oConfig)
	ctx.SetEndpointConfig(mockConfig)

	opts := requestOptions{}
	err := opt(ctx, &opts)
	assert.NotNil(t, err, "Should have failed for invalid target peer")
}

func TestWithTargetURLsValid(t *testing.T) {
	ctx := setupTestContext("test", "Org1MSP")
	opt := WithTargetEndpoints("127.0.0.1:7050")

	mockConfig := &fcmocks.MockConfig{}

	pConfig1 := fab.PeerConfig{
		URL: "127.0.0.1:7050",
	}

	npConfig1 := fab.NetworkPeer{
		PeerConfig: pConfig1,
		MSPID:      "MYMSP",
	}

	pConfig2 := fab.PeerConfig{
		URL: "127.0.0.1:7051",
	}

	npConfig2 := fab.NetworkPeer{
		PeerConfig: pConfig2,
		MSPID:      "OTHERMSP",
	}

	mockConfig.SetCustomPeerCfg(&pConfig1)
	mockConfig.SetCustomNetworkPeerCfg([]fab.NetworkPeer{npConfig2, npConfig1})
	ctx.SetEndpointConfig(mockConfig)

	opts := requestOptions{}
	err := opt(ctx, &opts)
	assert.Nil(t, err, "Should have failed for invalid target peer")

	assert.Equal(t, 1, len(opts.Targets), "should have one peer")
	assert.Equal(t, pConfig1.URL, opts.Targets[0].URL(), "", "Wrong URL")
	assert.Equal(t, npConfig1.MSPID, opts.Targets[0].MSPID(), "", "Wrong MSP")
}

func TestTimeoutOptions(t *testing.T) {

	opts := requestOptions{}

	options := []RequestOption{WithTimeout(fab.PeerResponse, 20*time.Second),
		WithTimeout(fab.ResMgmt, 25*time.Second), WithTimeout(fab.OrdererResponse, 30*time.Second),
		WithTimeout(fab.PeerConnection, 35*time.Second), WithTimeout(fab.Execute, 40*time.Second),
		WithTimeout(fab.Query, 45*time.Second)}

	for _, option := range options {
		option(nil, &opts)
	}

	assert.True(t, opts.Timeouts[fab.PeerResponse] == 20*time.Second, "timeout value by type didn't match with one supplied")
	assert.True(t, opts.Timeouts[fab.ResMgmt] == 25*time.Second, "timeout value by type didn't match with one supplied")
	assert.True(t, opts.Timeouts[fab.OrdererResponse] == 30*time.Second, "timeout value by type didn't match with one supplied")
	assert.True(t, opts.Timeouts[fab.PeerConnection] == 35*time.Second, "timeout value by type didn't match with one supplied")
	assert.True(t, opts.Timeouts[fab.Execute] == 40*time.Second, "timeout value by type didn't match with one supplied")
	assert.True(t, opts.Timeouts[fab.Query] == 45*time.Second, "timeout value by type didn't match with one supplied")

}
