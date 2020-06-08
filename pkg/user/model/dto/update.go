package dto

type UserUpdateInput struct {
	ID      string                 `json:"id"`
	Version string                 `json:"version"`
	Values  *UserUpdateValuesInput `json:"values"`
}

type UserUpdateValuesInput struct {
	Password *UserPasswordInput `json:"password"`
}
