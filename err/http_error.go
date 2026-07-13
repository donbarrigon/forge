package err

import (
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/donbarrigon/forge/config"
)

type HttpError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
	Name    string `json:"name,omitempty"`
	Cause   string `json:"cause,omitempty"`
	Stack   string `json:"stack,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func New(status int, message string, e error) *HttpError {
	if config.Env.App.Debug {
		return &HttpError{
			Status:  status,
			Message: message,
			Name:    "Error",
			Cause:   causeData(e),
			Stack:   string(debug.Stack()),
			Data:    errorData(e),
		}
	} else {
		return &HttpError{
			Status:  status,
			Message: message,
			Name:    "Error",
			Cause:   "",
			Stack:   "",
			Data:    nil,
		}
	}
}

// ================================
// Funciones para la interfaz de error
// ================================

func (self *HttpError) Error() string {
	return fmt.Sprintf("[%d:%s] %s", self.Status, self.Name, self.Message)
}

// ================================
// 4xx Errores del cliente
// ================================

func BadRequest(e error) *HttpError {
	return New(BAD_REQUEST, BAD_REQUEST_MSG, e)
}

func BadRequestMsg(msg string, e error) *HttpError {
	return New(BAD_REQUEST, msg, e)
}

func Unauthorized(e error) *HttpError {
	return New(UNAUTHORIZED, UNAUTHORIZED_MSG, e)
}

func UnauthorizedMsg(msg string, e error) *HttpError {
	return New(UNAUTHORIZED, msg, e)
}

func PaymentRequired(e error) *HttpError {
	return New(PAYMENT_REQUIRED, PAYMENT_REQUIRED_MSG, e)
}

func PaymentRequiredMsg(msg string, e error) *HttpError {
	return New(PAYMENT_REQUIRED, msg, e)
}

func Forbidden(e error) *HttpError {
	return New(FORBIDDEN, FORBIDDEN_MSG, e)
}

func ForbiddenMsg(msg string, e error) *HttpError {
	return New(FORBIDDEN, msg, e)
}

func NotFound(e error) *HttpError {
	return New(NOT_FOUND, NOT_FOUND_MSG, e)
}

func NotFoundMsg(msg string, e error) *HttpError {
	return New(NOT_FOUND, msg, e)
}

func MethodNotAllowed(e error) *HttpError {
	return New(METHOD_NOT_ALLOWED, METHOD_NOT_ALLOWED_MSG, e)
}

func MethodNotAllowedMsg(msg string, e error) *HttpError {
	return New(METHOD_NOT_ALLOWED, msg, e)
}

func NotAcceptable(e error) *HttpError {
	return New(NOT_ACCEPTABLE, NOT_ACCEPTABLE_MSG, e)
}

func NotAcceptableMsg(msg string, e error) *HttpError {
	return New(NOT_ACCEPTABLE, msg, e)
}

func ProxyAuthRequired(e error) *HttpError {
	return New(PROXY_AUTHENTICATION_REQUIRED, PROXY_AUTHENTICATION_REQUIRED_MSG, e)
}

func ProxyAuthRequiredMsg(msg string, e error) *HttpError {
	return New(PROXY_AUTHENTICATION_REQUIRED, msg, e)
}

func RequestTimeout(e error) *HttpError {
	return New(REQUEST_TIMEOUT, REQUEST_TIMEOUT_MSG, e)
}

func RequestTimeoutMsg(msg string, e error) *HttpError {
	return New(REQUEST_TIMEOUT, msg, e)
}

func Conflict(e error) *HttpError {
	return New(CONFLICT, CONFLICT_MSG, e)
}

func ConflictMsg(msg string, e error) *HttpError {
	return New(CONFLICT, msg, e)
}

func Gone(e error) *HttpError {
	return New(GONE, GONE_MSG, e)
}

func GoneMsg(msg string, e error) *HttpError {
	return New(GONE, msg, e)
}

func LengthRequired(e error) *HttpError {
	return New(LENGTH_REQUIRED, LENGTH_REQUIRED_MSG, e)
}

func LengthRequiredMsg(msg string, e error) *HttpError {
	return New(LENGTH_REQUIRED, msg, e)
}

func PreconditionFailed(e error) *HttpError {
	return New(PRECONDITION_FAILED, PRECONDITION_FAILED_MSG, e)
}

func PreconditionFailedMsg(msg string, e error) *HttpError {
	return New(PRECONDITION_FAILED, msg, e)
}

func RequestEntityTooLarge(e error) *HttpError {
	return New(PAYLOAD_TOO_LARGE, PAYLOAD_TOO_LARGE_MSG, e)
}

func RequestEntityTooLargeMsg(msg string, e error) *HttpError {
	return New(PAYLOAD_TOO_LARGE, msg, e)
}

func RequestURITooLong(e error) *HttpError {
	return New(URI_TOO_LONG, URI_TOO_LONG_MSG, e)
}

func RequestURITooLongMsg(msg string, e error) *HttpError {
	return New(URI_TOO_LONG, msg, e)
}

func UnsupportedMediaType(e error) *HttpError {
	return New(UNSUPPORTED_MEDIA_TYPE, UNSUPPORTED_MEDIA_TYPE_MSG, e)
}

func UnsupportedMediaTypeMsg(msg string, e error) *HttpError {
	return New(UNSUPPORTED_MEDIA_TYPE, msg, e)
}

func RequestedRangeNotSatisfiable(e error) *HttpError {
	return New(RANGE_NOT_SATISFIABLE, RANGE_NOT_SATISFIABLE_MSG, e)
}

func RequestedRangeNotSatisfiableMsg(msg string, e error) *HttpError {
	return New(RANGE_NOT_SATISFIABLE, msg, e)
}

func ExpectationFailed(e error) *HttpError {
	return New(EXPECTATION_FAILED, EXPECTATION_FAILED_MSG, e)
}

func ExpectationFailedMsg(msg string, e error) *HttpError {
	return New(EXPECTATION_FAILED, msg, e)
}

func ImATeapot(e error) *HttpError {
	return New(IM_A_TEAPOT, IM_A_TEAPOT_MSG, e)
}

func ImATeapotMsg(msg string, e error) *HttpError {
	return New(IM_A_TEAPOT, msg, e)
}

func MisdirectedRequest(e error) *HttpError {
	return New(MISDIRECTED_REQUEST, MISDIRECTED_REQUEST_MSG, e)
}

func MisdirectedRequestMsg(msg string, e error) *HttpError {
	return New(MISDIRECTED_REQUEST, msg, e)
}

func UnprocessableEntity(e error) *HttpError {
	return New(UNPROCESSABLE_ENTITY, UNPROCESSABLE_ENTITY_MSG, e)
}

func UnprocessableEntityMsg(msg string, e error) *HttpError {
	return New(UNPROCESSABLE_ENTITY, msg, e)
}

func Locked(e error) *HttpError {
	return New(LOCKED, LOCKED_MSG, e)
}

func LockedMsg(msg string, e error) *HttpError {
	return New(LOCKED, msg, e)
}

func FailedDependency(e error) *HttpError {
	return New(FAILED_DEPENDENCY, FAILED_DEPENDENCY_MSG, e)
}

func FailedDependencyMsg(msg string, e error) *HttpError {
	return New(FAILED_DEPENDENCY, msg, e)
}

func TooEarly(e error) *HttpError {
	return New(TOO_EARLY, TOO_EARLY_MSG, e)
}

func TooEarlyMsg(msg string, e error) *HttpError {
	return New(TOO_EARLY, msg, e)
}

func UpgradeRequired(e error) *HttpError {
	return New(UPGRADE_REQUIRED, UPGRADE_REQUIRED_MSG, e)
}

func UpgradeRequiredMsg(msg string, e error) *HttpError {
	return New(UPGRADE_REQUIRED, msg, e)
}

func PreconditionRequired(e error) *HttpError {
	return New(PRECONDITION_REQUIRED, PRECONDITION_REQUIRED_MSG, e)
}

func PreconditionRequiredMsg(msg string, e error) *HttpError {
	return New(PRECONDITION_REQUIRED, msg, e)
}

func TooManyRequests(e error) *HttpError {
	return New(TOO_MANY_REQUESTS, TOO_MANY_REQUESTS_MSG, e)
}

func TooManyRequestsMsg(msg string, e error) *HttpError {
	return New(TOO_MANY_REQUESTS, msg, e)
}

func RequestHeaderFieldsTooLarge(e error) *HttpError {
	return New(REQUEST_HEADER_FIELDS_TOO_LARGE, REQUEST_HEADER_FIELDS_TOO_LARGE_MSG, e)
}

func RequestHeaderFieldsTooLargeMsg(msg string, e error) *HttpError {
	return New(REQUEST_HEADER_FIELDS_TOO_LARGE, msg, e)
}

func UnavailableForLegalReasons(e error) *HttpError {
	return New(UNAVAILABLE_FOR_LEGAL_REASONS, UNAVAILABLE_FOR_LEGAL_REASONS_MSG, e)
}

func UnavailableForLegalReasonsMsg(msg string, e error) *HttpError {
	return New(UNAVAILABLE_FOR_LEGAL_REASONS, msg, e)
}

// ================================
// 5xx Server errors
// ================================

func Internal(e error) *HttpError {
	return New(INTERNAL, INTERNAL_MSG, e)
}

func InternalMsg(msg string, e error) *HttpError {
	return New(INTERNAL, msg, e)
}

// Panic es un caso especial: recibe el stack ya capturado en el punto del
// recover(), en vez de dejar que New() capture uno nuevo (que solo mostraría
// la pila dentro de New(), inútil para depurar el panic real).
func Panic(rec any, stack string) *HttpError {
	return PanicMsg("Oops, something went very wrong", rec, stack)
}

func PanicMsg(msg string, rec any, stack string) *HttpError {
	if config.Env.App.Debug {
		return &HttpError{
			Status:  INTERNAL,
			Message: msg,
			Name:    "Panic",
			Cause:   "panic",
			Stack:   stack,
			Data:    errorData(rec),
		}
	}
	return &HttpError{
		Status:  INTERNAL,
		Message: msg,
		Name:    "Panic",
	}
}

func NotImplemented(e error) *HttpError {
	return New(NOT_IMPLEMENTED, NOT_IMPLEMENTED_MSG, e)
}

func NotImplementedMsg(msg string, e error) *HttpError {
	return New(NOT_IMPLEMENTED, msg, e)
}

func BadGateway(e error) *HttpError {
	return New(BAD_GATEWAY, BAD_GATEWAY_MSG, e)
}

func BadGatewayMsg(msg string, e error) *HttpError {
	return New(BAD_GATEWAY, msg, e)
}

func ServiceUnavailable(e error) *HttpError {
	return New(SERVICE_UNAVAILABLE, SERVICE_UNAVAILABLE_MSG, e)
}

func ServiceUnavailableMsg(msg string, e error) *HttpError {
	return New(SERVICE_UNAVAILABLE, msg, e)
}

func GatewayTimeout(e error) *HttpError {
	return New(GATEWAY_TIMEOUT, GATEWAY_TIMEOUT_MSG, e)
}

func GatewayTimeoutMsg(msg string, e error) *HttpError {
	return New(GATEWAY_TIMEOUT, msg, e)
}

func HTTPVersionNotSupported(e error) *HttpError {
	return New(HTTP_VERSION_NOT_SUPPORTED, HTTP_VERSION_NOT_SUPPORTED_MSG, e)
}

func HTTPVersionNotSupportedMsg(msg string, e error) *HttpError {
	return New(HTTP_VERSION_NOT_SUPPORTED, msg, e)
}

func VariantAlsoNegotiates(e error) *HttpError {
	return New(VARIANT_ALSO_NEGOTIATES, VARIANT_ALSO_NEGOTIATES_MSG, e)
}

func VariantAlsoNegotiatesMsg(msg string, e error) *HttpError {
	return New(VARIANT_ALSO_NEGOTIATES, msg, e)
}

func InsufficientStorage(e error) *HttpError {
	return New(INSUFFICIENT_STORAGE, INSUFFICIENT_STORAGE_MSG, e)
}

func InsufficientStorageMsg(msg string, e error) *HttpError {
	return New(INSUFFICIENT_STORAGE, msg, e)
}

func LoopDetected(e error) *HttpError {
	return New(LOOP_DETECTED, LOOP_DETECTED_MSG, e)
}

func LoopDetectedMsg(msg string, e error) *HttpError {
	return New(LOOP_DETECTED, msg, e)
}

func NotExtended(e error) *HttpError {
	return New(NOT_EXTENDED, NOT_EXTENDED_MSG, e)
}

func NotExtendedMsg(msg string, e error) *HttpError {
	return New(NOT_EXTENDED, msg, e)
}

func NetworkAuthenticationRequired(e error) *HttpError {
	return New(NETWORK_AUTHENTICATION_REQUIRED, NETWORK_AUTHENTICATION_REQUIRED_MSG, e)
}

func NetworkAuthenticationRequiredMsg(msg string, e error) *HttpError {
	return New(NETWORK_AUTHENTICATION_REQUIRED, msg, e)
}

// ================================
// Helpers
// ================================

func errorData(e any) any {
	if e == nil {
		return nil
	}

	v := reflect.ValueOf(e)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		result := make(map[string]any)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)

			if field.IsExported() {
				fieldValue := v.Field(i)
				result[field.Name] = fieldValue.Interface()
			}
		}

		if len(result) > 0 {
			return result
		}
	}

	if er, ok := e.(error); ok {
		return er.Error()
	}

	return e
}

func causeData(e error) string {
	if e == nil {
		return ""
	}
	return e.Error()
}
