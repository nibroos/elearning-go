package dtos

import (
	"time"

	"github.com/nibroos/elearning-go/service/internal/utils"
)

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
	CreatedAt   *string                `json:"created_at"`
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
	ID            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	DeleteAt      *string `json:"deleted_at" db:"deleted_at"`
}

type SubscribeDetailDTO struct {
	ID            uint    `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	CreatedByID   uint    `json:"created_by_id" db:"created_by_id"`
	UpdatedByID   *uint   `json:"updated_by_id" db:"updated_by_id"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	DeletedAt     *string `json:"deleted_at" db:"deleted_at"`
}
type GetSubscribesResult struct {
	Subscribes []SubscribeListDTO
	Total      int
	Err        error
}

type ClassesRequest struct {
	Global         string                 `json:"global"`
	Name           string                 `json:"name"`
	SubscribeID    utils.Nullable[string] `json:"subscribe_id"`
	InchargeID     utils.Nullable[string] `json:"incharge_id"`
	PerPage        utils.Nullable[string] `json:"per_page" default:"10"`         // Default per_page to 10
	Page           utils.Nullable[string] `json:"page" default:"1"`              // Default page to 1
	OrderColumn    string                 `json:"order_column" default:"id"`     // Default order column to "id"
	OrderDirection string                 `json:"order_direction" default:"asc"` // Default order direction to "asc"
}

type CreateClassRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SubscribeID uint   `json:"subscribe_id"`
	InchargeID  uint   `json:"incharge_id"`
}

type UpdateClassRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SubscribeID uint   `json:"subscribe_id"`
	InchargeID  uint   `json:"incharge_id"`
}

type GetClassByIDRequest struct {
	ID uint `json:"id"`
}

type DeleteClassRequest struct {
	ID uint `json:"id"`
}

type ClassListDTO struct {
	ID            uint    `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   *string `json:"description" db:"description"`
	BannerURL     *string `json:"banner_url" db:"banner_url"`
	LogoURL       *string `json:"logo_url" db:"logo_url"`
	VideoURL      *string `json:"video_url" db:"video_url"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	InchargeName  *string `json:"incharge_name" db:"incharge_name"`
	SubscribeName *string `json:"subscribe_name" db:"subscribe_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	SubscribeID   uint    `json:"subscribe_id" db:"subscribe_id"`
	InchargeID    uint    `json:"incharge_id" db:"incharge_id"`
}

type ClassDetailDTO struct {
	ID            uint    `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   *string `json:"description" db:"description"`
	BannerURL     *string `json:"banner_url" db:"banner_url"`
	LogoURL       *string `json:"logo_url" db:"logo_url"`
	VideoURL      *string `json:"video_url" db:"video_url"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	InchargeName  *string `json:"incharge_name" db:"incharge_name"`
	SubscribeName *string `json:"subscribe_name" db:"subscribe_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	DeletedAt     *string `json:"deleted_at" db:"deleted_at"`
	CreatedByID   uint    `json:"created_by_id" db:"created_by_id"`
	UpdatedByID   *uint   `json:"updated_by_id" db:"updated_by_id"`
	SubscribeID   uint    `json:"subscribe_id" db:"subscribe_id"`
	InchargeID    uint    `json:"incharge_id" db:"incharge_id"`
}
type GetClassesResult struct {
	Classes []ClassListDTO
	Total   int
	Err     error
}

type GetModulesRequest struct {
	Global         string                 `json:"global"`
	Name           string                 `json:"name"`
	PerPage        utils.Nullable[string] `json:"per_page" default:"10"`         // Default per_page to 10
	Page           utils.Nullable[string] `json:"page" default:"1"`              // Default page to 1
	OrderColumn    string                 `json:"order_column" default:"id"`     // Default order column to "id"
	OrderDirection string                 `json:"order_direction" default:"asc"` // Default order direction to "asc"
}

type CreateModuleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ClassID     uint   `json:"class_id"`
}

type UpdateModuleRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ClassID     uint   `json:"class_id"`
}

type GetModuleByIDRequest struct {
	ID uint `json:"id"`
}

type DeleteModuleRequest struct {
	ID uint `json:"id"`
}

