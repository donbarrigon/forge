package handler

import "strings"

type ContentType uint8

const (
	TYPE_MSGPACK      ContentType = iota // application/msgpack
	TYPE_JSON                            // application/json
	TYPE_HTML                            // text/html
	TYPE_CSV                             // text/csv
	TYPE_XML                             // application/xml
	TYPE_TEXT                            // text/plain
	TYPE_FORM                            // application/x-www-form-urlencoded
	TYPE_MULTIPART                       // multipart/form-data
	TYPE_JS                              // text/javascript
	TYPE_CSS                             // text/css
	TYPE_PDF                             // application/pdf
	TYPE_OCTET_STREAM                    // application/octet-stream
	TYPE_YAML                            // application/x-yaml
	TYPE_EVENT_STREAM                    // text/event-stream (SSE)
)

// mimeByType mapea cada ContentType a su string MIME.
// el indice del array coincide con el valor de la constante,
// por eso el orden de las constantes de arriba no se debe alterar
// (agregar nuevas siempre al final).
var mimeByType = [...]string{
	TYPE_MSGPACK:      "application/msgpack",
	TYPE_JSON:         "application/json",
	TYPE_HTML:         "text/html",
	TYPE_CSV:          "text/csv",
	TYPE_XML:          "application/xml",
	TYPE_TEXT:         "text/plain",
	TYPE_FORM:         "application/x-www-form-urlencoded",
	TYPE_MULTIPART:    "multipart/form-data",
	TYPE_JS:           "text/javascript",
	TYPE_CSS:          "text/css",
	TYPE_PDF:          "application/pdf",
	TYPE_OCTET_STREAM: "application/octet-stream",
	TYPE_YAML:         "application/x-yaml",
	TYPE_EVENT_STREAM: "text/event-stream",
}

// typeByMime es el mapa inverso, se construye una sola vez al cargar el paquete.
var typeByMime map[string]ContentType

func init() {
	typeByMime = make(map[string]ContentType, len(mimeByType))
	for t, mime := range mimeByType {
		typeByMime[mime] = ContentType(t)
	}
}

// String devuelve el mime string correspondiente a un ContentType.
// si el valor esta fuera de rango el fallback es "application/msgpack".
func (self ContentType) String() string {
	if int(self) < len(mimeByType) {
		return mimeByType[self]
	}
	return mimeByType[TYPE_MSGPACK]
}

// ParseContentType parsea un header tipo "application/json; charset=utf-8" y devuelve
// el ContentType correspondiente. si no lo encuentra el fallback es TYPE_MSGPACK (0).
func ParseContentType(s string) ContentType {
	if i := strings.IndexByte(s, ';'); i != -1 {
		s = s[:i]
	}
	s = strings.TrimSpace(strings.ToLower(s))

	if t, ok := typeByMime[s]; ok {
		return t
	}
	return TYPE_MSGPACK
}
