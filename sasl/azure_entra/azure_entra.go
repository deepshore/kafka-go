package azure_entra

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/segmentio/kafka-go/sasl"
)

type Mechanism struct {
	tokenCredential     azcore.TokenCredential
	tokenRequestOptions policy.TokenRequestOptions
}

func (*Mechanism) Name() string {
	return "OAUTHBEARER"
}

func (m *Mechanism) Start(ctx context.Context) (sasl.StateMachine, []byte, error) {
	saslMeta := sasl.MetadataFromContext(ctx)

	if saslMeta == nil {
		return nil, nil, errors.New("missing sasl metadata")
	}

	entraToken, err := m.getEntraToken(ctx, saslMeta)

	if err != nil {
		return nil, nil, err
	}

	// See https://datatracker.ietf.org/doc/html/rfc7628
	saslResponse := fmt.Sprintf("n,,\x01auth=Bearer %s\x01\x01", entraToken.Token)

	return m, []byte(saslResponse), nil
}

func (m *Mechanism) getEntraToken(ctx context.Context, saslMeta *sasl.Metadata) (azcore.AccessToken, error) {

	entraToken, err := m.tokenCredential.GetToken(ctx, m.tokenRequestOptions)

	if err == nil {
		return entraToken, nil
	} else {
		err := fmt.Errorf("failed to request an Azure Entra Token: %w", err)
		return entraToken, err
	}

}

func (m *Mechanism) Next(ctx context.Context, challenge []byte) (done bool, response []byte, err error) {
	return true, nil, nil
}

func NewMechanism(tokenCredential azcore.TokenCredential, options *policy.TokenRequestOptions) *Mechanism {
	defaultTokenRequestOptions := policy.TokenRequestOptions{
		Scopes:    []string{"https://graph.microsoft.com/.default"},
		EnableCAE: false,
	}

	// If options is nil, use the default options
	if options == nil {
		options = &defaultTokenRequestOptions
	}

	return &Mechanism{
		tokenCredential:     tokenCredential,
		tokenRequestOptions: *options,
	}
}
