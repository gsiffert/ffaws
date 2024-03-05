# ffaws

This package provides multiple feature-flags sources for AWS services to be used with the [featureflag](https://github.com/gsiffert/featureflag) package.

## Usage

```go
package main

import (
    "context"
    "time"
	
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
    "github.com/gsiffert/ffaws/ffsecretsmanager"
    "github.com/gsiffert/featureflag"
)

func main() {
    cfg := aws.NewConfig()
    secretManagerClient := secretsmanager.NewFromConfig(cfg)
    
    // Create a new SourceReader that will read the secret from AWS Secret Manager.
    secretRequest := secretsmanager.GetSecretValueInput{
        SecretId: aws.String("secretARN"),
    }
    source := ffsecretsmanager.NewStringReader(secretManagerClient, secretRequest)
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // This featureFlag will be updated every 30 seconds with the latest secret from AWS Secret Manager.
    featureFlag, err := featureflag.New(ctx, 30*time.Second, source)
    if err != nil {
        panic(err)
    }
}
```
