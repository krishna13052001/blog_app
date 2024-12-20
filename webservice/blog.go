package webservice

import (
	"blog_app/models"
	"blog_app/mycontext"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func (s *WebService) createBlog(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	var blog models.Blog

	err := s.GetContent(&blog, r)

	blog.BatchArray = strings.Split(blog.Batch, ",")
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	err = s.Domain.CreateBlog(ctx, blog)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to create blog", http.StatusInternalServerError, err)
		return
	}

	s.ReturnResponse(w, http.StatusCreated, "Blog created successfully")
}

func (s *WebService) getBlog(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())
	start := mux.Vars(r)["start"]
	blogs, err := s.Domain.GetBlog(ctx, start)

	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to get blog", http.StatusInternalServerError, err)
		return
	}

	value, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		value = 0
	}
	value += int64(len(blogs))
	s.ReturnOKResponse(w, map[string]interface{}{"blogs": blogs, "status": "success", "end": value})
}

func (s *WebService) searchLocationAndBatches(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())
	start := r.URL.Query().Get("start")
	location := r.URL.Query().Get("location")
	batch := r.URL.Query().Get("batch")
	jobtype := r.URL.Query().Get("jobtype")
	typeVal := 0
	queryValue := ""
	if location != "" {
		typeVal = 1
		queryValue = location
	} else if batch != "" {
		typeVal = 2
		queryValue = batch
	} else if jobtype != "" {
		typeVal = 3
		queryValue = jobtype
	}
	blogs, err := s.Domain.GetBlogsByFilter(ctx, queryValue, typeVal, start)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to get blog", http.StatusInternalServerError, err)
		return
	}
	value, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		value = 0
	}
	value += int64(len(blogs))
	s.ReturnOKResponse(w, map[string]interface{}{"blogs": blogs, "status": "success", "value": value})
}

func (s *WebService) getBlogByID(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	id := mux.Vars(r)["id"]

	if id == "" {
		s.ReturnErrorResponse(ctx, w, "Invalid blog id", http.StatusBadRequest, nil)
		return
	}

	blog, err := s.Domain.GetBlogById(ctx, id)

	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to get blog", http.StatusInternalServerError, err)
		return
	}

	s.ReturnOKResponse(w, map[string]interface{}{"blog": blog, "status": "success"})
}

func (s *WebService) deleteBlog(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	id := mux.Vars(r)["id"]

	if id == "" {
		s.ReturnErrorResponse(ctx, w, "Invalid blog id", http.StatusBadRequest, nil)
		return
	}

	err := s.Domain.DeleteBlog(ctx, id)

	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to delete blog", http.StatusInternalServerError, err)
		return
	}
	s.ReturnOKResponse(w, "Blog deleted successfully")
}

func (s *WebService) updateBlog(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	var blog models.Blog

	err := s.GetContent(&blog, r)

	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	err = s.Domain.UpdateBlog(ctx, blog)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to update blog", http.StatusInternalServerError, err)
		return
	}

	s.ReturnOKResponse(w, "Blog updated successfully")
}
