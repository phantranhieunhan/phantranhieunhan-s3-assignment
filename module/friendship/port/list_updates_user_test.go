package port

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	mockHandler "github.com/phantranhieunhan/s3-assignment/mock/friendship/handler"
	"github.com/phantranhieunhan/s3-assignment/module/friendship/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCase_ListUpdatesUser struct {
	name        string
	hasFinalErr bool
	bodyRequest ListUpdatesUserReq

	ListUpdatesUserHandlerError error
	ListUpdatesUserData         []string

	hasValidateErr bool
}

func TestListUpdatesUser(t *testing.T) {
	t.Parallel()

	mockListUpdatesUserHandler := new(mockHandler.MockListUpdatesUserHandler)
	commandHandlerErr := errors.New("command handler error")

	req := ListUpdatesUserReq{
		Sender: "lisa@example.com",
		Text:   "Hello World! email1@example.com email2@example.com",
	}
	tcs := []TestCase_ListUpdatesUser{
		{
			name:                "successful",
			bodyRequest:         req,
			ListUpdatesUserData: []string{"john@example.com", "kate@example.com"},
		},
		{
			name:           "fail because email empty",
			bodyRequest:    ListUpdatesUserReq{},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name:           "fail because email invalid",
			bodyRequest:    ListUpdatesUserReq{Sender: "lisa-example.com"},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name:                        "fail because list friends handle has error",
			bodyRequest:                 req,
			ListUpdatesUserHandlerError: commandHandlerErr,
			hasFinalErr:                 true,
		},
	}

	for _, tc := range tcs {
		dataReq := tc.bodyRequest
		if !tc.hasValidateErr {
			mockListUpdatesUserHandler.On("Handle", mock.Anything, tc.bodyRequest.Sender, tc.bodyRequest.Text).Once().Return(tc.ListUpdatesUserData, tc.ListUpdatesUserHandlerError)
		}

		server := NewServer(app.Application{
			Queries: app.Queries{
				ListUpdatesUser: mockListUpdatesUserHandler,
			},
		})
		router := gin.Default()

		router.POST("/test", server.ListUpdatesUser)

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
			resBody := &ListUpdatesUserRes{}
			err = json.Unmarshal(res.Body.Bytes(), resBody)
			assert.NoError(t, err)
			assert.Equal(t, &ListUpdatesUserRes{
				Recipients: tc.ListUpdatesUserData,
			}, resBody)
		}
		mock.AssertExpectationsForObjects(t, mockListUpdatesUserHandler)
	}

}
