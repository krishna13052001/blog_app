package constants

import (
	"blog_app/env"
	"blog_app/mycontext"
)

var (
	ServiceName    string
	ServicePort    string
	ServiceRoute   string
	MongoHost      string
	MongoHostAtlas string
)

func LoadEnv(ctx mycontext.Context) {
	ServiceName = env.GetEnv(ctx, "SERVICE_NAME")
	ServicePort = env.GetEnv(ctx, "SERVICE_PORT")
	ServiceRoute = env.GetEnv(ctx, "SERVICE_ROUTE")
	MongoHost = env.GetEnv(ctx, "MONGO_HOST")
	MongoHostAtlas = env.GetEnv(ctx, "MONGO_ATLAS")
}
