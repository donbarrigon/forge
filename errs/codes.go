package errs

// Códigos de error HTTP 4xx y 5xx
const (
	// --- 4xx: Errores del cliente ---
	BAD_REQUEST                     = 400
	UNAUTHORIZED                    = 401
	PAYMENT_REQUIRED                = 402
	FORBIDDEN                       = 403
	NOT_FOUND                       = 404
	METHOD_NOT_ALLOWED              = 405
	NOT_ACCEPTABLE                  = 406
	PROXY_AUTHENTICATION_REQUIRED   = 407
	REQUEST_TIMEOUT                 = 408
	CONFLICT                        = 409
	GONE                            = 410
	LENGTH_REQUIRED                 = 411
	PRECONDITION_FAILED             = 412
	PAYLOAD_TOO_LARGE               = 413
	URI_TOO_LONG                    = 414
	UNSUPPORTED_MEDIA_TYPE          = 415
	RANGE_NOT_SATISFIABLE           = 416
	EXPECTATION_FAILED              = 417
	IM_A_TEAPOT                     = 418
	MISDIRECTED_REQUEST             = 421
	UNPROCESSABLE_ENTITY            = 422
	LOCKED                          = 423
	FAILED_DEPENDENCY               = 424
	TOO_EARLY                       = 425
	UPGRADE_REQUIRED                = 426
	PRECONDITION_REQUIRED           = 428
	TOO_MANY_REQUESTS               = 429
	REQUEST_HEADER_FIELDS_TOO_LARGE = 431
	UNAVAILABLE_FOR_LEGAL_REASONS   = 451

	// --- 5xx: Errores del servidor ---
	INTERNAL                        = 500
	NOT_IMPLEMENTED                 = 501
	BAD_GATEWAY                     = 502
	SERVICE_UNAVAILABLE             = 503
	GATEWAY_TIMEOUT                 = 504
	HTTP_VERSION_NOT_SUPPORTED      = 505
	VARIANT_ALSO_NEGOTIATES         = 506
	INSUFFICIENT_STORAGE            = 507
	LOOP_DETECTED                   = 508
	NOT_EXTENDED                    = 510
	NETWORK_AUTHENTICATION_REQUIRED = 511
)

// Mensajes legibles para cada código.
const (
	// --- 4xx: Errores del cliente ---
	BAD_REQUEST_MSG                     = "Something's wrong with your request"
	UNAUTHORIZED_MSG                    = "You need to sign in to continue"
	PAYMENT_REQUIRED_MSG                = "Payment is required to access this"
	FORBIDDEN_MSG                       = "You don't have permission to do this"
	NOT_FOUND_MSG                       = "We couldn't find what you're looking for"
	METHOD_NOT_ALLOWED_MSG              = "This action isn't allowed"
	NOT_ACCEPTABLE_MSG                  = "We can't process your request in this format"
	PROXY_AUTHENTICATION_REQUIRED_MSG   = "Proxy authentication is required"
	REQUEST_TIMEOUT_MSG                 = "The request took too long"
	CONFLICT_MSG                        = "There's a conflict with the existing data"
	GONE_MSG                            = "This resource is no longer available"
	LENGTH_REQUIRED_MSG                 = "The request is missing size information"
	PRECONDITION_FAILED_MSG             = "The required conditions weren't met"
	PAYLOAD_TOO_LARGE_MSG               = "The file or data is too large"
	URI_TOO_LONG_MSG                    = "The address is too long"
	UNSUPPORTED_MEDIA_TYPE_MSG          = "This file type isn't supported"
	RANGE_NOT_SATISFIABLE_MSG           = "The requested range isn't valid"
	EXPECTATION_FAILED_MSG              = "We couldn't meet what was expected"
	IM_A_TEAPOT_MSG                     = "I'm a teapot, I can't brew coffee"
	MISDIRECTED_REQUEST_MSG             = "Your request was sent to the wrong place"
	UNPROCESSABLE_ENTITY_MSG            = "We couldn't process the information you sent"
	LOCKED_MSG                          = "This resource is locked"
	FAILED_DEPENDENCY_MSG               = "Something we needed didn't work correctly"
	TOO_EARLY_MSG                       = "It's too soon for this request"
	UPGRADE_REQUIRED_MSG                = "You need to update your app"
	PRECONDITION_REQUIRED_MSG           = "Some prerequisites are missing"
	TOO_MANY_REQUESTS_MSG               = "You've made too many requests, please wait a moment"
	REQUEST_HEADER_FIELDS_TOO_LARGE_MSG = "Your request headers are too large"
	UNAVAILABLE_FOR_LEGAL_REASONS_MSG   = "Unavailable for legal reasons"

	// --- 5xx: Errores del servidor ---
	INTERNAL_MSG                        = "Something went wrong on our end"
	NOT_IMPLEMENTED_MSG                 = "This feature isn't available yet"
	BAD_GATEWAY_MSG                     = "There was a problem with our servers"
	SERVICE_UNAVAILABLE_MSG             = "The service isn't available right now"
	GATEWAY_TIMEOUT_MSG                 = "The server took too long to respond"
	HTTP_VERSION_NOT_SUPPORTED_MSG      = "Your browser version isn't supported"
	VARIANT_ALSO_NEGOTIATES_MSG         = "There was a server negotiation error"
	INSUFFICIENT_STORAGE_MSG            = "Not enough storage space available"
	LOOP_DETECTED_MSG                   = "An infinite loop was detected on the server"
	NOT_EXTENDED_MSG                    = "The server is missing required extensions"
	NETWORK_AUTHENTICATION_REQUIRED_MSG = "You need to authenticate on the network"
)

