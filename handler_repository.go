package gohf

type handlerRepository struct {
	handlers []*handler
}

func newHandlerRepository() *handlerRepository {
	return &handlerRepository{}
}

func (hr *handlerRepository) addHandler(f HandlerFunc, owner *Router, pattern string, all bool) {
	hr.handlers = append(hr.handlers, newHandler(f, owner, pattern, all))
}

func (hr *handlerRepository) getHandlers() []*handler {
	return hr.handlers
}

type handler struct {
	f       HandlerFunc
	owner   *Router
	pattern string
	all     bool
}

func newHandler(f HandlerFunc, owner *Router, pattern string, all bool) *handler {
	return &handler{
		f:       f,
		owner:   owner,
		pattern: pattern,
		all:     all,
	}
}
