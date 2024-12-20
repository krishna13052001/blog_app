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

/*
	// Define the ITestFramework interface
	type ITestFramework interface {
		Setup(t *testing.T) *TestFramework
		Run(name string, test FrameworkTest) bool
		CreateGateway(gateway Gateway) (map[string]any, error)
		DeleteGateway(name string) (map[string]any, error)
		API() *API
		Ctx() context.Context
		T() *testing.T
	}

	// Refactor TestFramework to implement ITestFramework
	type TestFramework struct {
		api *API
		ctx context.Context
		t *testing.T
		controller  testcontainers.Container
		gateway     testcontainers.Container
		gwSimClient gatewaySimulator.SimulatorServiceClient
	}

	func (tf *TestFramework) API() *API {
		return tf.api
	}

	func (tf *TestFramework) Ctx() context.Context {
		return tf.ctx
	}

	func (tf *TestFramework) T() *testing.T {
		return tf.t
	}

	func (tf *TestFramework) Setup(t *testing.T) *TestFramework {
		tf.api = api
		tf.ctx = logging.InitTestCtx(t)
		tf.t = t
		return tf
	}

	func TestFrameworkFactory(env string) ITestFramework {
		switch env {
		case "live":
			return &LiveTestFramework{}
		case "aviatrix_8.0":
			return &Aviatrix80TestFramework{}
		default:
			return &TestFramework{}
		}
	}

	type LiveTestFramework struct {
		*TestFramework
	}

	func (ltf *LiveTestFramework) Setup(t *testing.T) *TestFramework {
		return ltf.TestFramework.Setup(t)
	}

	type Aviatrix80TestFramework struct {
		*TestFramework
	}

	func (atf *Aviatrix80TestFramework) Setup(t *testing.T) *TestFramework {
		return atf.TestFramework.Setup(t)
	}

*/
