package repository

import "gorm.io/gorm"

func GetListOfFriends(db *gorm.DB, userId, page int64) (users []User, err error) {
	err = db.Model(&User{ID: userId}).Offset(int(page * 6)).Limit(6).Association("Friends").Find(&users)
	return
	//err = db.Raw(`select * from users where id in
	//(select friend_id from friends where user_id = $1 order by friend_id LIMIT 6 OFFSET $2)`, userId, page*6).Scan(&users).Error
	//return
}

func GetListOfFriendsSize(db *gorm.DB, userId int64) (size int64, err error) {
	ass := db.Model(&User{ID: userId}).Association("Friends")
	size = ass.Count()
	return size, ass.Error
}

func AddFriend(db *gorm.DB, userId, friendId int64) error {
	err := db.Model(&User{ID: userId}).Association("Friends").Append(&User{ID: friendId})
	if err != nil {
		return err
	}
	return db.Model(&User{ID: friendId}).Association("Friends").Append(&User{ID: userId})
}

func DeleteFriend(db *gorm.DB, userId, friendId int64) error {
	err := db.Model(&User{ID: userId}).Association("Friends").Delete(&User{ID: friendId})
	if err != nil {
		return err
	}
	return db.Model(&User{ID: friendId}).Association("Friends").Delete(&User{ID: userId})
}
