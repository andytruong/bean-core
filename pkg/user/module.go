package user

func NewUserService() *UserModule {
	return &UserModule{}
}

type UserModule struct {
	mutationResolver *UserMutationResolver
}

func (this *UserModule) MutationResolver() *UserMutationResolver {
	if nil == this.mutationResolver {
		this.mutationResolver = &UserMutationResolver{
		}
	}

	return this.mutationResolver
}
