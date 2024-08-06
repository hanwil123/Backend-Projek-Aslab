package Auth

type UserBaru struct {
	Id       uint64 `gorm:"primaryKey" json:"id"`
	Userid   uint64 `gorm:"foreignKey" json:"userid"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password []byte
}
