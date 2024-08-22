package webservice

import (
	"blog_app/models"
	"blog_app/mycontext"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *WebService) createBlog(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	var blog models.Blog

	err := s.GetContent(&blog, r)

	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	err = s.Domain.CreateBlog(ctx, blog)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to create blog", http.StatusInternalServerError, err)
		return
	}

	s.ReturnOKResponse(w, "Blog created successfully")
}

func (s *WebService) getBlog(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	blogs, err := s.Domain.GetBlog(ctx)

	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to get blog", http.StatusInternalServerError, err)
		return
	}

	s.ReturnOKResponse(w, map[string]interface{}{"blogs": blogs, "status": "success"})
}

func (s *WebService) getBlogByID(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	id := mux.Vars(r)["id"]
	fmt.Println("id", id)

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
	fmt.Println("id", id)

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
