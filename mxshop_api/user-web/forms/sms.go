package forms

type SendSmsForm struct {
	Mobile string `json:"mobile" form:"mobile" binding:"required,mobile"`
	Type   uint   `json:"type" form:"type" binding:"required,oneof=1 2"`
}
