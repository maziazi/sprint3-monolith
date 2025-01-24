package model

type User struct {
	Id        uint   `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"-"`
	CreatedAt string `json:"createdAt"`
}
type UserProfile struct {
	Id     uint   `json:"id"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	FileId uint   `json:"fileId"`
}
