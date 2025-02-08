package handlers

import (
	"errors"
	"fmt"
	"forum/models"
	"forum/pkg/cookie"
	"forum/pkg/validator"
	"net/http"
	"strconv"
	"strings"
)

func (h *handler) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		h.app.NotFound(w)
		return
	}
	methodResolver(w, r, h.postCreateGet, h.postCreatePost)
}

func (h *handler) postCreateGet(w http.ResponseWriter, r *http.Request) {
	var err error
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data.Form = models.PostForm{}
	data.Categories, err = h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	h.app.Render(w, http.StatusOK, "create.html", data)
}

func (h *handler) postCreatePost(w http.ResponseWriter, r *http.Request) {
	form := models.PostForm{
		Title:            r.FormValue("title"),
		Content:          r.FormValue("content"),
		CategoriesString: r.Form["categories"],
	}
	categories, err := h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	trim(&form.Title, &form.Content)
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.NotSelected(form.CategoriesString), "categories", "At least one must be selected")
	form.CheckField(validator.IsError(form.ConverCategories(categories)), "categories", "This field is not correct")

	if !form.Valid() {
		data, err := h.NewTemplateData(r)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		data.Form = form

		data.Categories = categories
		h.app.Render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}
	cookies := cookie.GetSessionCookie(r)
	postID, err := h.service.CreatePost(form.Title, form.Content, cookies.Value, form.Categories)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

func (h *handler) postView(w http.ResponseWriter, r *http.Request) {
	id, _ := strings.CutPrefix(r.URL.Path, "/post/")
	if strings.Contains(id, "/") {
		h.app.ClientError(w, 404)
		return
	}
	ID, err := strconv.Atoi(id)
	if err != nil || ID < 1 || id[0] == '0' {
		h.app.ClientError(w, 400)
		return
	}

	post, err := h.service.GetPostByID(ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.app.ClientError(w, http.StatusNotFound)
		} else {
			h.app.ServerError(w, err)
		}
		return
	}

	data, err := h.NewTemplateData(r)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data.Post = post
	token := cookie.GetSessionCookie(r)
	if token != nil {
		exists, reaction, err := h.service.GetReactionPost(token.Value, ID)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		if exists {
			if reaction == true {
				data.Post.IsLiked = 1
			} else {
				data.Post.IsLiked = -1
			}
		}
		reactions, err := h.service.GetReactionComment(token.Value, ID)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
		data.Post = h.service.IsLikedComment(data.Post, reactions)
	}

	data.Form = models.CommentForm{}
	data.Categories, err = h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	h.app.Render(w, http.StatusOK, "post.html", data)
}

func (h *handler) PostByUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/posts" {
		h.app.NotFound(w)
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
		} else {
			h.app.ServerError(w, err)
		}
		return
	}
	c := cookie.GetSessionCookie(r)
	posts, err := h.service.GetAllPostByUserPaginated(c.Value, data.CurrentPage, data.Limit)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data.Categories, err = h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data.Posts = posts

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
}

func (h *handler) LikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/liked" {
		h.app.NotFound(w)
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
		} else {
			h.app.ServerError(w, err)
		}
		return
	}
	c := cookie.GetSessionCookie(r)
	posts, err := h.service.GetLikedPostsPaginated(c.Value, data.CurrentPage, data.Limit)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data.Categories, err = h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	data.Posts = posts

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
}
