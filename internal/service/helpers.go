package service

import (
	"forum/models"
	"net/http"
	"strconv"
	"strings"
)

const (
	pageSize    = 5
	defaultPage = 1
)

var limitVariation = []int{5, 10, 15, 20, 50}

func AddCategory(arr []int) []int {
	for i, nb := range arr {
		arr[i] = nb + 1
	}
	return arr
}

func (s *service) SetUpPage(data *models.TemplateData, r *http.Request) (*models.TemplateData, error) {
	var err error
	currentPageStr := r.URL.Query().Get("page")
	data.Limit = validateLimit(r.URL.Query().Get("limit"))

	data.Category = strings.Title(r.URL.Query().Get("category"))
	data.Categories, err = s.GetAllCategory()
	if err != nil {
		return nil, err
	}
	if data.Category != "" {
		for key, value := range data.Categories {
			if data.Category == value {
				data.Category_id = key + 1
				break
			}
		}
		if data.Category_id == 0 {
			return nil, models.ErrNoRecord
		}
	}

	if r.URL.Path == "/user/posts" {
		data.NumberOfPage, err = s.repo.GetPageNumberMyPosts(data.Limit, int(data.User.ID))
	} else if r.URL.Path == "/user/liked" {
		data.NumberOfPage, err = s.repo.GetPageNumberLikedPosts(data.Limit, int(data.User.ID))
	} else {
		data.NumberOfPage, err = s.repo.GetPageNumber(data.Limit, data.Category_id)
	}
	if err != nil {
		return nil, err
	}

	data.CurrentPage, err = strconv.Atoi(currentPageStr)
	if err != nil || data.CurrentPage < 1 || data.CurrentPage > data.NumberOfPage {
		data.CurrentPage = defaultPage
	}
	data.URL = r.URL.Path
	data.LimitVariation = limitVariation
	return data, nil
}

func validateLimit(limitStr string) int {
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = pageSize
	}
	if limit <= 0 {
		limit = pageSize
	}
	return limit
}
