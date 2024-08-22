package webservice

import (
	"blog_app/db"
	"blog_app/domain"
	"blog_app/log"
	"blog_app/mycontext"
	"blog_app/server"
	"bytes"
	"encoding/json"
	"net/http"
)

type WebService struct {
	Domain domain.Service
	DB     db.Service
	server server.Server
}

type responseMetaData struct { //to-do: it must be renamed to a generic response struct
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
	TraceID string `json:"traceID,omitempty"`
}

func NewWebservices(domainService domain.Service, repo db.Service, serviceRoute, servicePort string) *WebService {
	webserver := server.NewServer(servicePort, serviceRoute)
	return &WebService{
		Domain: domainService,
		DB:     repo,
		server: webserver,
	}
}

func (s *WebService) GetContent(v interface{}, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func (s *WebService) ReturnErrorResponse(ctx mycontext.Context,
	w http.ResponseWriter,
	responseErrorMessage string,
	statusCode int,
	logError error,
	fields ...log.FieldsMap) {
	var buf = new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if logError != nil {
		if len(fields) > 0 {
			log.GenericError(ctx, logError, fields[0])
		} else {
			log.GenericError(ctx, logError, nil)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	_ = encoder.Encode(responseMetaData{Code: statusCode, Message: responseErrorMessage, TraceID: ctx.RequestID})
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}

func (s *WebService) ReturnOKResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var buf = new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(response)
	if err != nil {
		log.GenericError(mycontext.New(), err)
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.GenericError(mycontext.New(), err)
	}
}

func (s *WebService) ReturnResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var buf = new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(response)
	if err != nil {
		log.GenericError(mycontext.New(), err)
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.GenericError(mycontext.New(), err)
	}

}
