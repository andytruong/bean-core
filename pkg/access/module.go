package access

func NewAccessModule() *AccessModule {
	return &AccessModule{}
}

type AccessModule struct {
}

func (this AccessModule) MutationResolver() *AccessMutationResolver {
	return &AccessMutationResolver{}
}
