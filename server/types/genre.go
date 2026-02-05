package types

type GenreInput struct {
	Name string `json:"name" validate:"required,min=1,max=32"`
}
