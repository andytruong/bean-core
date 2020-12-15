package dto

type MailerMutation struct {
	Account  *MailerAccountMutation  `json:"account"`
	Template *MailerTemplateMutation `json:"template"`
}

type MailerQuery struct {
	Account *MailerQueryAccount `json:"account"`
}

type MailerQueryAccount struct {
}
