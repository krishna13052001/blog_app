package webservice

import (
	"blog_app/mycontext"
	"net/http"
)

func (s *WebService) Start(ctx mycontext.Context) {
	s.registerRoutes()
	s.server.Start(ctx)
}

func (s *WebService) registerRoutes() {
	s.server.AddNoAuthRoute("ping request", "GET", "/ping", s.ping)
	s.server.AddNoAuthRoute("create token", "POST", "/create-token", s.createToken)
	s.server.AddBasicRoute("create blog", "POST", "/blog", s.createBlog)
	s.server.AddNoAuthRoute("Get blog", "GET", "/blog", s.getBlog)
	s.server.AddNoAuthRoute("Get blog by id", "GET", "/blog/{id}", s.getBlogByID)
	s.server.AddBasicRoute("Delete blog", "DELETE", "/blog/{id}", s.deleteBlog)
	s.server.AddBasicRoute("Update blog", "PUT", "/blog", s.updateBlog)
	//TODO:: add routes for login and register
	s.server.AddNoAuthRoute("login", "POST", "/login", s.login)
	s.server.AddNoAuthRoute("register", "POST", "/register", s.register)
}

func (s *WebService) ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
