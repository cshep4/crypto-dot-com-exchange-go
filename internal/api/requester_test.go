package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cdcerrors "github.com/cshep4/crypto-dot-com-exchange-go/errors"
	"github.com/cshep4/crypto-dot-com-exchange-go/internal/api"
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

func TestRequester_Post_Error(t *testing.T) {
	type args struct {
		ctx    context.Context
		body   api.Request
		method string
	}
	tests := []struct {
		name   string
		client http.Client
		args
		expectedStatusCode int
		expectedErr        error
	}{
		{
			name: "returns error if error creating request",
			args: args{
				ctx:    nil,
				body:   api.Request{},
				method: "some method",
			},
			expectedErr: errors.New("net/http: nil Context"),
		},
		{
			name: "returns error if error from server",
			args: args{
				ctx:    context.Background(),
				body:   api.Request{},
				method: "some method",
			},
			client: http.Client{
				Transport: roundTripper{
					err: errors.New("some error"),
				},
			},
			expectedErr: errors.New("some error"),
		},
		{
			name: "returns errors if nil body returned",
			args: args{
				ctx:    context.Background(),
				body:   api.Request{},
				method: "some method",
			},
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusOK,
					response:   nil,
				},
			},
			expectedErr: errors.New("unexpected end of JSON input"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			requester := api.Requester{
				Client: &tt.client,
			}

			var response api.BaseResponse
			statusCode, err := requester.Post(tt.ctx, tt.body, tt.method, &response)
			require.Error(t, err)

			assert.Empty(t, response)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestRequester_Post_Success(t *testing.T) {
	type args struct {
		body   api.Request
		method string
	}
	tests := []struct {
		name   string
		client http.Client
		args
		expectedStatusCode int
		expectedResponse   api.BaseResponse
	}{
		{
			name: "returns status code and unmarshals response",
			args: args{
				body:   api.Request{},
				method: "some method",
			},
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: api.BaseResponse{
						ID:     "1234",
						Method: "some method",
						Code:   "5678",
					},
				},
			},
			expectedStatusCode: http.StatusTeapot,
			expectedResponse: api.BaseResponse{
				ID:     "1234",
				Method: "some method",
				Code:   "5678",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			t.Cleanup(ctrl.Finish)

			var response api.BaseResponse
			statusCode, err := api.Requester{Client: &tt.client}.Post(ctx, tt.body, tt.method, &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResponse, response)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
		})
	}
}

func TestRequester_Get_Error(t *testing.T) {
	type args struct {
		ctx    context.Context
		body   api.Request
		method string
	}
	tests := []struct {
		name   string
		client http.Client
		args
		expectedStatusCode int
		expectedErr        error
	}{
		{
			name: "returns error if error creating request",
			args: args{
				ctx:    nil,
				body:   api.Request{},
				method: "some method",
			},
			expectedErr: errors.New("net/http: nil Context"),
		},
		{
			name: "returns error if error from server",
			args: args{
				ctx:    context.Background(),
				body:   api.Request{},
				method: "some method",
			},
			client: http.Client{
				Transport: roundTripper{
					err: errors.New("some error"),
				},
			},
			expectedErr: errors.New("some error"),
		},
		{
			name: "returns errors if nil body returned",
			args: args{
				ctx:    context.Background(),
				body:   api.Request{},
				method: "some method",
			},
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusOK,
					response:   nil,
				},
			},
			expectedErr: errors.New("unexpected end of JSON input"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			requester := api.Requester{
				Client: &tt.client,
			}

			var response api.BaseResponse
			statusCode, err := requester.Post(tt.ctx, tt.body, tt.method, &response)
			require.Error(t, err)

			assert.Empty(t, response)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
			assert.Contains(t, err.Error(), tt.expectedErr.Error())
		})
	}
}

func TestRequester_Get_Success(t *testing.T) {
	type args struct {
		body   api.Request
		method string
	}
	tests := []struct {
		name   string
		client http.Client
		args
		expectedStatusCode int
		expectedResponse   api.BaseResponse
	}{
		{
			name: "returns status code and unmarshals response",
			args: args{
				body:   api.Request{},
				method: "some method",
			},
			client: http.Client{
				Transport: roundTripper{
					statusCode: http.StatusTeapot,
					response: api.BaseResponse{
						ID:     "1234",
						Method: "some method",
						Code:   "5678",
					},
				},
			},
			expectedStatusCode: http.StatusTeapot,
			expectedResponse: api.BaseResponse{
				ID:     "1234",
				Method: "some method",
				Code:   "5678",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			t.Cleanup(ctrl.Finish)

			var response api.BaseResponse
			statusCode, err := api.Requester{Client: &tt.client}.Get(ctx, tt.body, tt.method, &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedResponse, response)
			assert.Equal(t, tt.expectedStatusCode, statusCode)
		})
	}
}

func TestRequester_CheckErrorResponse_Error(t *testing.T) {
	type args struct {
		statusCode   int
		responseCode json.Number
	}
	tests := []struct {
		name string
		args
		expectedHTTPStatusCode int
		expectedCode           int64
		expectedErr            error
		underlyingErr          error
	}{
		{
			name: "returns unexpected error when response code is invalid",
			args: args{
				statusCode:   http.StatusTeapot,
				responseCode: "invalid code",
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedErr:            errors.New("invalid response code: invalid code"),
		},
		{
			name: "returns unexpected error when response code is invalid",
			args: args{
				statusCode:   http.StatusTeapot,
				responseCode: "10002",
			},
			expectedHTTPStatusCode: http.StatusTeapot,
			expectedCode:           10002,
			expectedErr:            cdcerrors.ErrUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.Requester{}.CheckErrorResponse(tt.statusCode, tt.responseCode)
			require.Error(t, err)

			var responseError cdcerrors.ResponseError
			require.True(t, errors.As(err, &responseError))

			if tt.underlyingErr != nil {
				assert.True(t, errors.Is(err, tt.underlyingErr))
			}

			assert.Equal(t, tt.expectedHTTPStatusCode, responseError.HTTPStatusCode)
			assert.Equal(t, tt.expectedCode, responseError.Code)
			assert.Equal(t, tt.expectedErr, responseError.Err)
		})
	}
}

func TestRequester_CheckErrorResponse_Success(t *testing.T) {
	type args struct {
		statusCode   int
		responseCode json.Number
	}
	tests := []struct {
		name string
		args
	}{
		{
			name: "returns nil when status code is 1xx",
			args: args{
				statusCode: 199,
			},
		},
		{
			name: "returns nil when status code is 2xx",
			args: args{
				statusCode: 299,
			},
		},
		{
			name: "returns nil when status code is 3xx",
			args: args{
				statusCode: 399,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.Requester{}.CheckErrorResponse(tt.statusCode, tt.responseCode)
			require.NoError(t, err)
		})
	}
}
