package ffsecretsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// StringReader implements a SourceReader interface for AWS SecretManager.
// SecretManager can store secrets as strings or as binaries, this reader will read the secret as a string.
type StringReader struct {
	client Client
	params *secretsmanager.GetSecretValueInput
}

// NewStringReader creates a new StringReader.
func NewStringReader(client Client, params secretsmanager.GetSecretValueInput) *StringReader {
	return &StringReader{client: client, params: &params}
}

// Read the secret from AWS SecretManager.
func (s *StringReader) Read(ctx context.Context) ([]byte, error) {
	out, err := s.client.GetSecretValue(ctx, s.params)
	if err != nil {
		return nil, fmt.Errorf("secretsmanager get secret value: %w", err)
	}

	return []byte(*out.SecretString), nil
}
