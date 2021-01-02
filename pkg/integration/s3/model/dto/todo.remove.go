package dto

// TODO: remove
type S3ApplicationCreateInput struct {
	IsActive    bool                             `json:"isActive"`
	Credentials S3ApplicationCredentialsInput    `json:"credentials"`
	Policies    []S3ApplicationPolicyCreateInput `json:"policies"`
}

// TODO: Remove
type S3ApplicationUpdateInput struct {
	Id          string                               `json:"id"`
	Version     string                               `json:"version"`
	IsActive    *bool                                `json:"isActive"`
	Credentials *S3ApplicationCredentialsUpdateInput `json:"credentials"`
	Policies    *S3ApplicationPolicyMutationInput    `json:"policies"`
}
