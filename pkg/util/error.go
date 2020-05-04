package util

type Error struct {
	Code   *ErrorCode `json:"code"`
	Fields []string   `json:"fields"`
}
