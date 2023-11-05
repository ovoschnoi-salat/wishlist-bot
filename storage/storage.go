package storage

import "database/sql"

type User struct {
	Id       int64
	Username string
}

type Wish struct {
	Id         int64
	Owner      int64
	Title      string
	Url        sql.NullString
	ReservedBy sql.NullInt64
}

type Storage interface {
	AddUser(id int64, username string) error
	GetUserById(id int64, username string) (User, error)
	GetUserByUsername(id int64, username string) (User, error)
	DeleteUser(id int64) error
	AddWish(userId int64, title, url string) (int64, error)
	GetWish(wishId int64) (Wish, error)
	EditWish(wishId int64, title, url string) error
	EditWishUrl(wishId int64, url string) error
	EditWishTitle(wishId int64, title string) error
	ReserveWish(wishId, userId int64) error
	UndoReservation(wishId int64) error
	DeleteWish(wishId int64) error
	GetWishlist(id, page int64) ([]Wish, error)
	GetWishlistSize(id int64) (int64, error)
	GetFriendsList(id, page int64) ([]User, error)
	GetFriendsListSize(id int64) (int64, error)
	AddFriend(id, friendId int64) error
	DeleteFriend(id, friendId int64) error
}