type ModuleListDTO struct {
	ID            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	ClassID       uint    `json:"class_id" db:"class_id"`
	ClassName     string  `json:"class_name" db:"class_name"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	DeleteAt      *string `json:"deleted_at" db:"deleted_at"`
}

type ModuleDetailDTO struct {
	ID            uint    `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	LogoURL       *string `json:"logo_url" db:"logo_url"`
	VideoURL      *string `json:"video_url" db:"video_url"`
	ClassID       uint    `json:"class_id" db:"class_id"`
	ClassName     string  `json:"class_name" db:"class_name"`
	CreatedByID   uint    `json:"created_by_id" db:"created_by_id"`
	UpdatedByID   *uint   `json:"updated_by_id" db:"updated_by_id"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	DeletedAt     *string `json:"deleted_at" db:"deleted_at"`
}
type GetModulesResult struct {
	Modules []ModuleListDTO
	Total   int
	Err     error
}
type CreateEducationRequest struct {
	ModuleID    uint   `json:"module_id"`
	NoUrut      uint   `json:"no_urut"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TextMateri  string `json:"text_materi"`
	// AttachmentUrls []string `json:"attachment_urls"`
}

type UpdateEducationRequest struct {
	ID          uint   `json:"id"`
	ModuleID    uint   `json:"module_id"`
	NoUrut      uint   `json:"no_urut"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TextMateri  string `json:"text_materi"`
	// AttachmentUrls []string `json:"attachment_urls"`
}

type GetEducationByIDRequest struct {
	ID uint `json:"id"`
}

type GetEducationParams struct {
	ID        uint
	IsDeleted *int
}

func NewGetEducationParams(id uint) *GetEducationParams {
	defaultIsDeleted := 0
	return &GetEducationParams{
		ID:        id,
		IsDeleted: &defaultIsDeleted,
	}
}

type DeleteEducationRequest struct {
	ID uint `json:"id"`
}

type EducationListDTO struct {
	ID            int     `json:"id" db:"id"`
	NoUrut        uint    `json:"no_urut" db:"no_urut"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	ThumbnailURL  *string `json:"thumbnail_url" db:"thumbnail_url"`
	VideoURL      *string `json:"video_url" db:"video_url"`
	ModuleID      uint    `json:"module_id" db:"module_id"`
	ModuleName    string  `json:"module_name" db:"module_name"`
	TextMaterial  string  `json:"text_materi" db:"text_materi"`
	CreatedByName *string `json:"created_by_name" db:"created_by_name"`
	UpdatedByName *string `json:"updated_by_name" db:"updated_by_name"`
	CreatedAt     *string `json:"created_at" db:"created_at"`
	UpdatedAt     *string `json:"updated_at" db:"updated_at"`
	DeleteAt      *string `json:"deleted_at" db:"deleted_at"`
}

type EducationDetailDTO struct {
	ID             uint                  `json:"id" db:"id"`
	NoUrut         uint                  `json:"no_urut" db:"no_urut"`
	Name           string                `json:"name" db:"name"`
	Description    string                `json:"description" db:"description"`
	ThumbnailURL   *string               `json:"thumbnail_url" db:"thumbnail_url"`
	VideoURL       *string               `json:"video_url" db:"video_url"`
	AttachmentURLs utils.JSONStringArray `json:"attachment_urls" db:"attachment_urls"`
	ModuleID       uint                  `json:"module_id" db:"module_id"`
	ModuleName     string                `json:"module_name" db:"module_name"`
	TextMaterial   string                `json:"text_materi" db:"text_materi"`
	CreatedByID    uint                  `json:"created_by_id" db:"created_by_id"`
	UpdatedByID    *uint                 `json:"updated_by_id" db:"updated_by_id"`
	CreatedByName  *string               `json:"created_by_name" db:"created_by_name"`
	UpdatedByName  *string               `json:"updated_by_name" db:"updated_by_name"`
	CreatedAt      *time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time            `json:"deleted_at" db:"deleted_at"`
}
type GetEducationsResult struct {
	Educations []EducationListDTO
	Total      int
	Err        error
}
