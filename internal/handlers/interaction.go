package handlers

import (
	"errors"
	"fmt"
	"forum/models"
	"forum/pkg/cookie"
	"forum/pkg/validator"
	"net/http"
	"strings"
)

func (h *handler) postReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction" {
		h.app.NotFound(w)
		return
	}
	if r.Method != http.MethodPost {
		h.app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.app.ServerError(w, err)
		return
	}

	token := cookie.GetSessionCookie(r)
	postID, err := GetIntForm(r, "postID")
	if err != nil {
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}
	form := models.ReactionForm{
		ID:    postID,
		Token: token.Value,
	}
	reaction := r.FormValue("reaction")

	switch reaction {
	case "true":
		form.Reaction = true
	case "false":
		form.Reaction = false
	default:
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}
	err = h.service.PostReaction(form)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.app.ClientError(w, http.StatusBadRequest)
			return
		}
		h.app.ServerError(w, err)
		return
	}
	url := strings.TrimPrefix(r.Header.Get("Referer"), r.Header.Get("Origin"))

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *handler) commentPost(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.URL.Path != "/comment/post" {
		h.app.NotFound(w)
		return
	}
	if r.Method != http.MethodPost {
		h.app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	token := cookie.GetSessionCookie(r)
	postID, err := GetIntForm(r, "postID")
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	form := models.CommentForm{
		Content: r.FormValue("comment"),
		PostID:  postID,
		Token:   token.Value,
	}
	trim(&form.Content)
	form.CheckField(validator.NotBlank(form.Content), "comment", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Content, 2), "comment", "This field must be at least 2 characters long")
	form.CheckField(validator.MaxChars(form.Content, 100), "comment", "This field must be maximum 100 characters")

	if !form.Valid() {
		data, err := h.NewTemplateData(r)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		data.Form = form
		data.Categories, err = h.service.GetAllCategory()
		if err != nil {
			h.app.ServerError(w, err)
			return
		}

		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		post, err := h.service.GetPostByID(form.PostID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				h.app.ClientError(w, http.StatusNotFound)
				return
			} else {
				h.app.ServerError(w, err)
				return
			}
		}
		data.Post = post
		h.app.Render(w, http.StatusUnprocessableEntity, "post.html", data)
		return
	}

	err = h.service.CommentPost(form)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", form.PostID), http.StatusSeeOther)
}

func (h *handler) commentReaction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/reaction" {
		h.app.NotFound(w)
		return
	}
	if r.Method != http.MethodPost {
		h.app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.app.ServerError(w, err)
		return
	}

	postID, err := GetIntForm(r, "postID")
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	commentID, err := GetIntForm(r, "commentID")
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	token := cookie.GetSessionCookie(r)
	form := models.ReactionForm{
		ID:    commentID,
		Token: token.Value,
	}
	reaction := r.FormValue("reaction")

	switch reaction {
	case "true":
		form.Reaction = true
	case "false":
		form.Reaction = false
	default:
		h.app.ClientError(w, http.StatusBadRequest)
		return
	}
	err = h.service.CommentReaction(form)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.app.ClientError(w, http.StatusBadRequest)
			return
		}
		h.app.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
