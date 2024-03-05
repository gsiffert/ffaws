package ffsecretsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// BinaryReader implements a SourceReader interface for AWS SecretManager.
// SecretManager can store secrets as strings or as binaries, this reader will read the secret as a binary.
type BinaryReader struct {
	client Client
	params *secretsmanager.GetSecretValueInput
}

// NewBinaryReader creates a new BinaryReader.
func NewBinaryReader(client Client, params *secretsmanager.GetSecretValueInput) *BinaryReader {
	return &BinaryReader{client: client, params: params}
}

// Read the secret from AWS SecretManager.
func (s *BinaryReader) Read(ctx context.Context) ([]byte, error) {
	out, err := s.client.GetSecretValue(ctx, s.params)
	if err != nil {
		return nil, fmt.Errorf("secretsmanager get secret value: %w", err)
	}

	return out.SecretBinary, nil
}
