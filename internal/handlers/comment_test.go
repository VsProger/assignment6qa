package handlers

import (
	"fmt"
	mock "forum/internal/repo/mocks"
	"net/http"
	"net/url"
	"testing"
)

func TestCommentCreate(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	const (
		validPostID  = "1"
		validComment = "Naaah, this one is not thaaaat great"
	)

	tests := []struct {
		name     string
		comment  string
		postID   string
		wantCode int
	}{
		{
			name:     "Valid comment",
			comment:  validComment,
			postID:   validPostID,
			wantCode: http.StatusOK,
		},
		{
			name:     "Blank content",
			comment:  "",
			postID:   validPostID,
			wantCode: http.StatusOK,
		},
		{
			name:     "Invalid postID (negative)",
			comment:  validComment,
			postID:   "-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Invalid postID (non-digit)",
			comment:  validComment,
			postID:   "nah",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("comment", tt.comment)

			code, _, _ := ts.postForm(t, fmt.Sprintf("/post/%s", tt.postID), form)
			mock.Equal(t, code, tt.wantCode)
		})
	}
}
