package handler

import (
	"reflect"
	"runtime"
	"strings"
)

// === Path ===================================================================
// guarda la informacion de un endpoint
// ============================================================================
type Path struct {
	path        string
	method      string
	pathName    string
	controller  ContolerFunc
	middlewares []MiddlewareFunc
}

// Name renombra la ruta completa
func (self *Path) Name(name string) *Path {
	self.pathName = name
	return self
}

// Namep reemplaza solo el último segmento del nombre (después del último ".")
// por el string entrante. Si no hay punto, reemplaza el nombre completo.
func (self *Path) Namep(name string) *Path {
	if i := strings.LastIndex(self.pathName, "."); i != -1 {
		self.pathName = self.pathName[:i+1] + name
	} else {
		self.pathName = name
	}
	return self
}

type RouterBuilder struct {
	notFound    ContolerFunc // funcion de no encontrado
	prefixes    []string
	middlewares []MiddlewareFunc
	paths       []*Path
}

func NewRouterBuilder(notFoundFunc ContolerFunc) *RouterBuilder {
	if notFoundFunc == nil {
		notFoundFunc = DefaultNotFoundFunc
	}
	return &RouterBuilder{
		notFound:    notFoundFunc,
		prefixes:    []string{},
		middlewares: []MiddlewareFunc{},
		paths:       []*Path{},
	}
}

func DefaultNotFoundFunc(c *Context) {
	c.Writer.WriteHeader(404)
}

// crea un grupo de rutas con un prefijo dentro de la funcion fn
func (self *RouterBuilder) Prefix(prefix string, fn func()) {
	segments := strings.Split(strings.ToLower(strings.Trim(prefix, " /")), "/")
	self.prefixes = append(self.prefixes, segments...)
	fn()
	self.prefixes = self.prefixes[:len(self.prefixes)-len(segments)]
}

// crea un grupo de rutas con un middleware dentro de la funcion fn
func (self *RouterBuilder) Middleware(middleware MiddlewareFunc, fn func()) {
	self.middlewares = append(self.middlewares, middleware)
	fn()
	self.middlewares = self.middlewares[:len(self.middlewares)-1]
}

// crea un grupo de rutas con varios middlewares que se agreguen dentro de la funcion fn
func (self *RouterBuilder) Middlewares(middlewares []MiddlewareFunc, fn func()) {
	self.middlewares = append(self.middlewares, middlewares...)
	fn()
	self.middlewares = self.middlewares[:len(self.middlewares)-len(middlewares)]
}

// crea un grupo de rutas con un prefijo y un middleware dentro de la funcion fn
func (self *RouterBuilder) Group(prefix string, middleware MiddlewareFunc, fn func()) {
	segments := strings.Split(strings.ToLower(strings.Trim(prefix, " /")), "/")
	self.prefixes = append(self.prefixes, segments...)
	self.middlewares = append(self.middlewares, middleware)
	fn()
	self.prefixes = self.prefixes[:len(self.prefixes)-len(segments)]
	self.middlewares = self.middlewares[:len(self.middlewares)-1]
}

// crea un grupo de rutas con un prefijo y varios middlewares dentro de la funcion fn
func (self *RouterBuilder) Groups(prefix string, middlewares []MiddlewareFunc, fn func()) {
	segments := strings.Split(strings.ToLower(strings.Trim(prefix, " /")), "/")
	self.prefixes = append(self.prefixes, segments...)
	self.middlewares = append(self.middlewares, middlewares...)
	fn()
	self.prefixes = self.prefixes[:len(self.prefixes)-len(segments)]
	self.middlewares = self.middlewares[:len(self.middlewares)-len(middlewares)]
}

// === Métodos HTTP ===========================================================

func (self *RouterBuilder) Get(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "GET", controller, middlewares...)
}

func (self *RouterBuilder) Post(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "POST", controller, middlewares...)
}

func (self *RouterBuilder) Put(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "PUT", controller, middlewares...)
}

func (self *RouterBuilder) Patch(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "PATCH", controller, middlewares...)
}

func (self *RouterBuilder) Delete(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "DELETE", controller, middlewares...)
}

func (self *RouterBuilder) Options(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "OPTIONS", controller, middlewares...)
}

func (self *RouterBuilder) Head(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "HEAD", controller, middlewares...)
}

func (self *RouterBuilder) Connect(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "CONNECT", controller, middlewares...)
}

func (self *RouterBuilder) Trace(path string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	return self.append(path, "TRACE", controller, middlewares...)
}

// === Interno =================================================================

func (self *RouterBuilder) append(path string, method string, controller ContolerFunc, middlewares ...MiddlewareFunc) *Path {
	path = strings.ToLower(strings.Trim(path, " /"))
	newPath := strings.Join(append(self.prefixes, path), "/")

	cpMiddlewares := make([]MiddlewareFunc, 0, len(self.middlewares)+len(middlewares))
	cpMiddlewares = append(cpMiddlewares, self.middlewares...)
	cpMiddlewares = append(cpMiddlewares, middlewares...)

	p := &Path{
		path:        newPath,
		method:      method,
		pathName:    self.pathName(newPath, controller),
		controller:  controller,
		middlewares: cpMiddlewares,
	}
	self.paths = append(self.paths, p)
	return p
}

