package dto

import "bean/pkg/integration/s3/model"

type S3ApplicationPolicyMutationInput struct {
	Create []S3ApplicationPolicyCreateInput `json:"create"`
	Update []S3ApplicationPolicyUpdateInput `json:"update"`
	Delete []S3ApplicationPolicyDeleteInput `json:"delete"`
}

type S3ApplicationPolicyCreateInput struct {
	Kind  model.PolicyKind `json:"kind"`
	Value string           `json:"value"`
}

type S3ApplicationPolicyUpdateInput struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type S3ApplicationPolicyDeleteInput struct {
	Id string `json:"id"`
}
