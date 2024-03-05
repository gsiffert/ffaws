// Package ffsecretsmanager implements a SourceReader interface for AWS SecretManager.
// This package is intended to be used with the github.com/gsiffert/featureflag package.
package ffsecretsmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// Client is a subset of secretsmanager.Client.
type Client interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}
