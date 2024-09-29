package server

import (
	"blog_app/log"
	"blog_app/mycontext"
	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"net/http/pprof"
)

func (s *server) Start(ctx mycontext.Context) {
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) // Allowing all origin as of now

	allowedHeaders := handlers.AllowedHeaders([]string{
		"Accept",
		"Content-Type",
		"contentType",
		"Content-Length",
		"Accept-Encoding",
		"Client-Security-Token",
		"X-CSRF-Token",
		"X-Auth-Token",
		"processData",
		"Authorization",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Connection",
		"Host",
		"Origin",
		"User-Agent",
		"Referer",
		"Cache-Control",
		"X-header",
		"X-Requested-With",
		"timezone",
		"email",
		"apitoken",
		"application",
		"username",
	})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"})

	allowCredential := handlers.AllowCredentials()

	serverHandler := handlers.CORS(
		allowedHeaders,
		allowedMethods,
		allowedOrigins,
		allowCredential)(
		context.ClearHandler(
			s.newRouter(s.subRoute),
		),
	)
	log.GenericInfo(ctx, "Starting Server",
		log.FieldsMap{
			"Port":     s.port,
			"SubRoute": s.subRoute,
		})

	err := http.ListenAndServe(":"+s.port, serverHandler)
	if err != nil {
		log.GenericError(ctx, errors.New("failed to start server"),
			log.FieldsMap{
				"Port":     s.port,
				"SubRoute": s.subRoute,
			})

		return
	}
}

func (s *server) newRouter(subRoute string) *mux.Router {
	muxRouter := mux.NewRouter().StrictSlash(true)
	muxRouter.HandleFunc(subRoute+"/debug/pprof", pprof.Index)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/cmdline", pprof.Cmdline)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/profile", pprof.Profile)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/symbol", pprof.Symbol)
	muxRouter.HandleFunc(subRoute+"/debug/pprof/trace", pprof.Trace)
	muxRouter.Handle(subRoute+"/debug/pprof/goroutine", pprof.Handler("goroutine"))
	muxRouter.Handle(subRoute+"/debug/pprof/heap", pprof.Handler("heap"))
	muxRouter.Handle(subRoute+"/debug/pprof/thread/create", pprof.Handler("threadcreate"))
	muxRouter.Handle(subRoute+"/debug/pprof/block", pprof.Handler("block"))
	for _, r := range s.routes {
		muxRouter.HandleFunc(subRoute+r.Pattern, r.HandlerFunc).Methods(r.Method)
	}

	return muxRouter
}

func (s *server) AddBasicRoute(apiDescription, methodType, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        apiDescription,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, recovery, enableCorsMiddleware, enableLogging, validateToken, createContext),
	}
	s.routes = append(s.routes, r)
}

func (s *server) AddNoAuthRoute(apiDescription, methodType, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        apiDescription,
		Method:      methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, recovery, enableCorsMiddleware, enableLogging, createContext),
	}
	s.routes = append(s.routes, r)
}
