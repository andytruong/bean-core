package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/infra/api"
	"bean/pkg/space/model"
)

// Query
type (
	SpaceQuery           struct{}
	SpaceMembershipQuery struct{}

	SpaceFilters struct {
		ID     *string     `json:"id"`
		Domain *scalar.Uri `json:"domain"`
	}

	MembershipsFilter struct {
		UserID    string                  `json:"userId"`
		Space     *MembershipsFilterSpace `json:"space"`
		IsActive  bool                    `json:"isActive"`
		ManagerId *string                 `json:"managerId"`
	}

	MembershipsFilterSpace struct {
		Title      *string `json:"title"`
		DomainName *string `json:"domainName"`
	}
)

// Mutation
type (
	SpaceMutation           struct{}
	SpaceMembershipMutation struct{}

	SpaceCreateOutcome struct {
		Errors []util.Error `json:"errors"`
		Space  *model.Space `json:"space"`
	}

	// space.create
	SpaceCreateInput struct {
		Object SpaceCreateInputObject `json:"object"`
	}

	SpaceCreateInputObject struct {
		Kind        model.SpaceKind    `json:"kind"`
		Title       *string            `json:"title"`
		Language    api.Language       `json:"language"`
		IsActive    bool               `json:"isActive"`
		DomainNames *DomainNamesInput  `json:"domainNames"`
		Features    SpaceFeaturesInput `json:"features"`

		// Internal field
		ParentId *string `json:"parentId"`
	}

	DomainNameInput struct {
		Verified *bool   `json:"verified"`
		Value    *string `json:"value"`
		IsActive *bool   `json:"isActive"`
	}

	DomainNamesInput struct {
		Primary   *DomainNameInput   `json:"primary"`
		Secondary []*DomainNameInput `json:"secondary"`
	}

	SpaceFeaturesInput struct {
		Register bool `json:"register"`
	}

	// space.update
	SpaceUpdateInput struct {
		SpaceID      string                  `json:"spaceId"`
		SpaceVersion string                  `json:"spaceVersion"`
		Object       *SpaceUpdateInputObject `json:"object"`
	}

	SpaceUpdateInputFeatures struct {
		Register *bool `json:"register"`
	}

	SpaceUpdateInputObject struct {
		Features *SpaceUpdateInputFeatures `json:"features"`
		Language *api.Language             `json:"language"`
	}

	// membership.create
	SpaceMembershipCreateInput struct {
		SpaceID          string   `json:"spaceId"`
		UserID           string   `json:"userId"`
		IsActive         bool     `json:"isActive"`
		ManagerMemberIds []string `json:"managerMemberIds"`
	}

	SpaceMembershipCreateOutcome struct {
		Errors     []*util.Error     `json:"errors"`
		Membership *model.Membership `json:"membership"`
	}

	// membership.update
	SpaceMembershipUpdateInput struct {
		Id       string        `json:"id"`
		Version  string        `json:"version"`
		IsActive bool          `json:"isActive"`
		Language *api.Language `json:"language"`
	}

	SpaceMembershipUpdateOutcome struct {
		Errors     []*util.Error     `json:"errors"`
		Membership *model.Membership `json:"membership"`
	}
)
