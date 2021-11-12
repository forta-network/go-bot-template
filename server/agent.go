package server

import (
	"context"
	"time"

	"forta-protocol/go-agent/protocol"
)

type Agent struct {
	protocol.UnimplementedAgentServer
}

func (a *Agent) Initialize(ctx context.Context, request *protocol.InitializeRequest) (*protocol.InitializeResponse, error) {
	panic("implement me")
}

func (a *Agent) EvaluateTx(ctx context.Context, request *protocol.EvaluateTxRequest) (*protocol.EvaluateTxResponse, error) {
	return &protocol.EvaluateTxResponse{
		Status: protocol.ResponseStatus_SUCCESS,
		Findings: []*protocol.Finding{
			{
				Protocol:    "ethereum",
				Severity:    protocol.Finding_MEDIUM,
				Type:        protocol.Finding_INFORMATION,
				AlertId:     "test-alert",
				Name:        "Test Alert",
				Metadata:    map[string]string{},
				Description: "This is a test alert",
			},
		},
		Metadata:  map[string]string{},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (a *Agent) EvaluateBlock(ctx context.Context, request *protocol.EvaluateBlockRequest) (*protocol.EvaluateBlockResponse, error) {
	return &protocol.EvaluateBlockResponse{
		Status:    protocol.ResponseStatus_SUCCESS,
		Findings:  nil,
		Metadata:  map[string]string{},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}
