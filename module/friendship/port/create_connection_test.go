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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCase_ConnectFriendship struct {
	name        string
	hasFinalErr bool
	bodyRequest ConnectFriendshipReq

	commandHandlerError error

	hasValidateErr bool
}

func TestConnectFriendship(t *testing.T) {
	t.Parallel()

	mockConnectFriendshipHandler := new(mockHandler.MockConnectFriendshipHandler)
	commandHandlerErr := errors.New("command handler error")
	tcs := []TestCase_ConnectFriendship{
		{
			name: "successful",
			bodyRequest: ConnectFriendshipReq{
				Friends: []string{"lisa@example.com", "common@example.com"},
			},
		},
		{
			name: "fail because request emails is not 2",
			bodyRequest: ConnectFriendshipReq{
				Friends: []string{"lisa@example.com"},
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because request emails is the same",
			bodyRequest: ConnectFriendshipReq{
				Friends: []string{"lisa@example.com", "lisa@example.com"},
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name: "fail because command handle has error",
			bodyRequest: ConnectFriendshipReq{
				Friends: []string{"lisa@example.com", "common@example.com"},
			},
			commandHandlerError: commandHandlerErr,
			hasFinalErr:         true,
		},
	}

	for _, tc := range tcs {
		dataReq := tc.bodyRequest
		if !tc.hasValidateErr {
			mockConnectFriendshipHandler.On("Handle", mock.Anything, dataReq.Friends[0], dataReq.Friends[1]).Once().Return("", tc.commandHandlerError)
		}

		server := NewServer(app.Application{
			Commands: app.Commands{
				ConnectFriendship: mockConnectFriendshipHandler,
			},
		})
		router := gin.Default()

		router.POST("/test", server.ConnectFriendship)

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
