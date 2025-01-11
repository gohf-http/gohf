package gohf

type httpHandlerMap struct {
	table map[string]*httpHandler
}

func newHttpHandlerMap() *httpHandlerMap {
	return &httpHandlerMap{table: make(map[string]*httpHandler)}
}

func (httpHandlerMap *httpHandlerMap) addPatternString(pattern string) {
	if _, ok := httpHandlerMap.table[pattern]; !ok {
		httpHandlerMap.table[pattern] = newHttpHandler()
	}
}

func (httpHandlerMap *httpHandlerMap) getMap() map[string]*httpHandler {
	return httpHandlerMap.table
}

func (httpHandlerMap *httpHandlerMap) addHandlerFuncToPatternString(pattern string, handlerFunc HandlerFunc) {
	if _, ok := httpHandlerMap.table[pattern]; !ok {
		httpHandlerMap.table[pattern] = newHttpHandler()
	}
	httpHandlerMap.table[pattern].addHandlerFunc(handlerFunc)
}

func (httpHandlerMap *httpHandlerMap) addHandlerFuncToAll(handlerFunc HandlerFunc) {
	for _, handler := range httpHandlerMap.table {
		handler.addHandlerFunc(handlerFunc)
	}
}

func writeHandlerMap(target *httpHandlerMap, source *httpHandlerMap) {
	for pattern, httpHandler := range source.table {
		if httpHandler.len() > 0 {
			target.table[pattern] = httpHandler
		}
	}
}
