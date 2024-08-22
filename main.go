package main

import (
	"blog_app/connectionMgr"
	"blog_app/constants"
	"blog_app/db"
	"blog_app/domain"
	"blog_app/log"
	"blog_app/mycontext"
	"blog_app/webservice"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := mycontext.New()
	constants.LoadEnv(ctx)
	log.GenericInfo(ctx, "Starting Server for the service "+constants.ServiceName)
	mongoClient, err := connectionMgr.NewMongoClient(constants.MongoHost, "blogapp", map[string]interface{}{})
	if err != nil {
		log.GenericError(ctx, errors.WithMessage(err, "can't connect to Mongodb"))
		return
	}
	mongoService := db.NewMongoService(mongoClient)
	domainService := domain.NewDomainService(mongoService)
	service := webservice.NewWebservices(domainService, mongoService, constants.ServiceRoute, constants.ServicePort)
	go func(ctx mycontext.Context) {
		service.Start(ctx)
	}(mycontext.CopyContext(ctx))
	log.GenericInfo(ctx, "Server Started and run on port "+constants.ServicePort)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	err = mongoClient.Disconnect(ctx)
	if err != nil {
		log.GenericError(mycontext.New(), errors.WithMessage(err, "can't disconnect to Mongodb"))
	}
	log.GenericInfo(ctx, constants.ServiceName+"Server Stopped")
}
