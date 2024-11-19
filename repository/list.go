package repository

import "gorm.io/gorm"

type List struct {
	ID      int64  `gorm:"autoIncrement"`
	OwnerID int64  `gorm:"not null"`
	Title   string `gorm:"not null"`
	Open    bool
	Access  []User `gorm:"many2many:list_access;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

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

type UserAccess struct {
	User
	HasAccess bool
}

func GetFriendsAccessListForList(db *gorm.DB, userId, listId, page int64) (users []UserAccess, err error) {
	err = db.Raw(`SELECT *, id IN (SELECT user_id FROM list_access WHERE list_id = $2) as has_access FROM users 
         WHERE id IN (SELECT friend_id FROM friends WHERE user_id = $1)
         LIMIT 6 OFFSET $3`, userId, listId, page*6).Scan(&users).Error
	return
}

func GrantAccess(db *gorm.DB, listId, friendId int64) error {
	err := db.Model(&List{ID: listId}).Association("Access").Append(&User{ID: friendId})
	return err
}
