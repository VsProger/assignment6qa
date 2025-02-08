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

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		h.app.NotFound(w)
		return
	}
	methodResolver(w, r, h.loginGet, h.loginPost)
}

func (h *handler) loginGet(w http.ResponseWriter, r *http.Request) {
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data.Categories, err = h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data.Form = models.UserLoginForm{}
	h.app.Render(w, http.StatusOK, "login.html", data)
}

func (h *handler) loginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ping")
	form := models.UserLoginForm{
		Email:    strings.ToLower(r.FormValue("email")),
		Password: r.FormValue("password"),
	}
	fmt.Println(form)
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

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
		h.app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}
	session, err := h.service.Authenticate(form.Email, form.Password)

	fmt.Println(session, err)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			form.AddFieldError("email", "email doesn't exist")
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
			h.app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("password", models.ErrInvalidCredentials.Error())
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
			h.app.Render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			h.app.ServerError(w, err)
		}
		return
	}
	cookie.SetSessionCookie(w, session.Token, session.ExpTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handler) signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		h.app.NotFound(w)
		return
	}
	methodResolver(w, r, h.signupGet, h.signupPost)
}

func (h *handler) signupGet(w http.ResponseWriter, r *http.Request) {
	data, err := h.NewTemplateData(r)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data.Categories, err = h.service.GetAllCategory()
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	data.Form = models.UserSignupForm{}
	h.app.Render(w, http.StatusOK, "signup.html", data)
}

func (h *handler) signupPost(w http.ResponseWriter, r *http.Request) {
	form := models.UserSignupForm{
		Name:     r.FormValue("name"),
		Email:    strings.ToLower(r.FormValue("email")),
		Password: r.FormValue("password"),
	}
	fmt.Println(form)
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 12), "name", "This field must be 12 characters long maximum")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.IsEmail(form.Email), "email", "This field must be an email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

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
		h.app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	//
	user := form.FormToUser()
	err := h.service.CreateUser(user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data, err := h.NewTemplateData(r)
			if err != nil {
				h.app.ServerError(w, err)
				return
			}
			data.Form = form
			h.app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else if errors.Is(err, models.ErrDuplicateName) {
			form.AddFieldError("name", "Name is already in use")
			data, err := h.NewTemplateData(r)
			if err != nil {
				h.app.ServerError(w, err)
				return
			}
			data.Form = form
			h.app.Render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			h.app.ServerError(w, err)
		}
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *handler) logoutPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		h.app.NotFound(w)
		return
	}
	c := cookie.GetSessionCookie(r)
	if c != nil {
		h.service.DeleteSession(c.Value)
		cookie.ExpireSessionCookie(w)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
