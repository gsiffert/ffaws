// Package ffappconfig implements a SourceReader interface for AWS AppConfig.
// This package is intended to be used with the github.com/gsiffert/featureflag package.
package ffappconfig

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
)

// Client is a subset of appconfigdata.Client.
type Client interface {
	StartConfigurationSession(ctx context.Context, params *appconfigdata.StartConfigurationSessionInput, optFns ...func(*appconfigdata.Options)) (*appconfigdata.StartConfigurationSessionOutput, error)
	GetLatestConfiguration(ctx context.Context, params *appconfigdata.GetLatestConfigurationInput, optFns ...func(*appconfigdata.Options)) (*appconfigdata.GetLatestConfigurationOutput, error)
}

// AppConfigReader implements a SourceReader interface for AWS AppConfig.
// AppConfig restricts how regularly you can poll for configuration changes,
// this reader will cache the results from AppConfig and only poll when it's allowed to.
type AppConfigReader struct {
	client     Client
	params     *appconfigdata.StartConfigurationSessionInput
	once       sync.Once
	token      *string
	nextRun    time.Time
	lastResult lastResult
}

type lastResult struct {
	configuration []byte
	err           error
}

// New creates a new AppConfigReader.
func New(client Client, params appconfigdata.StartConfigurationSessionInput) *AppConfigReader {
	return &AppConfigReader{client: client, params: &params}
}

func (s *AppConfigReader) startSession(ctx context.Context) error {
	var err error

	s.once.Do(func() {
		result, e := s.client.StartConfigurationSession(ctx, s.params)
		if e != nil {
			err = fmt.Errorf("appconfig start configuration session: %w", e)
			return
		}

		s.token = result.InitialConfigurationToken
	})
	return err
}

// Read the latest configuration from AWS AppConfig.
func (s *AppConfigReader) Read(ctx context.Context) ([]byte, error) {
	if err := s.startSession(ctx); err != nil {
		return nil, err
	}

	// If the next run is in the future, it means AWS expects us to wait before polling again.
	if time.Now().UTC().Before(s.nextRun) {
		return s.lastResult.configuration, s.lastResult.err
	}

	params := &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: s.token,
	}
	out, err := s.client.GetLatestConfiguration(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("appconfig get latest configuration: %w", err)
	}

	// Update the token and the next run time.
	s.token = out.NextPollConfigurationToken
	s.nextRun = time.Now().UTC().Add(time.Duration(out.NextPollIntervalInSeconds) * time.Second)
	// We also cache the result in case the next call to Get is before the next run time.
	s.lastResult.configuration = out.Configuration
	s.lastResult.err = nil
	return out.Configuration, nil
}
