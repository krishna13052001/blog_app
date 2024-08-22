package webservice

import (
	"blog_app/models"
	"blog_app/mycontext"
	"net/http"
)

func (s *WebService) createToken(w http.ResponseWriter, r *http.Request) {
	ctx := mycontext.UpgradeCtx(r.Context())

	var token models.Token

	err := s.GetContent(&token, r)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to parse request", http.StatusBadRequest, err)
		return
	}

	var tokenValue string
	tokenValue, err = mycontext.GenerateJWT(token.UserEmail, token.UserName)
	if err != nil {
		s.ReturnErrorResponse(ctx, w, "Failed to create token", http.StatusInternalServerError, err)
		return
	}
	tokenRsp := models.TokenResponse{
		Token:   tokenValue,
		Message: "Token created successfully",
	}
	s.ReturnOKResponse(w, tokenRsp)
}
