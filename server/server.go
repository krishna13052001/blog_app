package server

import (
	"blog_app/auth"
	"blog_app/log"
	"blog_app/mycontext"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

const (
	requestID = "requestId"
	userName  = "userName"
	userEmail = "userEmail"
	apiToken  = "apiToken"

	MyCtx = "myCtx"
)

type server struct {
	port     string
	subRoute string
	routes   []route
}

type Server interface {
	Start(ctx mycontext.Context)
	AddNoAuthRoute(apiDescription, methodType, mRoute string, handlerFunc http.HandlerFunc)
	AddBasicRoute(apiDescription, methodType, mRoute string, handlerFunc http.HandlerFunc)
}

func NewServer(port, subRoute string) Server {
	return &server{
		port:     port,
		subRoute: subRoute,
	}
}

func useMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			ctx := mycontext.UpgradeCtx(r.Context())
			rec := recover()
			if rec != nil {
				trace := string(debug.Stack())
				trace = strings.Replace(trace, "\n", "    ", -1)
				trace = strings.Replace(trace, "\t", "    ", -1)
				log.GenericError(ctx, fmt.Errorf("%v", rec),
					log.FieldsMap{
						"msg":        "recovering from panic",
						"stackTrace": trace,
					})
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func enableCorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}

}

func enableLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("enableLogging")
		ctx := mycontext.UpgradeCtx(r.Context())
		// Avoid logging ping API
		rawBody, _ := ioutil.ReadAll(r.Body)
		if len(rawBody) > 0 {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		}
		log.GenericInfo(ctx, "content type", log.FieldsMap{"value": http.DetectContentType(rawBody)})
		if strings.Contains(r.RequestURI, "ping") {
			return
		}

		var body log.FieldsMap
		body = log.FieldsMap{
			"method":  r.Method,
			"url":     r.RequestURI,
			"reqBody": string(rawBody),
		}
		log.GenericInfo(ctx, "trace", body)
		next.ServeHTTP(w, r)
	}
}

func createContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header
		ctx := r.Context()
		reqID := header.Get(requestID)
		if reqID == "" {
			reqID = strings.ReplaceAll(uuid.NewString(), "-", "")
		}
		name, email, token := header.Get(userName), header.Get(userEmail), header.Get(apiToken)

		myCtx := mycontext.MyContext{
			RequestID: reqID,
			UserName:  name,
			UserEmail: email,
			ApiToken:  token,
		}

		ctx = mycontext.WithCtx(ctx, myCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func validateToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		myCtx, exists := mycontext.GetMyCtx(ctx)

		if !exists || myCtx.ApiToken == "" {
			http.Error(w, "Authorization token is missing", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJWT(myCtx.ApiToken)
		if err != nil || !claims.Valid {
			http.Error(w, "Invalid or expired authorization token", http.StatusUnauthorized)
			return
		}

		if claims.Email != myCtx.UserEmail {
			http.Error(w, "Token email does not match provided email", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
