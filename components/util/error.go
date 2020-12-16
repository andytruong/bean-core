package util

func NewErrors(code ErrorCode, fields []string, message string) []*Error {
	return []*Error{NewError(code, fields, message)}
}

func NewError(code ErrorCode, fields []string, message string) *Error {
	return &Error{
		Code:    &code,
		Fields:  fields,
		Message: message,
	}
}

type (
	Error struct {
		Code    *ErrorCode `json:"code"`
		Fields  []string   `json:"fields"`
		Message string     `json:"message"`
	}

	// TODO: Review https://blog.golang.org/go1.13-errors
	Err string
)

func (e Err) Error() string { return string(e) }

const (
	ErrorInvalidArgument = Err("invalid argument")
	ErrorConfig          = Err("configuration error")
	ErrorNilPointer      = Err("nil pointer error")
	ErrorQuery           = Err("query error")
	ErrorVersionConflict = Err("version conflict")
	ErrorLocked          = Err("locked")
	ErrorAuthRequired    = Err("auth required")
	ErrorAccessDenied    = Err("access denied")
	ErrorUselessInput    = Err("useless input")
)

func NilPointerErrorValidate(values ...interface{}) error {
	for _, value := range values {
		if nil == value {
			return ErrorNilPointer
		}
	}

	return nil
}
