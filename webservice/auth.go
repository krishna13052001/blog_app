package webservice

import (
	"blog_app/auth"
	"blog_app/models"
	"blog_app/mycontext"
	"github.com/pkg/errors"
	"net/http"
)

func (s *WebService) login(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	var credentials models.Credentials
	err := s.GetContent(&credentials, r)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	user, err := s.Domain.ValidateUser(ctx, credentials)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Invalid username or password", http.StatusUnauthorized, err)
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to generate token", http.StatusInternalServerError, err)
		return
	}

	s.ReturnOKResponse(w, map[string]interface{}{"token": token, "status": "success"})
}

func (s *WebService) register(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	var user models.User
	err := s.GetContent(&user, r)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	if user.Approval == true {
		s.ReturnErrorResponse(ctx, w, "User already exists", http.StatusConflict, nil)
		return
	}

	if isExist, _ := s.Domain.UserExists(ctx, user.Email); isExist {
		s.ReturnErrorResponse(ctx, w, "User already exists", http.StatusConflict, nil)
		return
	}

	err = s.Domain.RegisterUser(ctx, user)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, errors.WithMessage(err, "Error while registering the user").Error(), http.StatusInternalServerError, nil)
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Username)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to generate token", http.StatusInternalServerError, err)
		return
	}
	s.ReturnOKResponse(w, map[string]interface{}{"token": token, "status": "success"})
}
