package handler

import (
	"encoding/json"
	"net/http"

	"github.com/donbarrigon/forge/errs"
	"github.com/vmihailenco/msgpack/v5"
)

type Context struct {
	Request   *http.Request
	Writer    http.ResponseWriter
	MapRoutes map[string]string // [name]method:path mapa de rutas para buscar por nombre
	Params    Params
}

func NewContext(w http.ResponseWriter, r *http.Request, MapRoutes map[string]string, params Params) *Context {
	return &Context{
		Writer:    w,
		Request:   r,
		MapRoutes: MapRoutes,
		Params:    params,
	}
}

var (
	fallbackJSON    = []byte(`{"message":"failed to encode response"}`)
	fallbackMsgpack []byte
)

func init() {
	// se calcula una sola vez al cargar el paquete, no en cada request
	fallbackMsgpack, _ = msgpack.Marshal(map[string]string{"message": "failed to encode response"})
}

// === Response: quienes realmente serializan y escriben =====================

func (self *Context) Response(status int, a any) {
	body, encErr := msgpack.Marshal(a)
	self.Writer.Header().Set("Content-Type", "application/msgpack")
	if encErr != nil {
		self.Writer.WriteHeader(http.StatusInternalServerError)
		self.Writer.Write(fallbackMsgpack)
		return
	}

	self.Writer.WriteHeader(status)
	self.Writer.Write(body)
}

func (self *Context) ResponseJson(status int, a any) {
	body, encErr := json.Marshal(a)
	self.Writer.Header().Set("Content-Type", "application/json")
	if encErr != nil {
		self.Writer.WriteHeader(http.StatusInternalServerError)
		self.Writer.Write(fallbackJSON)
		return
	}

	self.Writer.WriteHeader(status)
	self.Writer.Write(body)
}

// === Atajos: solo aportan el status code, delegan en Response/ResponseJson ===

func (self *Context) ResponseOk(a any) {
	self.Response(http.StatusOK, a)
}

func (self *Context) ResponseCreated(a any) {
	self.Response(http.StatusCreated, a)
}

func (self *Context) ResponseNoContent() {
	self.Writer.WriteHeader(http.StatusNoContent)
}

func (self *Context) ResponseOkJson(a any) {
	self.ResponseJson(http.StatusOK, a)
}

func (self *Context) ResponseCreatedJson(a any) {
	self.ResponseJson(http.StatusCreated, a)
}

func (self *Context) ResponseError(e *errs.Error) {
	self.Response(e.Status, e)
}

func (self *Context) ResponseErrorJson(e *errs.Error) {
	self.ResponseJson(e.Status, e)
}
