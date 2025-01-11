package gohf

import (
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

func (r *Router) SubRouter(pattern string) *Router {
	parsedPattern, err := parsePattern(pattern)
	if err != nil {
		panic(err)
	}

	pat, err := mergePattern(r.pattern, parsedPattern)
	if err != nil {
		panic(err)
	}
	if len(pat.path) > 0 && !strings.HasSuffix(pat.path, "/") {
		pat.path += "/"
	}

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
	localHandlerMap.addPatternString(r.pattern.String())

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
