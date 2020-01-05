package admin

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	gorm.Model
	UserName   string `db:"user_name"`
	HashedPass string `db:"hashed_pass"`
}

type AdminManager struct {
	db *gorm.DB
}

func New(db *gorm.DB) *AdminManager {
	db.AutoMigrate(&Admin{})
	adminManager := &AdminManager{
		db: db,
	}

	// TODO: pubsub

	return adminManager
}

func (m *AdminManager) IsAdminExists() bool {
	var count int
	m.db.Table("admins").Count(&count)
	return count > 0
}

func (m *AdminManager) CreateAdmin(username, password string) error {
	hashedPass,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return m.db.Create(&Admin{
		UserName: username,
		HashedPass: string(hashedPass),
	}).Error
}

func (m *AdminManager) Authorize(username, password string) (*Admin, error) {
	adminUser := &Admin{}
	if err := m.db.Find(adminUser, "user_name = ?", username).Error; err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.HashedPass), []byte(password)); err != nil {
		return nil, err
	}

	return adminUser, nil
}



