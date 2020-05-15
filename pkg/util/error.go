package util

type (
	Error struct {
		Code   *ErrorCode `json:"code"`
		Fields []string   `json:"fields"`
	}
)

// TODO: Review https://blog.golang.org/go1.13-errors
type Err string

func (e Err) Error() string { return string(e) }

const ErrorInvalidArgument = Err("invalid argument")
const ErrorConfig = Err("configuration error")
const ErrorNilPointer = Err("nil pointer error")
const ErrorQuery = Err("query error")

func NilPointerErrorValidate(values ...interface{}) error {
	for _, value := range values {
		if nil == value {
			return ErrorNilPointer
		}
	}

	return nil
}
