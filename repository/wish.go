package repository

import "gorm.io/gorm"

type Wish struct {
	ID              int64  `gorm:"autoIncrement"`
	ListID          int64  `gorm:"not null"`
	List            List   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Title           string `gorm:"not null"`
	Url             string
	Description     string
	Price           string
	ReservationFree bool
	ReservedBy      int64
}

func AddWish(db *gorm.DB, userId, listId int64, title string) (int64, error) {
	wish := Wish{ListID: listId, Title: title}
	err := db.Omit("reserved_by").Create(&wish).Error
	return wish.ID, err
}

func GetWish(db *gorm.DB, wishId int64) (wish Wish, err error) {
	err = db.First(&wish, wishId).Error
	return
}

func UpdateWish(db *gorm.DB, wish Wish) error {
	return db.Save(&wish).Error
}

func UpdateWishField(db *gorm.DB, wishId int64, fieldName string, value any) error {
	return db.Model(&Wish{ID: wishId}).Update(fieldName, value).Error
}

func ReserveWish(db *gorm.DB, wishId, friendId int64) error {
	return db.Exec(`UPDATE wishes SET reserved_by = $1 WHERE id = $2 and reservation_free = false`, friendId, wishId).Error
	//return db.Model(&Wish{ID: wishId}).Update("reserved_by", friendId).Error
}

func UndoReservation(db *gorm.DB, wishId int64) error {
	return db.Model(&Wish{ID: wishId}).Update("reserved_by", nil).Error // TODO test null
}

func MakeWishReservationFree(db *gorm.DB, wishId int64) error {
	return db.Exec(`UPDATE wishes SET reservation_free = true, reserved_by = null WHERE id = $1`, wishId).Error
}

func MakeWishReservable(db *gorm.DB, wishId int64) error {
	return db.Exec(`UPDATE wishes SET reservation_free = false WHERE id = $1`, wishId).Error
}

func DeleteWish(db *gorm.DB, wishId int64) error {
	return db.Delete(&Wish{ID: wishId}).Error
}
