package mycontext

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type MyContext struct {
	RequestID string
	ApiToken  string
	UserName  string
	UserEmail string
}

type Context struct {
	context.Context
	MyContext
}

type customContextType string

const (
	MyCtx customContextType = "myCtx"
)

func NewContext() Context {
	return Context{
		Context:   context.Background(),
		MyContext: MyContext{},
	}
}

func GetMyCtx(ctx context.Context) (MyContext, bool) {
	if ctx == nil {
		return MyContext{}, false
	}
	myCtx, exists := ctx.Value(MyCtx).(MyContext)
	return myCtx, exists
}

func WithCtx(ctx context.Context, myctx MyContext) context.Context {
	return context.WithValue(ctx, MyCtx, myctx)
}

func UpgradeCtx(ctx context.Context) Context {
	var myContext Context
	myCtx, _ := GetMyCtx(ctx)

	myContext.Context = ctx
	myContext.MyContext = myCtx
	return myContext
}

func New(id ...string) Context {
	var requestID string
	if len(id) > 0 {
		requestID = id[0]
	}
	if len(requestID) == 0 {
		requestID = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	myCtx := MyContext{
		RequestID: requestID,
	}
	ctx := UpgradeCtx(WithCtx(context.Background(), myCtx))
	return ctx
}

func CopyContext(ctx context.Context) Context {
	myCtx, _ := GetMyCtx(ctx)
	return Context{
		Context:   context.Background(),
		MyContext: myCtx,
	}
}

func (c *Context) GetUserEmail() string {
	return c.UserEmail
}

func (c *Context) GetUserName() string {
	return c.UserName
}

func (c *Context) GetRequestId() string {
	return c.RequestID
}
