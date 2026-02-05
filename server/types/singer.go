package types

type SingerInput struct {
	Name    string       `json:"name" validate:"required,min=1,max=64"`
	Aliases []AliasInput `json:"aliases" validate:"omitempty,dive,min=1,max=64"`
}

type AliasInput struct {
	Name     string `json:"name" validate:"required,min=1,max=64"`
	Language string `json:"language" validate:"required,len=2"` // ISO 639-1 codes
}
