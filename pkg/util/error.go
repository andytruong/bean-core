package util

type (
	Error struct {
		Code   *ErrorCode `json:"code"`
		Fields []string   `json:"fields"`
	}

	Err string
)

func (e Err) Error() string { return string(e) }

const NilPointerError = Err("nil pointer error")
