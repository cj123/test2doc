package test

import (
	"github.com/cj123/test2doc/doc/parse"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

var infoFunc = DefaultHandlerInfo

type HandlerInfo struct {
	FuncName string
}

type handlerInfoFunc func(r *http.Request) *HandlerInfo

func RegisterHandlerInfoFunc(fn handlerInfoFunc) {
	infoFunc = fn
}

func DefaultHandlerInfo(r *http.Request) *HandlerInfo {
	i := 1
	max := 15

	var pc uintptr
	var fnName string
	var ok, fnInPkg, sawPkg bool

	// iterate until we find the top level func in this pkg (the handler)
	for i < max {
		pc, _, _, ok = runtime.Caller(i)
		if !ok {
			log.Println("test2doc: DefaultHandlerInfo: !ok")
			return nil
		}

		fn := runtime.FuncForPC(pc)
		fnName = fn.Name()

		fnInPkg = parse.IsFuncInPkg(fnName)
		if sawPkg && !fnInPkg {
			pc, _, _, ok = runtime.Caller(i - 1)
			fn := runtime.FuncForPC(pc)
			fnName = fn.Name()
			break
		}

		sawPkg = fnInPkg
		i++
	}

	return &HandlerInfo{
		FuncName: fnName,
	}
}

// GorillaMuxHandlerInfo takes a mux.Router and finds handler info from it.
func GorillaMuxHandlerInfo(router *mux.Router) func(r *http.Request) *HandlerInfo {
	return func(r *http.Request) *HandlerInfo {
		var handlerInfo *HandlerInfo

		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			var match mux.RouteMatch

			if route.Match(r, &match) {
				fn := runtime.FuncForPC(reflect.ValueOf(match.Handler).Pointer())
				handlerInfo = &HandlerInfo{FuncName: strings.Replace(fn.Name(), "-fm", " ", 1)}
				return nil
			}

			return nil
		})

		return handlerInfo
	}
}
