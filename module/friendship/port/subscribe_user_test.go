package port

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/phantranhieunhan/s3-assignment/common"
	mockHandler "github.com/phantranhieunhan/s3-assignment/mock/friendship/handler"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app/command/payload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCase_SubscribeUser struct {
	name        string
	hasFinalErr bool
	bodyRequest SubscribeUserReq

	commandHandlerError error

	hasValidateErr bool
}

func TestSubscribeUser(t *testing.T) {
	t.Parallel()

	mockSubscribeUserHandler := new(mockHandler.MockSubscribeUserHandler)
	commandHandlerErr := errors.New("command handler error")
	tcs := []TestCase_SubscribeUser{
		{
			name: "successful",
			bodyRequest: SubscribeUserReq{
				Requestor: "lisa@example.com",
				Target:    "john@example.com",
			},
		},
		{
			name: "fail because request emails is not 2",
			bodyRequest: SubscribeUserReq{
				Requestor: "lisa@example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because request emails is invalid",
			bodyRequest: SubscribeUserReq{
				Requestor: "lisa-example.com",
				Target:    "john@example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because target emails is invalid",
			bodyRequest: SubscribeUserReq{
				Requestor: "lisa@example.com",
				Target:    "john-example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because command handle has error",
			bodyRequest: SubscribeUserReq{
				Requestor: "lisa@example.com",
				Target:    "john@example.com",
			},
			commandHandlerError: commandHandlerErr,
			hasFinalErr:         true,
		},
	}

	for _, tc := range tcs {
		dataReq := tc.bodyRequest
		if !tc.hasValidateErr {
			mockSubscribeUserHandler.On("Handle", mock.Anything, payload.SubscriberUserPayloads{{
				Requestor: tc.bodyRequest.Requestor,
				Target:    tc.bodyRequest.Target,
			}}).Once().Return(tc.commandHandlerError)
		}

		server := NewServer(app.Application{
			Commands: app.Commands{
				SubscribeUser: mockSubscribeUserHandler,
			},
		})
		router := gin.Default()

		router.POST("/test", server.SubscribeUser)

		jsonBody, err := json.Marshal(dataReq)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		assert.NoError(t, err)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		if tc.hasFinalErr {
			assert.Equal(t, http.StatusBadRequest, res.Code)
		} else {
			assert.Equal(t, http.StatusOK, res.Code)
			resBody := &common.SuccessRes{}
			err = json.Unmarshal(res.Body.Bytes(), resBody)
			assert.NoError(t, err)
			assert.Equal(t, common.SimpleSuccessResponse(nil), resBody)
		}
	}
}
