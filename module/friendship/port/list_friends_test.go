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

type TestCase_ListFriends struct {
	name        string
	hasFinalErr bool
	bodyRequest ListFriendsReq

	listFriendsHandlerError error
	listFriendsData         []string

	hasValidateErr bool
}

func TestListFriends(t *testing.T) {
	t.Parallel()

	mockListFriendsHandler := new(mockHandler.MockListFriendsHandler)
	commandHandlerErr := errors.New("command handler error")

	req := ListFriendsReq{
		Email: "lisa@example.com",
	}
	tcs := []TestCase_ListFriends{
		{
			name:            "successful",
			bodyRequest:     req,
			listFriendsData: []string{"john@example.com", "kate@example.com"},
		},
		{
			name:           "fail because email empty",
			bodyRequest:    ListFriendsReq{},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name:           "fail because email invalid",
			bodyRequest:    ListFriendsReq{Email: "lisa-example.com"},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name:                    "fail because list friends handle has error",
			bodyRequest:             req,
			listFriendsHandlerError: commandHandlerErr,
			hasFinalErr:             true,
		},
	}

	for _, tc := range tcs {
		dataReq := tc.bodyRequest
		if !tc.hasValidateErr {
			mockListFriendsHandler.On("Handle", mock.Anything, tc.bodyRequest.Email).Once().Return(tc.listFriendsData, tc.listFriendsHandlerError)
		}

		server := NewServer(app.Application{
			Queries: app.Queries{
				ListFriends: mockListFriendsHandler,
			},
		})
		router := gin.Default()

		router.POST("/test", server.ListFriends)

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
			resBody := &ListFriendsRes{}
			err = json.Unmarshal(res.Body.Bytes(), resBody)
			assert.NoError(t, err)
			assert.Equal(t, &ListFriendsRes{
				Friends: tc.listFriendsData,
				Count:   len(tc.listFriendsData),
			}, resBody)
		}
		mock.AssertExpectationsForObjects(t, mockListFriendsHandler)
	}

}
