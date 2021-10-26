package cdcexchange_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cshep4/crypto-dot-com-exchange-go"
	"github.com/cshep4/crypto-dot-com-exchange-go/errors"
)

type roundTripper struct {
	statusCode int
	response   interface{}
	err        error
}

func (rt roundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	if rt.statusCode == 0 {
		rt.statusCode = 200
	}

	var body io.ReadCloser
	if rt.response != nil {
		b, err := json.Marshal(rt.response)
		if err != nil {
			return nil, err
		}

		body = ioutil.NopCloser(bytes.NewBufferString(string(b)))
	}

	return &http.Response{
		StatusCode: rt.statusCode,
		Status:     http.StatusText(rt.statusCode),
		Body:       body,
	}, rt.err
}

func TestNew_Error(t *testing.T) {
	type args struct {
		apiKey    string
		secretKey string
	}
	tests := []struct {
		name string
		args
		expectedErr error
	}{
		{
			name: "error when api key is empty",
			args: args{
				apiKey: "",
			},
			expectedErr: errors.InvalidParameterError{Parameter: "apiKey", Reason: "cannot be empty"},
		},
		{
			name: "error when secret key is empty",
			args: args{
				apiKey:    "api key",
				secretKey: "",
			},
			expectedErr: errors.InvalidParameterError{Parameter: "secretKey", Reason: "cannot be empty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := cdcexchange.New(tt.apiKey, tt.secretKey)
			require.Error(t, err)

			assert.Empty(t, client)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestNew_Success(t *testing.T) {
	type args struct {
		apiKey     string
		secretKey  string
		httpClient *http.Client
		opts       []cdcexchange.ClientOption
	}
	tests := []struct {
		name string
		args
		expectedBaseURL string
	}{
		{
			name: "successfully creates UAT client",
			args: args{
				apiKey:    "api key",
				secretKey: "secret key",
				opts:      []cdcexchange.ClientOption{cdcexchange.WithUATEnvironment()},
			},
			expectedBaseURL: cdcexchange.UATSandboxBaseURL,
		},
		{
			name: "successfully creates production client",
			args: args{
				apiKey:    "api key",
				secretKey: "secret key",
			},
			expectedBaseURL: cdcexchange.ProductionBaseURL,
		},
		{
			name: "successfully creates client with custom http client",
			args: args{
				apiKey:     "api key",
				secretKey:  "secret key",
				httpClient: &http.Client{Timeout: time.Minute},
				opts:       []cdcexchange.ClientOption{cdcexchange.WithHTTPClient(&http.Client{Timeout: time.Minute})},
			},
			expectedBaseURL: cdcexchange.ProductionBaseURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := cdcexchange.New(tt.apiKey, tt.secretKey, tt.opts...)
			require.NoError(t, err)
			require.NotEmpty(t, client)

			assert.Equal(t, tt.apiKey, client.APIKey())
			assert.Equal(t, tt.secretKey, client.SecretKey())
			assert.Equal(t, tt.expectedBaseURL, client.BaseURL())

			if tt.httpClient == nil {
				assert.Equal(t, http.DefaultClient, client.HTTPClient())
			} else {
				assert.Equal(t, tt.httpClient, client.HTTPClient())
			}
		})
	}
}

func TestClient_UpdateConfig_Error(t *testing.T) {
	type args struct {
		apiKey    string
		secretKey string
	}
	tests := []struct {
		name string
		args
		expectedErr error
	}{
		{
			name: "error when api key is empty",
			args: args{
				apiKey: "",
			},
			expectedErr: errors.InvalidParameterError{Parameter: "apiKey", Reason: "cannot be empty"},
		},
		{
			name: "error when secret key is empty",
			args: args{
				apiKey:    "api key",
				secretKey: "",
			},
			expectedErr: errors.InvalidParameterError{Parameter: "secretKey", Reason: "cannot be empty"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := cdcexchange.New("another api key", "another secret key")
			require.NoError(t, err)

			err = client.UpdateConfig(tt.apiKey, tt.secretKey)
			require.Error(t, err)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestClient_UpdateConfig_Success(t *testing.T) {
	type args struct {
		apiKey     string
		secretKey  string
		httpClient *http.Client
		opts       []cdcexchange.ClientOption
	}
	tests := []struct {
		name string
		args
		expectedBaseURL string
	}{
		{
			name: "successfully updates UAT client",
			args: args{
				apiKey:    "api key",
				secretKey: "secret key",
				opts:      []cdcexchange.ClientOption{cdcexchange.WithUATEnvironment()},
			},
			expectedBaseURL: cdcexchange.UATSandboxBaseURL,
		},
		{
			name: "successfully updates production client",
			args: args{
				apiKey:    "api key",
				secretKey: "secret key",
				opts:      []cdcexchange.ClientOption{cdcexchange.WithProductionEnvironment()},
			},
			expectedBaseURL: cdcexchange.ProductionBaseURL,
		},
		{
			name: "successfully updates production client",
			args: args{
				apiKey:    "api key",
				secretKey: "secret key",
			},
			expectedBaseURL: cdcexchange.ProductionBaseURL,
		},
		{
			name: "successfully updates http client",
			args: args{
				apiKey:     "api key",
				secretKey:  "secret key",
				httpClient: &http.Client{Timeout: time.Minute},
				opts:       []cdcexchange.ClientOption{cdcexchange.WithHTTPClient(&http.Client{Timeout: time.Minute})},
			},
			expectedBaseURL: cdcexchange.ProductionBaseURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := cdcexchange.New("another api key", "another secret key")
			require.NoError(t, err)

			err = client.UpdateConfig(tt.apiKey, tt.secretKey, tt.opts...)
			require.NoError(t, err)
			require.NotEmpty(t, client)

			assert.Equal(t, tt.apiKey, client.APIKey())
			assert.Equal(t, tt.secretKey, client.SecretKey())
			assert.Equal(t, tt.expectedBaseURL, client.BaseURL())

			if tt.httpClient != nil {
				assert.Equal(t, tt.httpClient, client.HTTPClient())
			}
		})
	}
}
