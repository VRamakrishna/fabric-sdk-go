/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package preferpeer

import (
	"testing"

	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/fab"
	clientmocks "github.com/VRamakrishna/fabric-sdk-go/pkg/fab/events/client/mocks"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/fab/events/client/peerresolver/minblockheight"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/fab/mocks"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/msp/test/mockmsp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	channelID = "testchannel"
	org1MSP   = "Org1MSP"
	org2MSP   = "Org2MSP"

	p1O1 = clientmocks.NewMockStatefulPeer("p1_O1", "peer1.org1.com:7051", clientmocks.WithBlockHeight(100), clientmocks.WithMSP(org1MSP))
	p2O1 = clientmocks.NewMockStatefulPeer("p2_O1", "peer2.org1.com:7051", clientmocks.WithBlockHeight(109), clientmocks.WithMSP(org1MSP))
	p1O2 = clientmocks.NewMockStatefulPeer("p1_O2", "peer1.org2.com:7051", clientmocks.WithBlockHeight(111), clientmocks.WithMSP(org2MSP))
	p2O2 = clientmocks.NewMockStatefulPeer("p2_O2", "peer2.org2.com:7051", clientmocks.WithBlockHeight(112), clientmocks.WithMSP(org2MSP))

	peers = []fab.Peer{p1O1, p2O1, p1O2, p2O2}
)

func TestResolve(t *testing.T) {
	dispatcher := &clientmocks.MockDispatcher{}
	ctx := mocks.NewMockContext(mockmsp.NewMockSigningIdentity("test", org1MSP))
	config := &mocks.MockConfig{}
	config.SetCustomChannelConfig(channelID, &fab.ChannelEndpointConfig{
		Policies: fab.ChannelPolicies{
			EventService: fab.EventServicePolicy{
				Balancer: fab.RoundRobin,
			},
		},
	})
	ctx.SetEndpointConfig(config)

	resolver := New(dispatcher, ctx, channelID, []string{p2O1.URL()}, minblockheight.WithBlockHeightLagThreshold(0))
	peer, err := resolver.Resolve(peers)
	require.NoError(t, err)
	assert.Equalf(t, p2O2.URL(), peer.URL(), "expected peer1 from org1 to be selected since threshold is set to 0 (highest block height)")

	resolver = New(dispatcher, ctx, channelID, []string{p2O1.URL()}, minblockheight.WithBlockHeightLagThreshold(5))
	peer, err = resolver.Resolve(peers)
	require.NoError(t, err)
	assert.Equalf(t, p2O1.URL(), peer.URL(), "expected peer2 from org1 to be selected since threshold is set to 5 and it's the preferred peer")

	resolver = New(dispatcher, ctx, channelID, []string{p1O1.URL()}, minblockheight.WithBlockHeightLagThreshold(5))
	peer, err = resolver.Resolve(peers)
	require.NoError(t, err)
	assert.Equalf(t, p2O1.URL(), peer.URL(), "expected peer2 from org1 to be selected since the preferred peer (peer1 in org1) is below the threshold but the other peer in org1 is above the threshold")

	resolver = New(dispatcher, ctx, channelID, []string{p1O1.URL(), p2O1.URL()}, minblockheight.WithBlockHeightLagThreshold(-1))

	chosenPeers := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		peer, err := resolver.Resolve(peers)
		require.NoError(t, err)
		assert.Equalf(t, org1MSP, peer.MSPID(), "expected a peer from org1 to be selected since threshold is set to -1 (disabled)")
		chosenPeers[peer.URL()] = struct{}{}
	}
	assert.Equalf(t, 2, len(chosenPeers), "expecting only 2 peers to be chosen")
}

func TestShouldDisconnect(t *testing.T) {
	dispatcher := &clientmocks.MockDispatcher{LastBlock: 100}
	ctx := mocks.NewMockContext(mockmsp.NewMockSigningIdentity("test", org1MSP))

	resolver := New(dispatcher, ctx, channelID, []string{p2O1.URL(), p2O1.URL()}, minblockheight.WithBlockHeightLagThreshold(5))
	disconnect := resolver.ShouldDisconnect(peers, p2O2)
	assert.Truef(t, disconnect, "expecting peer to be disconnected since at least one of the preferred peers is within the threshold")

	resolver = New(dispatcher, ctx, channelID, []string{p2O1.URL()}, minblockheight.WithBlockHeightLagThreshold(1))
	disconnect = resolver.ShouldDisconnect(peers, p2O2)
	assert.Falsef(t, disconnect, "expecting peer not to have disconnected since peer2 in org1 is not within the threshold")
}
