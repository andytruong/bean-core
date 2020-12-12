package dto

type MailerMutation struct {
	Version  string                  `json:"version"`
	Account  *MailerAccountMutation  `json:"account"`
	Template *MailerTemplateMutation `json:"template"`
}

type MailerQuery struct {
	Version string              `json:"version"`
	Account *MailerQueryAccount `json:"account"`
}

type MailerQueryAccount struct {
	Get         *MailerAccount   `json:"get"`
	GetMultiple []*MailerAccount `json:"getMultiple"`
}