func (self *RouterBuilder) pathName(path string, controller ContolerFunc) string {
	name := strings.Replace(path, "/", ".", -1)

	fnName := runtime.FuncForPC(reflect.ValueOf(controller).Pointer()).Name()
	if i := strings.LastIndex(fnName, "."); i != -1 {
		fnName = fnName[i+1:]
	}
	fnName = strings.TrimSuffix(fnName, "-fm")

	return name + "." + strings.ToLower(fnName)
}

// newTrieNode crea un nodo del trie de rutas dinámicas con su Controller por
// defecto en NotFound, así una ruta parcial que no llegue a ser un endpoint
// real cae en el 404 en vez de un Controller nil.
func newTrieNode(notFound ContolerFunc) *Route {
	return &Route{
		Controller:   notFound,
		StaticRoutes: map[string]*Route{},
	}
}

// isDynamicSegment reporta si un segmento del path es un parámetro (:id,
// {varname}) o un wildcard (*, *filepath).
func isDynamicSegment(segment string) bool {
	return strings.HasPrefix(segment, ":") ||
		strings.HasPrefix(segment, "{") ||
		strings.HasPrefix(segment, "*")
}

// paramName limpia el segmento dinámico y devuelve solo el nombre del
// parámetro (":id" -> "id", "{userId}" -> "userId", "*filepath" -> "filepath").
func paramName(segment string) string {
	segment = strings.TrimPrefix(segment, ":")
	segment = strings.TrimPrefix(segment, "*")
	segment = strings.TrimPrefix(segment, "{")
	segment = strings.TrimSuffix(segment, "}")
	return segment
}

// standardizePath convierte los segmentos con prefijo ":" al formato "{...}",
// que es el estándar para MapRoutes. Ambas sintaxis (":id" y "{id}") se
// aceptan al registrar la ruta, pero en el mapa siempre quedan como "{id}".
func standardizePath(path string) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if strings.HasPrefix(seg, ":") {
			segments[i] = "{" + strings.TrimPrefix(seg, ":") + "}"
		}
	}
	return strings.Join(segments, "/")
}

// Build construye el Router final a partir de la información de rutas y
// middlewares guardada en el RouterBuilder. Se llama al final de la
// configuración de rutas.
func (self *RouterBuilder) Build() *Router {

	router := &Router{
		NotFound:     self.notFound,
		MapRoutes:    make(map[string]string, len(self.paths)),
		StaticRoutes: make(map[string]*Route, len(self.paths)),
	}

	for _, p := range self.paths {
		router.MapRoutes[p.pathName] = p.method + ":" + standardizePath(p.path)

		segments := strings.Split(p.path, "/")

		dynamic := false
		for _, seg := range segments {
			if isDynamicSegment(seg) {
				dynamic = true
				break
			}
		}

		// --- ruta estática: va directo al map, sin pasar por el trie ---
		if !dynamic {
			key := strings.Join(segments, "/") + "/" + p.method
			router.StaticRoutes[key] = &Route{
				Controller:  p.controller,
				Middlewares: p.middlewares,
			}
			continue
		}

		// --- ruta dinámica: se recorre/crea el trie ---
		if router.DynamicRoutes == nil {
			router.DynamicRoutes = newTrieNode(router.NotFound)
		}

		node := router.DynamicRoutes
		var paramNames []string
		isWildcard := false

		for _, seg := range segments {
			switch {
			case strings.HasPrefix(seg, "*"):
				// el nodo actual (el padre) es quien captura todo lo que
				// siga; no se crea un hijo nuevo para el "*" en sí.
				node.IsFilePath = true
				if name := paramName(seg); name != "" {
					paramNames = append(paramNames, name)
				}
				isWildcard = true

			case strings.HasPrefix(seg, ":"), strings.HasPrefix(seg, "{"):
				paramNames = append(paramNames, paramName(seg))
				if node.DynamicRoutes == nil {
					node.DynamicRoutes = newTrieNode(router.NotFound)
				}
				node = node.DynamicRoutes

			default:
				if node.StaticRoutes[seg] == nil {
					node.StaticRoutes[seg] = newTrieNode(router.NotFound)
				}
				node = node.StaticRoutes[seg]
			}

			if isWildcard {
				break // nada debería venir después de un wildcard
			}
		}

		// el método se agrega como último segmento literal del trie, salvo
		// en wildcard: ahí el nodo padre ya quedó marcado como IsFilePath
		// y Find() corta ahí sin seguir bajando por el método.
		if !isWildcard {
			if node.StaticRoutes[p.method] == nil {
				node.StaticRoutes[p.method] = newTrieNode(router.NotFound)
			}
			node = node.StaticRoutes[p.method]
		}

		node.Controller = p.controller
		node.Middlewares = p.middlewares
		node.params = paramNames
	}

	return router
}
