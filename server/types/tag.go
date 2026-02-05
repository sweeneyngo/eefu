package types

type TagInput struct {
	Name        string  `json:"name" validate:"required,min=1,max=32"`
	Type        string  `json:"type" validate:"required,oneof=archive custom stage status"`
	Description *string `json:"description" validate:"omitempty,min=1,max=256"`
}
