package env

import (
	"blog_app/log"
	"blog_app/mycontext"
	"errors"
	"os"
)

func GetEnv(ctx mycontext.Context, key string, logValue ...bool) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.GenericError(ctx, errors.New("key not found"), log.FieldsMap{"key": key})
	}
	if len(logValue) > 0 {
		if logValue[0] {
			log.GenericInfo(ctx, "Value Found", log.FieldsMap{"key": key, "value": value})
		}
	}
	return value
}
