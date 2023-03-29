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

type TestCase_BlockUpdatesUser struct {
	name        string
	hasFinalErr bool
	bodyRequest BlockUpdatesUserReq

	commandHandlerError error

	hasValidateErr bool
}

func TestBlockUpdatesUser(t *testing.T) {
	t.Parallel()

	mockBlockUpdatesUserHandler := new(mockHandler.MockBlockUpdatesUserHandler)
	commandHandlerErr := errors.New("command handler error")
	tcs := []TestCase_BlockUpdatesUser{
		{
			name: "successful",
			bodyRequest: BlockUpdatesUserReq{
				Requestor: "lisa@example.com",
				Target:    "john@example.com",
			},
		},
		{
			name: "fail because target email is not provided",
			bodyRequest: BlockUpdatesUserReq{
				Requestor: "lisa@example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because target email invalid",
			bodyRequest: BlockUpdatesUserReq{
				Target:    "lisa-example.com",
				Requestor: "john@example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because requestor email is not provided",
			bodyRequest: BlockUpdatesUserReq{
				Target: "lisa@example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because requestor email invalid",
			bodyRequest: BlockUpdatesUserReq{
				Requestor: "lisa-example.com",
				Target:    "john@example.com",
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because command handle has error",
			bodyRequest: BlockUpdatesUserReq{
				Target:    "lisa@example.com",
				Requestor: "john@example.com",
			},
			commandHandlerError: commandHandlerErr,
			hasFinalErr:         true,
		},
	}

	for _, tc := range tcs {
		dataReq := tc.bodyRequest
		if !tc.hasValidateErr {
			mockBlockUpdatesUserHandler.On("Handle", mock.Anything, payload.BlockUpdatesUserPayload{
				Requestor: tc.bodyRequest.Requestor,
				Target:    tc.bodyRequest.Target,
			}).Once().Return(tc.commandHandlerError)
		}

		server := NewServer(app.Application{
			Commands: app.Commands{
				BlockUpdatesUser: mockBlockUpdatesUserHandler,
			},
		})
		router := gin.Default()

		router.POST("/test", server.BlockUpdatesUser)

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
