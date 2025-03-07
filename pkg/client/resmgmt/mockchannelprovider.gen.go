// Code generated by counterfeiter. DO NOT EDIT.
package resmgmt

import (
	"sync"

	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/fab"
)

type MockChannelProvider struct {
	ChannelServiceStub        func(ctx fab.ClientContext, channelID string) (fab.ChannelService, error)
	channelServiceMutex       sync.RWMutex
	channelServiceArgsForCall []struct {
		ctx       fab.ClientContext
		channelID string
	}
	channelServiceReturns struct {
		result1 fab.ChannelService
		result2 error
	}
	channelServiceReturnsOnCall map[int]struct {
		result1 fab.ChannelService
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MockChannelProvider) ChannelService(ctx fab.ClientContext, channelID string) (fab.ChannelService, error) {
	fake.channelServiceMutex.Lock()
	ret, specificReturn := fake.channelServiceReturnsOnCall[len(fake.channelServiceArgsForCall)]
	fake.channelServiceArgsForCall = append(fake.channelServiceArgsForCall, struct {
		ctx       fab.ClientContext
		channelID string
	}{ctx, channelID})
	fake.recordInvocation("ChannelService", []interface{}{ctx, channelID})
	fake.channelServiceMutex.Unlock()
	if fake.ChannelServiceStub != nil {
		return fake.ChannelServiceStub(ctx, channelID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.channelServiceReturns.result1, fake.channelServiceReturns.result2
}

func (fake *MockChannelProvider) ChannelServiceCallCount() int {
	fake.channelServiceMutex.RLock()
	defer fake.channelServiceMutex.RUnlock()
	return len(fake.channelServiceArgsForCall)
}

func (fake *MockChannelProvider) ChannelServiceArgsForCall(i int) (fab.ClientContext, string) {
	fake.channelServiceMutex.RLock()
	defer fake.channelServiceMutex.RUnlock()
	return fake.channelServiceArgsForCall[i].ctx, fake.channelServiceArgsForCall[i].channelID
}

func (fake *MockChannelProvider) ChannelServiceReturns(result1 fab.ChannelService, result2 error) {
	fake.ChannelServiceStub = nil
	fake.channelServiceReturns = struct {
		result1 fab.ChannelService
		result2 error
	}{result1, result2}
}

func (fake *MockChannelProvider) ChannelServiceReturnsOnCall(i int, result1 fab.ChannelService, result2 error) {
	fake.ChannelServiceStub = nil
	if fake.channelServiceReturnsOnCall == nil {
		fake.channelServiceReturnsOnCall = make(map[int]struct {
			result1 fab.ChannelService
			result2 error
		})
	}
	fake.channelServiceReturnsOnCall[i] = struct {
		result1 fab.ChannelService
		result2 error
	}{result1, result2}
}

func (fake *MockChannelProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.channelServiceMutex.RLock()
	defer fake.channelServiceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MockChannelProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ fab.ChannelProvider = new(MockChannelProvider)
