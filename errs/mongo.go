package errs

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func HexID(e error) *Error {
	return New(BAD_REQUEST, "The identifier isn't valid", e)
}

// Convierte los errores del driver de mongo a errores HTTP
func Mongo(e error) *Error {
	if e == nil {
		return nil
	}

	// No document found
	if e == mongo.ErrNoDocuments {
		return New(NOT_FOUND, NOT_FOUND_MSG, e)
	}

	// Client disconnected
	if e == mongo.ErrClientDisconnected {
		return New(SERVICE_UNAVAILABLE, "We lost the connection to the database", e)
	}

	// Deadline exceeded (timeout)
	if errors.Is(e, context.DeadlineExceeded) {
		return New(REQUEST_TIMEOUT, "The operation took too long", e)
	}

	// Context canceled
	if errors.Is(e, context.Canceled) {
		return New(BAD_REQUEST, "The operation was canceled", e)
	}

	// Handle WriteException for more detailed write errors
	var writeException mongo.WriteException
	if errors.As(e, &writeException) {
		return handleWriteException(writeException)
	}

	// Handle BulkWriteException
	var bulkWriteException mongo.BulkWriteException
	if errors.As(e, &bulkWriteException) {
		return handleBulkWriteException(bulkWriteException)
	}

	// Handle CommandError
	var commandError mongo.CommandError
	if errors.As(e, &commandError) {
		return handleCommandError(commandError)
	}

	// Handle ServerError
	var serverError mongo.ServerError
	if errors.As(e, &serverError) {
		return handleServerError(serverError)
	}

	// Duplicate key (legacy check as fallback)
	if mongo.IsDuplicateKeyError(e) {
		return New(CONFLICT, "This record already exists", e)
	}

	// Check for network/connection errors by error message
	errorMsg := strings.ToLower(e.Error())
	switch {
	case strings.Contains(errorMsg, "connection refused"),
		strings.Contains(errorMsg, "no reachable servers"),
		strings.Contains(errorMsg, "server selection timeout"):
		return New(SERVICE_UNAVAILABLE, "We couldn't connect to the database", e)

	case strings.Contains(errorMsg, "authentication failed"):
		return New(UNAUTHORIZED, "Database authentication failed", e)

	case strings.Contains(errorMsg, "not authorized"):
		return New(FORBIDDEN, "You don't have sufficient permissions", e)

	case strings.Contains(errorMsg, "invalid namespace"):
		return New(BAD_REQUEST, "The collection name isn't valid", e)

	case strings.Contains(errorMsg, "exceeds maximum"):
		return New(BAD_REQUEST, "The data is too large", e)
	}

	// Default case → Internal Server Error
	return Internal(e)
}

// handleWriteException processes MongoDB write exceptions
func handleWriteException(we mongo.WriteException) *Error {
	for _, writeError := range we.WriteErrors {
		switch writeError.Code {
		case 11000, 11001: // Duplicate key errors
			return New(CONFLICT, "This record already exists", we)
		case 2: // BadValue
			return New(BAD_REQUEST, "One of the field values isn't valid", we)
		case 9: // FailedToParse
			return New(BAD_REQUEST, "We couldn't process the data", we)
		case 14: // TypeMismatch
			return New(BAD_REQUEST, "The data type doesn't match", we)
		case 16755: // Location error
			return New(BAD_REQUEST, "The geographic location isn't valid", we)
		case 17280: // KeyTooLong
			return New(BAD_REQUEST, "The identifier is too long", we)
		case 10334: // BSONObjectTooLarge
			return New(BAD_REQUEST, "The document is too large", we)
		}
	}

	// Check write concern errors
	if we.WriteConcernError != nil {
		switch we.WriteConcernError.Code {
		case 64: // WriteConcernFailed
			return New(SERVICE_UNAVAILABLE, "We couldn't confirm the write", we)
		case 79: // UnknownReplWriteConcern
			return New(BAD_REQUEST, "The write configuration isn't valid", we)
		}
	}

	return New(INTERNAL, "We couldn't save the data", we)
}

// handleBulkWriteException processes bulk write exceptions
func handleBulkWriteException(bwe mongo.BulkWriteException) *Error {
	// Check for duplicate key errors in bulk operations
	for _, writeError := range bwe.WriteErrors {
		if writeError.Code == 11000 || writeError.Code == 11001 {
			return New(CONFLICT, "Some records already exist", bwe)
		}
	}

	// Check write concern errors
	if bwe.WriteConcernError != nil {
		return New(SERVICE_UNAVAILABLE, "We couldn't confirm all the changes", bwe)
	}

	return New(BAD_REQUEST, "We couldn't save all the records", bwe)
}

// handleCommandError processes MongoDB command errors
func handleCommandError(ce mongo.CommandError) *Error {
	switch ce.Code {
	case 2: // BadValue
		return New(BAD_REQUEST, "One of the parameters isn't valid", ce)
	case 9: // FailedToParse
		return New(BAD_REQUEST, "We couldn't understand the request", ce)
	case 13: // Unauthorized
		return New(FORBIDDEN, FORBIDDEN_MSG, ce)
	case 18: // AuthenticationFailed
		return New(UNAUTHORIZED, "Authentication failed", ce)
	case 26: // NamespaceNotFound
		return New(NOT_FOUND, "The collection doesn't exist", ce)
	case 59: // CommandNotFound
		return New(BAD_REQUEST, "The operation isn't valid", ce)
	case 61: // ShardKeyNotFound
		return New(BAD_REQUEST, "The shard key is missing", ce)
	case 72: // InvalidOptions
		return New(BAD_REQUEST, "The options aren't valid", ce)
	case 96: // OperationFailed
		return New(INTERNAL, "The operation failed", ce)
	case 11600: // InterruptedAtShutdown
		return New(SERVICE_UNAVAILABLE, "The server is restarting", ce)
	case 11601: // Interrupted
		return New(SERVICE_UNAVAILABLE, "The operation was interrupted", ce)
	case 13435: // ShardKeyTooBig
		return New(BAD_REQUEST, "The shard key is too large", ce)
	case 16550: // DocumentValidationFailure
		return New(BAD_REQUEST, "The data doesn't meet the validation rules", ce)
	case 50: // MaxTimeMSExpired
		return New(REQUEST_TIMEOUT, "The operation took too long", ce)
	}

	return New(INTERNAL, "The operation failed", ce)
}

// handleServerError processes general MongoDB server errors
func handleServerError(se mongo.ServerError) *Error {
	// Server errors are typically infrastructure issues
	return New(SERVICE_UNAVAILABLE, "There's a problem with the database server", se)
}

func MongoUpdateResult(result *mongo.UpdateResult) *Error {
	if result.MatchedCount == 0 {
		return New(NOT_FOUND, "The document to update doesn't exist", errors.New("!result.MatchedCount == 0"))
	}
	if result.ModifiedCount == 0 {
		return New(CONFLICT, "No changes were applied when saving the document", errors.New("!result.ModifiedCount == 0"))
	}
	return nil
}

func MongoDeleteResult(result *mongo.DeleteResult) *Error {
	if result.DeletedCount == 0 {
		return New(CONFLICT, "The document wasn't deleted", errors.New("!result.DeletedCount == 0"))
	}
	return nil
}
