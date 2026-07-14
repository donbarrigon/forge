package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/donbarrigon/forge/env"
	"github.com/donbarrigon/forge/errs"
	"github.com/vmihailenco/msgpack/v5"
)

type Context struct {
	Request     *http.Request
	Writer      http.ResponseWriter
	MapRoutes   map[string]string // [name]method:path mapa de rutas para buscar por nombre
	Params      Params
	contentType ContentType // formato del body de la request (header Content-Type)
	accept      ContentType // formato esperado de la response (header Accept)
}

func NewContext(w http.ResponseWriter, r *http.Request, MapRoutes map[string]string, params Params) *Context {
	return &Context{
		Writer:      w,
		Request:     r,
		MapRoutes:   MapRoutes,
		Params:      params,
		contentType: ParseContentType(r.Header.Get("Content-Type")),
		accept:      ParseContentType(r.Header.Get("Accept")),
	}
}

// ContentType retorna el formato en que llegó el body de la request.
func (self *Context) ContentType() ContentType {
	return self.contentType
}

// Accept retorna el formato en el que el cliente espera la response.
func (self *Context) Accept() ContentType {
	return self.accept
}

// Lang retorna el código de idioma de la cabecera "Accept-Language".
// Example: "en" from "en-US,en;q=0.9", or "fr" from "fr-CH, fr;q=0.8".
func (self *Context) Lang() string {
	header := self.Request.Header.Get("Accept-Language")
	var tag string
	if header == "" {
		tag = env.App.Locale
	} else {
		// Take first entry (before comma or end of string)
		tag = strings.TrimSpace(strings.SplitN(header, ",", 2)[0])
		// Remove quality factor (e.g., ";q=0.9")
		if i := strings.Index(tag, ";"); i != -1 {
			tag = strings.TrimSpace(tag[:i])
		}
	}

	// Return primary subtag (before "-")
	if i := strings.Index(tag, "-"); i != -1 {
		return tag[:i]
	}
	return tag
}

// Locale retorna el locale completo de la cabecera "Accept-Language".
// Example: "en-US" from "en-US,en;q=0.9", or "fr-CH" from "fr-CH, fr;q=0.8".
func (self *Context) Locale() string {
	header := self.Request.Header.Get("Accept-Language")
	if header == "" {
		return env.App.Locale
	}
	// Take first entry (before comma or end of string)
	tag := strings.TrimSpace(strings.SplitN(header, ",", 2)[0])
	// Remove quality factor (e.g., ";q=0.9")
	if i := strings.Index(tag, ";"); i != -1 {
		tag = strings.TrimSpace(tag[:i])
	}
	return tag
}

// GetBody lee el cuerpo de la petición y lo devuelve en un slice de bytes.
func (self *Context) GetBody() ([]byte, *errs.Error) {
	if self.Request.Body == nil {
		return nil, nil
	}
	data, err := io.ReadAll(self.Request.Body)
	if err != nil {
		return nil, errs.BadRequestMsg("Fail to read request body", err)
	}
	return data, nil
}

// UnmarshalBody lee el cuerpo de la petición y lo decodifica en v segun el
// Content-Type de la request (self.contentType). Si el formato no tiene un
// decoder soportado retorna un error 415 (Unsupported Media Type).
func (self *Context) UnmarshalBody(v any) *errs.Error {
	data, e := self.GetBody()
	if e != nil {
		return e
	}

	switch self.contentType {
	case TYPE_MSGPACK:
		if err := msgpack.Unmarshal(data, v); err != nil {
			return errs.BadRequestMsg("Fail to decode msgpack request body", err)
		}
	case TYPE_JSON:
		if err := json.Unmarshal(data, v); err != nil {
			return errs.BadRequestMsg("Fail to decode json request body", err)
		}
	default:
		return errs.UnsupportedMediaType(nil)
	}

	return nil
}

// === Response: quien realmente serializa y escribe =========================

var (
	fallbackJSON    = []byte(`{"message":"failed to encode response"}`)
	fallbackMsgpack []byte
)

func init() {
	// se calcula una sola vez al cargar el paquete, no en cada request
	fallbackMsgpack, _ = msgpack.Marshal(map[string]string{"message": "failed to encode response"})
}

// Response serializa "a" segun el Accept del cliente (self.accept) y lo escribe
// en el ResponseWriter con el status dado. Si el Accept no tiene un encoder
// propio (ej. text/html, text/csv) hace fallback a msgpack.
func (self *Context) Response(status int, a any) {
	switch self.accept {
	case TYPE_JSON:
		self.writeJSON(status, a)
	default:
		self.writeMsgpack(status, a)
	}
}

func (self *Context) writeMsgpack(status int, a any) {
	body, encErr := msgpack.Marshal(a)
	self.Writer.Header().Set("Content-Type", TYPE_MSGPACK.String())
	if encErr != nil {
		self.Writer.WriteHeader(http.StatusInternalServerError)
		self.Writer.Write(fallbackMsgpack)
		return
	}

	self.Writer.WriteHeader(status)
	self.Writer.Write(body)
}

func (self *Context) writeJSON(status int, a any) {
	body, encErr := json.Marshal(a)
	self.Writer.Header().Set("Content-Type", TYPE_JSON.String())
	if encErr != nil {
		self.Writer.WriteHeader(http.StatusInternalServerError)
		self.Writer.Write(fallbackJSON)
		return
	}

	self.Writer.WriteHeader(status)
	self.Writer.Write(body)
}

// === Atajos: solo aportan el status code, delegan en Response ==============

func (self *Context) ResponseOk(a any) {
	self.Response(http.StatusOK, a)
}

func (self *Context) ResponseCreated(a any) {
	self.Response(http.StatusCreated, a)
}

func (self *Context) ResponseNoContent() {
	self.Writer.WriteHeader(http.StatusNoContent)
}

func (self *Context) ResponseError(e *errs.Error) {
	self.Response(e.Status, e)
}
