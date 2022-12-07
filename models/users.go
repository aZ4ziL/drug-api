package models

import (
	"database/sql"
	"time"

	"github.com/aZ4ziL/drug-api/auth"
)

type User struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	FirstName    string        `gorm:"size:30" json:"first_name"`
	LastName     string        `gorm:"size:30" json:"last_name"`
	Username     string        `gorm:"size:30;index;unique" json:"username"`
	Password     string        `gorm:"size:100" json:"-"`
	IsAdmin      bool          `gorm:"default:0" json:"is_admin"`
	LastLogin    sql.NullTime  `gorm:"null" json:"last_login"`
	DateJoined   time.Time     `gorm:"autoCreateTime" json:"date_joined"`
	Transactions []Transaction `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"transactions"`
}

type userModel struct{}

func NewUserModel() userModel {
	return userModel{}
}

// Create new user
func (u userModel) CreateNewUser(user *User) error {
	hashed, err := auth.EncryptionPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashed
	return db.Create(user).Error
}

// GetUserByID
// getting user by passing the `ID`
func (u userModel) GetUserByID(id uint) (User, error) {
	var user User
	err := db.Model(&User{}).Where("id = ?", id).Preload("Transactions").First(&user).Error
	return user, err
}

// GetUserByUsername
// getting user by passing the `Username`
func (u userModel) GetUserByUsername(username string) (User, error) {
	var user User
	err := db.Model(&User{}).Where("username = ?", username).Preload("Transactions").First(&user).Error
	return user, err
}