var StatusMap = map[int]string{
	// --- 4xx: Errores del cliente ---
	BAD_REQUEST:                     BAD_REQUEST_MSG,
	UNAUTHORIZED:                    UNAUTHORIZED_MSG,
	PAYMENT_REQUIRED:                PAYMENT_REQUIRED_MSG,
	FORBIDDEN:                       FORBIDDEN_MSG,
	NOT_FOUND:                       NOT_FOUND_MSG,
	METHOD_NOT_ALLOWED:              METHOD_NOT_ALLOWED_MSG,
	NOT_ACCEPTABLE:                  NOT_ACCEPTABLE_MSG,
	PROXY_AUTHENTICATION_REQUIRED:   PROXY_AUTHENTICATION_REQUIRED_MSG,
	REQUEST_TIMEOUT:                 REQUEST_TIMEOUT_MSG,
	CONFLICT:                        CONFLICT_MSG,
	GONE:                            GONE_MSG,
	LENGTH_REQUIRED:                 LENGTH_REQUIRED_MSG,
	PRECONDITION_FAILED:             PRECONDITION_FAILED_MSG,
	PAYLOAD_TOO_LARGE:               PAYLOAD_TOO_LARGE_MSG,
	URI_TOO_LONG:                    URI_TOO_LONG_MSG,
	UNSUPPORTED_MEDIA_TYPE:          UNSUPPORTED_MEDIA_TYPE_MSG,
	RANGE_NOT_SATISFIABLE:           RANGE_NOT_SATISFIABLE_MSG,
	EXPECTATION_FAILED:              EXPECTATION_FAILED_MSG,
	IM_A_TEAPOT:                     IM_A_TEAPOT_MSG,
	MISDIRECTED_REQUEST:             MISDIRECTED_REQUEST_MSG,
	UNPROCESSABLE_ENTITY:            UNPROCESSABLE_ENTITY_MSG,
	LOCKED:                          LOCKED_MSG,
	FAILED_DEPENDENCY:               FAILED_DEPENDENCY_MSG,
	TOO_EARLY:                       TOO_EARLY_MSG,
	UPGRADE_REQUIRED:                UPGRADE_REQUIRED_MSG,
	PRECONDITION_REQUIRED:           PRECONDITION_REQUIRED_MSG,
	TOO_MANY_REQUESTS:               TOO_MANY_REQUESTS_MSG,
	REQUEST_HEADER_FIELDS_TOO_LARGE: REQUEST_HEADER_FIELDS_TOO_LARGE_MSG,
	UNAVAILABLE_FOR_LEGAL_REASONS:   UNAVAILABLE_FOR_LEGAL_REASONS_MSG,

	// --- 5xx: Errores del servidor ---
	INTERNAL:                        INTERNAL_MSG,
	NOT_IMPLEMENTED:                 NOT_IMPLEMENTED_MSG,
	BAD_GATEWAY:                     BAD_GATEWAY_MSG,
	SERVICE_UNAVAILABLE:             SERVICE_UNAVAILABLE_MSG,
	GATEWAY_TIMEOUT:                 GATEWAY_TIMEOUT_MSG,
	HTTP_VERSION_NOT_SUPPORTED:      HTTP_VERSION_NOT_SUPPORTED_MSG,
	VARIANT_ALSO_NEGOTIATES:         VARIANT_ALSO_NEGOTIATES_MSG,
	INSUFFICIENT_STORAGE:            INSUFFICIENT_STORAGE_MSG,
	LOOP_DETECTED:                   LOOP_DETECTED_MSG,
	NOT_EXTENDED:                    NOT_EXTENDED_MSG,
	NETWORK_AUTHENTICATION_REQUIRED: NETWORK_AUTHENTICATION_REQUIRED_MSG,
}

// CodeString devuelve el mensaje legible asociado a un código HTTP.
// Si el código no está registrado, devuelve un string vacío.
func CodeString(statusCode int) string {
	return StatusMap[statusCode]
}
