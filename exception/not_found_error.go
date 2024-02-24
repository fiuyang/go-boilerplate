package exception

// type NotFoundError struct {
//     message string
// }

// func (e *NotFoundError) Error() string {
//     return e.message
// }

// func NewNotFoundError(message string) *NotFoundError {
//     return &NotFoundError{message: message}
// }
type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}