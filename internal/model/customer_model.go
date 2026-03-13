package model

type CustomerResponse struct {
	CstID         int                  `json:"cstId"`
	NationalityID int                  `json:"nationalityId"`
	CstName       string               `json:"cstName"`
	CstDob        string               `json:"cstDob"`
	CstPhoneNum   string               `json:"cstPhoneNum"`
	CstEmail      string               `json:"cstEmail"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
	Nationality   *NationalityResponse `json:"nationality,omitempty"`
	FamilyLists   []FamilyListResponse `json:"family_lists,omitempty"`
}

type CreateCustomerRequest struct {
	NationalityID int    `json:"nationalityId"`
	CstName       string `json:"cstName" validate:"required,max=50"`
	CstDob        string `json:"cstDob" validate:"required,valid_date"`
	CstPhoneNum   string `json:"cstPhoneNum" validate:"required,max=20"`
	CstEmail      string `json:"cstEmail" validate:"required,email,max=50"`
}

type UpdateCustomerRequest struct {
	NationalityID int    `json:"nationalityId"`
	CstID         int    `json:"-" validate:"required,lte=100"`
	CstName       string `json:"cstName,omitempty" validate:"max=50"`
	CstDob        string `json:"cstDob,omitempty" validate:"valid_date"`
	CstPhoneNum   string `json:"cstPhoneNum,omitempty" validate:"max=20"`
	CstEmail      string `json:"cstEmail,omitempty" validate:"email,max=50"`
}

type SearchCustomerRequest struct {
	NationalityID int    `json:"-"`
	CstName       string `json:"cstName" validate:"max=100"`
	CstPhoneNum   string `json:"cstPhoneNum" validate:"max=20"`
	CstEmail      string `json:"cstEmail" validate:"max=200"`
	Page          int    `json:"page" validate:"min=1"`
	Size          int    `json:"size" validate:"min=1,max=100"`
}

type GetCustomerRequest struct {
	CstID int `json:"-" validate:"required,lte=100"`
}

type DeleteCustomerRequest struct {
	FamilyListId int `json:"-"`
	CstID        int `json:"-" validate:"required,lte=100"`
}
