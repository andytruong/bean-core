package app

func (bundle *AppBundle) newResolvers() map[string]interface{} {
	return map[string]interface{}{
		"Query":    map[string]interface{}{},
		"Mutation": map[string]interface{}{},
	}
}
