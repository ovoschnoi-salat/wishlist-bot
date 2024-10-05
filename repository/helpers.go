package repository

import (
	"gorm.io/gorm"
)

func AddList(db *gorm.DB, list *List) error {
	return db.Create(list).Error
}

func GetListById(db *gorm.DB, listId int64) (list List, err error) {
	err = db.First(&list, listId).Error
	return
}

func UpdateListField(db *gorm.DB, listId int64, fieldName string, value any) error {
	return db.Model(&List{ID: listId}).Update(fieldName, value).Error
}

func UpdateList(db *gorm.DB, list *List) error {
	return db.Save(list).Error
}

func DeleteList(db *gorm.DB, listId int64) error {
	return db.Delete(&List{ID: listId}).Error
}

func GetListSize(db *gorm.DB, listId int64) (size int64, err error) {
	err = db.Raw("SELECT count(*) FROM wishes WHERE list_id = ?", listId).Scan(&size).Error
	return
}

func GetUserLists(db *gorm.DB, userId int64) (lists []List, err error) {
	err = db.Where("owner_id =?", userId).Find(&lists).Error
	return
}

func CountUserLists(db *gorm.DB, userId int64) (size int64, err error) {
	ass := db.Model(&User{ID: userId}).Association("Lists")
	size = ass.Count()
	err = ass.Error
	return
}

func GetFriendLists(db *gorm.DB, userId, friendId int64) (lists []List, err error) {
	err = db.Raw(`SELECT * FROM lists WHERE owner_id = $1
           AND (id IN (SELECT list_id FROM list_access WHERE user_id = $2) OR open = true)`,
		friendId, userId).Scan(&lists).Error
	return
}

func ClearListAccess(db *gorm.DB, listId int64) error {
	return db.Exec(`DELETE FROM list_access WHERE list_id = $1`, listId).Error
}

func GetAvailableFriendsForList(db *gorm.DB, userId, listId, page int64) (users []User, err error) {
	err = db.Raw(`SELECT * FROM users 
         WHERE id IN (SELECT friend_id FROM friends WHERE user_id = $1) 
           AND id NOT IN (SELECT user_id FROM list_access WHERE list_id = $2) 
         LIMIT 6 OFFSET $3`, userId, listId, page*6).Scan(&users).Error
	return
}

func GetSizeOfAvailableFriendsForList(db *gorm.DB, userId, listId int64) (size int64, err error) {
	err = db.Raw(`SELECT count(*) FROM users 
         WHERE id IN (SELECT friend_id FROM friends WHERE user_id = $1) 
           AND id NOT IN (SELECT user_id FROM list_access WHERE list_id = $2)`, userId, listId).Scan(&size).Error
	return
}

func GrantAccess(db *gorm.DB, listId, friendId, userId int64) error {
	err := db.Model(&List{ID: listId}).Association("Access").Append(&User{ID: friendId})
	return err
}

func GetWishes(db *gorm.DB, listId, page int64) (wishes []Wish, err error) {
	err = db.Where("list_id = ?", listId).Offset(int(page * 6)).Limit(6).Find(&wishes).Error
	return
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

// Бот для создания списков желаний и бронирования подарков друзьями. Делитесь своими желаниями и избегайте дублирующихся подарков!
