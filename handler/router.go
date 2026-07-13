package handler

import (
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/donbarrigon/forge/err"
)

type ContolerFunc func(ctx *Context)
type MiddlewareFunc func(ctx *Context) *err.HttpError

type Param struct {
	Key   string
	Value string
}

type Params []Param

type Route struct {
	Controller    ContolerFunc      // la funcion que maneja esta ruta
	Middlewares   []MiddlewareFunc  // los middlewares de esta ruta
	params        []string          // los nombres de los parametros dinamicos de esta ruta
	IsFilePath    bool              // true si todas las siguientes rutas sin wildcards de esta
	StaticRoutes  map[string]*Route // [segment]Route siguientes rutas estaticas
	DynamicRoutes *Route            // siguiente ruta dinamica
}

type Router struct {
	NotFound      ContolerFunc      // la funcion que maneja las rutas no encontradas
	MapRoutes     map[string]string // [name]method:path mapa de rutas para buscar por nombre
	StaticRoutes  map[string]*Route // [path/method]Route rutas sin params
	DynamicRoutes *Route            // todas las rutas dinamicas y wildcard guardadas en ramas
}

func (self *Router) Find(url string, method string) (*Route, Params) {
	path := strings.ToLower(strings.Trim(url, " /")) + "/" + method

	// primero busco en las rutas staticas
	if route, ok := self.StaticRoutes[path]; ok {
		return route, []Param{}
	}

	// si no se encuentra en las rutas staticas, busco en las rutas dinamicas
	if self.DynamicRoutes == nil {
		return nil, nil
	}

	route := self.DynamicRoutes
	params := []Param{}
	segments := strings.Split(path, "/")
	// numSegments := len(segments)
	for i, segment := range segments {
		// si el siguiente segmento es un * filePath, se retorna y los segmentos siguientes va a params
		if route.IsFilePath {
			nextSegments := segments[i:]
			for j, s := range nextSegments {
				params = append(params, Param{Key: strconv.Itoa(j), Value: s})
			}
			return route, params
		}

		// busco si el segmento es estatico
		if r := route.StaticRoutes[segment]; r != nil {
			route = r
			continue
		}

		if route.DynamicRoutes == nil {
			return nil, nil
		}

		// tomo el segmento dinamico y lo agrego a los params
		route = route.DynamicRoutes
		params = append(params, Param{Key: "", Value: segment})

	}
	return route, params
}

func (self *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params := self.Find(r.URL.Path, r.Method)
	c := NewContext(w, r, self.MapRoutes, params)

	defer func() {
		if rec := recover(); rec != nil {
			c.ResponseError(err.Panic(rec, string(debug.Stack())))
		}
	}()

	if route == nil {
		self.NotFound(c)
		return
	}

	for i, key := range route.params {
		params[i].Key = key
	}

	for _, middleware := range route.Middlewares {
		if e := middleware(c); e != nil {
			c.ResponseError(e)
			return
		}
	}

	route.Controller(c)
}
