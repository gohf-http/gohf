package gohf

import (
	"net/http"
	"strings"
)

type Router struct {
	prefix            string
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
		r.handlerRepository.addHandler(f, r, "", true)
	}
}

func (r *Router) Handle(pattern string, handlerFuncs ...HandlerFunc) {
	for _, f := range handlerFuncs {
		r.handlerRepository.addHandler(f, r, pattern, false)
	}
}

func (r *Router) SubRouter(prefix string) *Router {
	prefix = strings.TrimSuffix(prefix, "/")

	subRouter := &Router{
		prefix:            prefix,
		handlerRepository: r.handlerRepository,
		parentRouter:      r,
	}

	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r *Router) CreateServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	httpHandlerMap := r.getHttpHandlerMap()

	for pattern, httpHandler := range httpHandlerMap {
		if httpHandler.len() > 0 {
			mux.Handle(pattern, httpHandler)
		}
	}

	for _, router := range r.subRouters {
		prefix := router.prefix
		mux.Handle(prefix+"/", http.StripPrefix(prefix, router.CreateServeMux()))
	}

	return mux
}

func (r *Router) getHttpHandlerMap() map[string]*httpHandler {
	httpHandlerMap := make(map[string]*httpHandler)
	httpHandlerMap["/"] = newHttpHandler()

	for _, h := range r.handlerRepository.getHandlers() {
		if !h.all && h.owner == r {
			if _, ok := httpHandlerMap[h.pattern]; !ok {
				httpHandlerMap[h.pattern] = newHttpHandler()
			}
		}
	}

	for _, handler := range r.handlerRepository.getHandlers() {
		if handler.all && (handler.owner == r || r.hasAncestor(handler.owner)) {
			for pattern := range httpHandlerMap {
				httpHandlerMap[pattern].addHandlerFunc(handler.f)
			}
		} else if !handler.all && handler.owner == r {
			pattern := handler.pattern
			httpHandlerMap[pattern].addHandlerFunc(handler.f)
		}
	}

	return httpHandlerMap
}

func (r *Router) hasAncestor(router *Router) bool {
	parentRouter := r.parentRouter
	if parentRouter == router {
		return true
	}
	return parentRouter != nil && parentRouter.hasAncestor(router)
}
