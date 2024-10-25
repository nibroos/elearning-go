package dtos

import "github.com/nibroos/elearning-go/service/internal/utils"

type GetUsersRequest struct {
	Global         string                 `json:"global"`
	Username       string                 `json:"username"`
	Name           string                 `json:"name"`
	Email          string                 `json:"email"`
	PerPage        utils.Nullable[string] `json:"per_page" default:"10"`         // Default per_page to 10
	Page           utils.Nullable[string] `json:"page" default:"1"`              // Default page to 1
	OrderColumn    string                 `json:"order_column" default:"id"`     // Default order column to "id"
	OrderDirection string                 `json:"order_direction" default:"asc"` // Default order direction to "asc"
}

type CreateUserRequest struct {
	Name     string                 `json:"name"`
	Username utils.Nullable[string] `json:"username"`
	Email    string                 `json:"email"`
	Address  utils.Nullable[string] `json:"address"`
	Password string                 `json:"password"`
	RoleIDs  []uint32               `json:"role_ids"`
}

type UpdateUserRequest struct {
	ID       uint                   `json:"id"`
	Username utils.Nullable[string] `json:"username"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Address  utils.Nullable[string] `json:"address"`
	Password utils.Nullable[string] `json:"password"`
	RoleIDs  []uint32               `json:"role_ids"`
}

type GetUserByIDRequest struct {
	ID uint `json:"id"`
}

type DeleteUserRequest struct {
	ID uint `json:"id"`
}

type UserListDTO struct {
	ID       int                    `json:"id"`
	Username utils.Nullable[string] `json:"username"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
}

type UserDetailDTO struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Username    utils.Nullable[string] `json:"username"`
	Email       string                 `json:"email"`
	Address     utils.Nullable[string] `json:"address"`
	Password    utils.Nullable[string] `json:"password"`
	Roles       []string               `json:"roles"`
	Permissions []string               `json:"permissions"`
}
type GetUsersResult struct {
	Users []UserListDTO
	Total int
	Err   error
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Name     string                 `json:"name"`
	Username utils.Nullable[string] `json:"username"`
	Email    string                 `json:"email"`
	Address  utils.Nullable[string] `json:"address"`
	Password string                 `json:"password"`
}

type GetSubscribesRequest struct {
	Global         string                 `json:"global"`
	Name           string                 `json:"name"`
	PerPage        utils.Nullable[string] `json:"per_page" default:"10"`         // Default per_page to 10
	Page           utils.Nullable[string] `json:"page" default:"1"`              // Default page to 1
	OrderColumn    string                 `json:"order_column" default:"id"`     // Default order column to "id"
	OrderDirection string                 `json:"order_direction" default:"asc"` // Default order direction to "asc"
}

type CreateSubscribeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateSubscribeRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetSubscribeByIDRequest struct {
	ID uint `json:"id"`
}

type DeleteSubscribeRequest struct {
	ID uint `json:"id"`
}

type SubscribeListDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SubscribeDetailDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type GetSubscribesResult struct {
	Subscribes []SubscribeListDTO
	Total      int
	Err        error
}

type ClassesRequest struct {
	Global         string                 `json:"global"`
	Name           string                 `json:"name"`
	SubcribeID     utils.Nullable[string] `json:"subcribe_id"`
	InchargeID     utils.Nullable[string] `json:"incharge_id"`
	PerPage        utils.Nullable[string] `json:"per_page" default:"10"`         // Default per_page to 10
	Page           utils.Nullable[string] `json:"page" default:"1"`              // Default page to 1
	OrderColumn    string                 `json:"order_column" default:"id"`     // Default order column to "id"
	OrderDirection string                 `json:"order_direction" default:"asc"` // Default order direction to "asc"
}

type CreateClassRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SubcribeID  uint   `json:"subcribe_id"`
	InchargeID  uint   `json:"incharge_id"`
}

type UpdateClassRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SubcribeID  uint   `json:"subcribe_id"`
	InchargeID  uint   `json:"incharge_id"`
}

type GetClassByIDRequest struct {
	ID uint `json:"id"`
}

type DeleteClassRequest struct {
	ID uint `json:"id"`
}

type ClassListDTO struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	BannerURL     string `json:"banner_url"`
	LogoURL       string `json:"logo_url"`
	VideoURL      string `json:"video_url"`
	CreatedByName string `json:"created_by_name"`
	UpdatedByName string `json:"updated_by_name"`
	InchargeName  string `json:"incharge_name"`
	SubcribeName  string `json:"subcribe_name"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type ClassDetailDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BannerURL   string `json:"banner_url"`
	LogoURL     string `json:"logo_url"`
	VideoURL    string `json:"video_url"`
	CreatedByID uint   `json:"created_by_id"`
	UpdatedByID uint   `json:"updated_by_id"`
	InchargeID  uint   `json:"incharge_id"`
	SubcribeID  uint   `json:"subcribe_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
type GetClassesResult struct {
	Classes []ClassListDTO
	Total   int
	Err     error
}
