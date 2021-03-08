package customerrors

type ErrorLevel int

const (
	Business ErrorLevel = 0
	Server   ErrorLevel = 1

	DetailBusiness = "check the input parameters"
	DetailServer   = "something went wrong"
)

type Error struct {
	Err    error
	Level  ErrorLevel
	Detail string
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func New(err error) *Error {
	return &Error{
		Err:    err,
		Level:  Server,
		Detail: DetailServer,
	}
}

func NewBusiness(err error, detail string) *Error {
	businessErr := &Error{
		Err:   err,
		Level: Business,
	}

	if detail == "" {
		businessErr.Detail = DetailBusiness
	} else {
		businessErr.Detail = detail
	}

	return businessErr
}
