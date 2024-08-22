package log

import (
	"blog_app/mycontext"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	INFO = iota
	HTTP
	ERROR
	TRACE
	WARNING
)

var (
	setLevel = WARNING
	trace    *log.Logger
	info     *log.Logger
	warning  *log.Logger
	httplog  *log.Logger
	errorlog *log.Logger
)

const (
	runModeType      = "RUN_MODE_TYPE"
	runModeTypeLocal = "local"
)

type FieldsMap map[string]interface{}

func init() {
	logInit(os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stderr)
}

func logInit(traceHandle, infoHandle, warningHandle, httpHandle, errorHandle io.Writer) {

	flagWithClusterType := log.LUTC | log.LstdFlags | log.Lshortfile
	flagWithoutClusterType := log.LUTC | log.LstdFlags

	var flag int

	if os.Getenv(runModeType) == runModeTypeLocal {
		flag = flagWithClusterType
	} else {
		flag = flagWithoutClusterType
	}

	trace = log.New(traceHandle, "TRACE|", flag)

	info = log.New(infoHandle, "INFO|", flag)

	warning = log.New(warningHandle, "WARNING|", flag)

	httplog = log.New(httpHandle, "HTTP|", flag)

	errorlog = log.New(errorHandle, "ERROR|", flagWithClusterType)
}

func doLog(cLog *log.Logger, level, callDepth int, v ...interface{}) {
	if level <= setLevel {
		if level == ERROR {
			cLog.SetOutput(os.Stderr)
			cLog.SetFlags(log.Llongfile)
		}
		//cLog.SetOutput(os.Stdout)
		cLog.Output(callDepth, fmt.Sprintln(v...))
	}
}

func generatePrefix(ctx mycontext.Context) string {
	return strings.Join([]string{ctx.UserName, ctx.UserEmail}, ":")
}

func generateTrackingIDs(ctx mycontext.Context) (retString string) {
	requestID := ctx.RequestID

	if requestID != "" {
		retString = "requestId=" + requestID
	}

	return
}

func traceLog(v ...interface{}) {
	doLog(trace, TRACE, 4, v...)
}

func infoLog(v ...interface{}) {
	doLog(info, INFO, 4, v...)
}

func warningLog(v ...interface{}) {
	doLog(warning, WARNING, 4, v...)
}

func errorLog(v ...interface{}) {
	doLog(errorlog, ERROR, 4, v...)
}

func GenericTrace(ctx mycontext.Context, traceMessage string, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	if os.Getenv("TEK_SERVICE_TRACE") == "true" {
		prefix := generatePrefix(ctx)
		trackingIDs := generateTrackingIDs(ctx)
		msg := fmt.Sprintf("|%s|%s|",
			prefix,
			trackingIDs)
		if fields != nil && len(fields) > 0 {
			fieldsBytes, _ := json.Marshal(fields)
			fieldsString := string(fieldsBytes)
			traceLog(msg, traceMessage, "|", fieldsString)
		} else {
			traceLog(msg, traceMessage)
		}
	}
}

func GenericInfo(ctx mycontext.Context, infoMessage string, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	fieldsBytes, _ := json.Marshal(fields)
	fieldsString := string(fieldsBytes)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		infoLog(msg, infoMessage, "|", fieldsString)
	} else {
		infoLog(msg, infoMessage)
	}

}

func GenericWarning(ctx mycontext.Context, warnMessage string, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	if os.Getenv("TEK_SERVICE_WARN") == "true" {
		prefix := generatePrefix(ctx)
		trackingIDs := generateTrackingIDs(ctx)
		msg := fmt.Sprintf("|%s|%s|",
			prefix,
			trackingIDs)
		if fields != nil && len(fields) > 0 {
			fieldsBytes, _ := json.Marshal(fields)
			fieldsString := string(fieldsBytes)
			warningLog(msg, warnMessage, "|", fieldsString)
		} else {
			warningLog(msg, warnMessage)
		}
	}
}

func GenericError(ctx mycontext.Context, e error, data ...FieldsMap) {
	var fields FieldsMap
	if len(data) > 0 {
		fields = data[0]
	}
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		errorLog(msg, "|", fieldsString)
	} else {
		errorLog(msg)
	}
}
