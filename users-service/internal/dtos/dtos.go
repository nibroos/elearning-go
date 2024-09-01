package dtos

type GetUsersRequest struct {
	Global         string `json:"global"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	PerPage        string `json:"per_page" default:"10"`         // Default per_page to 10
	Page           string `json:"page" default:"1"`              // Default page to 1
	OrderColumn    string `json:"order_column" default:"id"`     // Default order column to "id"
	OrderDirection string `json:"order_direction" default:"asc"` // Default order direction to "asc"
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Address  string `json:"address" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	RoleIDs  []uint `json:"role_ids" validate:"required"`
}

type GetUserByIDRequest struct {
	ID uint `json:"id" validate:"required"`
}

type UserListDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserDetailDTO struct {
	ID              uint     `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Address         string   `json:"address"`
	RoleIDs         []uint   `json:"role_ids"`
	PermissionIDs   []uint   `json:"permission_ids"`
	PermissionNames []string `json:"permission_names"`
}
