package errors

type Error string
func (e Error) Error() string { return string(e) }


const (
	ErrCannotOpenDB = Error("cannot open db")
)