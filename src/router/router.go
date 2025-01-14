package router

import (
	"context"
	"log"
	"net/http"
	"regexp"
)

type RouterEntry struct {
	Path *regexp.Regexp
	Method string
	HandlerFunc http.HandlerFunc
}

type Router struct {
	routes []RouterEntry
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR:", r)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
  }()
	for _, e := range rtr.routes {
    params := e.Match(r)
    if params == nil {
     continue
    }
    ctx := context.WithValue(r.Context(), "params", params)
    e.HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
    return
	}
	http.NotFound(w, r)
}

func (rtr *Router) Route(method, path string, handlerFunc http.HandlerFunc) {
	e := RouterEntry{
		Method: method,
		Path: regexp.MustCompile(path),
		HandlerFunc: handlerFunc,
	}
	rtr.routes = append(rtr.routes, e)
}

func (re *RouterEntry) Match(r *http.Request) map[string]string {
	match := re.Path.FindStringSubmatch(r.URL.Path)
	if match == nil {
		return nil
	}
	params := make(map[string]string)
	groupNames := re.Path.SubexpNames()
	for i, group := range match {
		params[groupNames[i]] = group
	}
	return params
}

func GetURLParam(r *http.Request, name string) string {
	ctx := r.Context()
	params := ctx.Value("params").(map[string]string)
	return params[name]
}