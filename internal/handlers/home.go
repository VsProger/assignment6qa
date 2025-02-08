package handlers

import (
	"errors"
	"forum/models"
	"forum/pkg/cookie"
	"net/http"
)

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.app.NotFound(w)
		return
	}

	if r.Method != http.MethodGet {
		h.app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	data, err := h.NewTemplateData(r)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data, err = h.service.SetUpPage(data, r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.app.NotFound(w)
			return
		} else {
			h.app.ServerError(w, err)
			return
		}
	}
	if data.Category_id == 0 {
		posts, err := h.service.GetAllPostPaginated(data.CurrentPage, data.Limit)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}

		data.Posts = posts
	} else {
		posts, err := h.service.GetAllPostByCategoryPaginated(data.CurrentPage, data.Limit, data.Category_id)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		data.Posts = posts
	}
	token := cookie.GetSessionCookie(r)
	if token != nil {
		reactions, err := h.service.GetReactionPosts(token.Value)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		data.Posts = h.service.IsLikedPost(data.Posts, reactions)
	}

	if len(*data.Posts) == 0 {
		data.Posts = nil
	}

	h.app.Render(w, http.StatusOK, "home.html", data)
	return
}

// SELECT count(*) FROM comments INNER JOIN posts ON comments.post_id=posts.id  GROUP by comments.post_id;
