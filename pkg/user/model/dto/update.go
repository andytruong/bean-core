package dto

type UserUpdateInput struct {
	ID      string                 `json:"id"`
	VersioN string                 `json:"versioN"`
	Values  *UserUpdateValuesInput `json:"values"`
}

type UserUpdateValuesInput struct {
	Password *UserPasswordInput `json:"password"`
}
