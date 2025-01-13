package gohf

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	pattern           pattern
	handlerRepository *handlerRepository
	parentRouter      *Router
	subRouters        []*Router
}

func New() *Router {
	return &Router{
		handlerRepository: newHandlerRepository(),
	}
}

func (r *Router) Use(handlerFuncs ...HandlerFunc) {
	for _, f := range handlerFuncs {
		if f == nil {
			continue
		}
		r.handlerRepository.addHandler(f, r, r.pattern, true)
	}
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	for _, f := range handlerFuncs {
		if f == nil {
			continue
		}
		parsedPattern, err := parsePattern(pattern)
		if err != nil {
			panic(err)
		}

		pat, err := mergePattern(r.pattern, parsedPattern)
		if err != nil {
			panic(err)
		}

		r.handlerRepository.addHandler(f, r, pat, false)
	}
}

func (r *Router) GET(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodGet, pattern), handlerFuncs...)
}

func (r *Router) POST(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodPost, pattern), handlerFuncs...)
}

func (r *Router) PUT(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodPut, pattern), handlerFuncs...)
}

func (r *Router) PATCH(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodPatch, pattern), handlerFuncs...)
}

func (r *Router) DELETE(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodDelete, pattern), handlerFuncs...)
}

func (r *Router) OPTIONS(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodOptions, pattern), handlerFuncs...)
}

func (r *Router) HEAD(pattern string, handlerFuncs ...HandlerFunc) {
	r.Handle(fmt.Sprintf("%s %s", http.MethodHead, pattern), handlerFuncs...)
}

func (r *Router) SubRouter(pattern string) *Router {
	parsedPattern, err := parsePattern(pattern)
	if err != nil {
		panic(err)
	}

	pat, err := mergePattern(r.pattern, parsedPattern)
	if err != nil {
		panic(err)
	}
	pat.path = strings.TrimSuffix(pat.path, "/")

	subRouter := &Router{
		pattern:           pat,
		handlerRepository: r.handlerRepository,
		parentRouter:      r,
	}

	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r *Router) CreateServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	httpHandlerMap := r.getHttpHandlerMap()

	for patternString, httpHandler := range httpHandlerMap.getMap() {
		if httpHandler.len() > 0 {
			mux.Handle(patternString, httpHandler)
		}
	}

	return mux
}

func (r *Router) hasAncestor(router *Router) bool {
	parentRouter := r.parentRouter
	if parentRouter == router {
		return true
	}
	return parentRouter != nil && parentRouter.hasAncestor(router)
}

func (r *Router) getHttpHandlerMap() *httpHandlerMap {
	httpHandlerMap := newHttpHandlerMap()
	r.setupHttpHandlerMap(httpHandlerMap)
	return httpHandlerMap
}

func (r *Router) setupHttpHandlerMap(targetHttpHandlerMap *httpHandlerMap) {
	localHandlerMap := newHttpHandlerMap()
	routerPatternString := r.pattern.String()
	if !strings.HasSuffix(routerPatternString, "/") {
		routerPatternString += "/"
	}
	localHandlerMap.addPatternString(routerPatternString)

	for _, h := range r.handlerRepository.getHandlers() {
		if !h.all && h.owner == r {
			localHandlerMap.addPatternString(h.pattern.String())
		}
	}

	for _, handler := range r.handlerRepository.getHandlers() {
		if handler.all && (handler.owner == r || r.hasAncestor(handler.owner)) {
			localHandlerMap.addHandlerFuncToAll(handler.f)
		} else if !handler.all && handler.owner == r {
			patternString := handler.pattern.String()
			localHandlerMap.addHandlerFuncToPatternString(patternString, handler.f)
		}
	}

	writeHandlerMap(targetHttpHandlerMap, localHandlerMap)

	for _, subRouter := range r.subRouters {
		subRouter.setupHttpHandlerMap(targetHttpHandlerMap)
	}
}
