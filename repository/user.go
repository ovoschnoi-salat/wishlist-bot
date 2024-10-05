package repository

import "gorm.io/gorm"

type User struct {
	ID           int64
	Username     string
	Name         string
	Language     string
	Lists        []List `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Friends      []User `gorm:"many2many:friends;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reservations []Wish `gorm:"foreignKey:ReservedBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func AddUser(db *gorm.DB, userId int64, username, lang string) error {
	return db.Create(&User{ID: userId, Username: username, Language: lang}).Error
}

func GetUserById(db *gorm.DB, userId int64) (user User, err error) {
	err = db.First(&user, userId).Error
	return
}

func GetUserByUsername(db *gorm.DB, username string) (user User, err error) {
	err = db.First(&user, "username = ?", username).Error
	return
}

func UpdateUserLanguage(db *gorm.DB, userId int64, lang string) error {
	return db.Model(&User{ID: userId}).Update("language", lang).Error
}

func DeleteUser(db *gorm.DB, userId int64) error {
	return db.Delete(&User{ID: userId}).Error
}
