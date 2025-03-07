/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/VRamakrishna/fabric-sdk-go/pkg/client/channel"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/errors/status"
	"github.com/VRamakrishna/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
)

// A Transaction represents a specific invocation of a transaction function, and provides
// flexibility over how that transaction is invoked. Applications should
// obtain instances of this class from a Contract using the
// Contract.CreateTransaction method.
//
// Instances of this class are stateful. A new instance <strong>must</strong>
// be created for each transaction invocation.
type Transaction struct {
	contract       *Contract
	request        *channel.Request
	endorsingPeers []string
	collections    []string
	isInit         bool
	eventch        chan *fab.TxStatusEvent
}

// TransactionOption functional arguments can be supplied when creating a transaction object
type TransactionOption = func(*Transaction) error

func newTransaction(name string, contract *Contract, options ...TransactionOption) (*Transaction, error) {
	qname := name
	if len(contract.name) > 0 {
		qname = contract.name + ":" + name
	}
	txn := &Transaction{
		contract: contract,
		request:  &channel.Request{ChaincodeID: contract.chaincodeID, Fcn: qname},
	}

	for _, option := range options {
		err := option(txn)
		if err != nil {
			return nil, err
		}
	}

	return txn, nil
}

// WithTransient is an optional argument to the CreateTransaction method which
// sets the transient data that will be passed to the transaction function
// but will not be stored on the ledger. This can be used to pass
// private data to a transaction function.
func WithTransient(data map[string][]byte) TransactionOption {
	return func(txn *Transaction) error {
		txn.request.TransientMap = data
		return nil
	}
}

// WithEndorsingPeers is an optional argument to the CreateTransaction method which
// sets the peers that should be used for endorsement of transaction submitted to the ledger using Submit()
func WithEndorsingPeers(peers ...string) TransactionOption {
	return func(txn *Transaction) error {
		txn.endorsingPeers = peers
		return nil
	}
}

// WithCollections is an optional argument to the CreateTransaction method which sets the collections
func WithCollections(collections ...string) TransactionOption {
	return func(txn *Transaction) error {
		txn.collections = collections
		return nil
	}
}

// WithInit is an option that makes the transaction fulfill the --init-required condition before the chaincode
// can be invoked to process regular transactions
func WithInit() TransactionOption {
	return func(txn *Transaction) error {
		txn.isInit = true
		return nil
	}
}

// Evaluate a transaction function and return its results.
// The transaction function will be evaluated on the endorsing peers but
// the responses will not be sent to the ordering service and hence will
// not be committed to the ledger. This can be used for querying the world state.
func (txn *Transaction) Evaluate(args ...string) ([]byte, error) {
	bytes := make([][]byte, len(args))
	for i, v := range args {
		bytes[i] = []byte(v)
	}
	txn.request.Args = bytes

	var options []channel.RequestOption
	if txn.endorsingPeers != nil {
		options = append(options, channel.WithTargetEndpoints(txn.endorsingPeers...))
	}
	options = append(options, channel.WithTimeout(fab.Query, txn.contract.network.gateway.options.Timeout))

	if txn.collections != nil {
		txn.request.InvocationChain = append(txn.request.InvocationChain, &fab.ChaincodeCall{ID: txn.contract.chaincodeID, Collections: txn.collections})
	}

	response, err := txn.contract.client.Query(
		*txn.request,
		options...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to evaluate")
	}

	return response.Payload, nil
}

// Submit a transaction to the ledger. The transaction function represented by this object
// will be evaluated on the endorsing peers and then submitted to the ordering service
// for committing to the ledger.
func (txn *Transaction) Submit(args ...string) ([]byte, error) {
	bytes := make([][]byte, len(args))
	for i, v := range args {
		bytes[i] = []byte(v)
	}
	txn.request.Args = bytes
	txn.request.IsInit = txn.isInit

	var options []channel.RequestOption
	if txn.endorsingPeers != nil {
		options = append(options, channel.WithTargetEndpoints(txn.endorsingPeers...))
	}
	options = append(options, channel.WithTimeout(fab.Execute, txn.contract.network.gateway.options.Timeout))
	options = append(options, channel.WithRetry(retry.DefaultChannelOpts))

	if txn.collections != nil {
		txn.request.InvocationChain = append(txn.request.InvocationChain, &fab.ChaincodeCall{ID: txn.contract.chaincodeID, Collections: txn.collections})
	}

	response, err := txn.contract.client.InvokeHandler(
		newSubmitHandler(txn.eventch),
		*txn.request,
		options...,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to submit")
	}

	return response.Payload, nil
}

// RegisterCommitEvent registers for a commit event for this transaction.
//  Returns:
//  the channel that is used to receive the event. The channel is closed after the event is queued.
func (txn *Transaction) RegisterCommitEvent() <-chan *fab.TxStatusEvent {
	txn.eventch = make(chan *fab.TxStatusEvent, 1)
	return txn.eventch
}

func newSubmitHandler(eventch chan *fab.TxStatusEvent) invoke.Handler {
	return invoke.NewSelectAndEndorseHandler(
		invoke.NewEndorsementValidationHandler(
			invoke.NewSignatureValidationHandler(&commitTxHandler{eventch}),
		),
	)
}

type commitTxHandler struct {
	eventch chan *fab.TxStatusEvent
}

//Handle handles commit tx
func (c *commitTxHandler) Handle(requestContext *invoke.RequestContext, clientContext *invoke.ClientContext) {
	txnID := requestContext.Response.TransactionID

	//Register Tx event
	reg, statusNotifier, err := clientContext.EventService.RegisterTxStatusEvent(string(txnID)) // TODO: Change func to use TransactionID instead of string
	if err != nil {
		requestContext.Error = errors.Wrap(err, "error registering for TxStatus event")
		return
	}
	defer clientContext.EventService.Unregister(reg)
	_, err = createAndSendTransaction(clientContext.Transactor, requestContext.Response.Proposal, requestContext.Response.Responses)
	if err != nil {
		requestContext.Error = errors.Wrap(err, "CreateAndSendTransaction failed")
		return
	}

	select {
	case txStatus := <-statusNotifier:
		requestContext.Response.TxValidationCode = txStatus.TxValidationCode

		if txStatus.TxValidationCode != peer.TxValidationCode_VALID {
			requestContext.Error = status.New(status.EventServerStatus, int32(txStatus.TxValidationCode),
				"received invalid transaction", nil)
			return
		}

		if c.eventch != nil {
			c.eventch <- txStatus
			close(c.eventch)
		}
	case <-requestContext.Ctx.Done():
		requestContext.Error = status.New(status.ClientStatus, status.Timeout.ToInt32(),
			"Execute didn't receive block event", nil)
		return
	}
}

func createAndSendTransaction(sender fab.Sender, proposal *fab.TransactionProposal, resps []*fab.TransactionProposalResponse) (*fab.TransactionResponse, error) {

	txnRequest := fab.TransactionRequest{
		Proposal:          proposal,
		ProposalResponses: resps,
	}

	tx, err := sender.CreateTransaction(txnRequest)
	if err != nil {
		return nil, errors.WithMessage(err, "CreateTransaction failed")
	}

	transactionResponse, err := sender.SendTransaction(tx)
	if err != nil {
		return nil, errors.WithMessage(err, "SendTransaction failed")

	}

	return transactionResponse, nil
}
