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

type TestCase_ListCommonFriends struct {
	name        string
	hasFinalErr bool
	bodyRequest ListCommonFriendsReq

	listCommonFriendsHandlerError error
	listCommonFriendsData         []string

	hasValidateErr bool
}

func TestListCommonFriends(t *testing.T) {
	t.Parallel()

	mockListCommonFriendsHandler := new(mockHandler.MockListCommonFriendsHandler)
	commandHandlerErr := errors.New("command handler error")

	req := ListCommonFriendsReq{
		Friends: []string{"lisa@example.com", "common@example.com"},
	}
	tcs := []TestCase_ListCommonFriends{
		{
			name:                  "successful",
			bodyRequest:           req,
			listCommonFriendsData: []string{"john@example.com", "kate@example.com"},
		},
		{
			name: "fail because request emails is not 2",
			bodyRequest: ListCommonFriendsReq{
				Friends: []string{"lisa@example.com"},
			},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name:           "fail because email invalid",
			bodyRequest:    ListCommonFriendsReq{Friends: []string{"lisa-example.com"}},
			hasValidateErr: true,
			hasFinalErr:    true,
		},
		{
			name:                          "fail because list friends handle has error",
			bodyRequest:                   req,
			listCommonFriendsHandlerError: commandHandlerErr,
			hasFinalErr:                   true,
		},
	}

	for _, tc := range tcs {
		dataReq := tc.bodyRequest
		if !tc.hasValidateErr {
			mockListCommonFriendsHandler.On("Handle", mock.Anything, tc.bodyRequest.Friends).Once().Return(tc.listCommonFriendsData, tc.listCommonFriendsHandlerError)
		}

		server := NewServer(app.Application{
			Queries: app.Queries{
				ListCommonFriends: mockListCommonFriendsHandler,
			},
		})
		router := gin.Default()

		router.POST("/test", server.ListCommonFriends)

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
			resBody := &ListCommonFriendsResp{}
			err = json.Unmarshal(res.Body.Bytes(), resBody)
			assert.NoError(t, err)
			assert.Equal(t, &ListCommonFriendsResp{
				Friends: tc.listCommonFriendsData,
				Count:   len(tc.listCommonFriendsData),
			}, resBody)
		}
		mock.AssertExpectationsForObjects(t, mockListCommonFriendsHandler)
	}

}
