package repository

//type DB interface {
//	// User
//	AddUser(userId int64, username, lang string) error
//	GetUserById(userId int64) (User, error)
//	GetUserByUsername(username string) (User, error)
//	UpdateUserLanguage(userId int64, lang string) error
//	DeleteUser(userId int64) error
//	// List
//	AddList(userId int64, title string) (int64, error)
//	UpdateListTitle(listId int64, title string) error
//	DeleteList(listId int64) error
//	GetListSize(listId int64) (int64, error)
//	GetUserLists(userId int64) ([]List, error)
//	GetFriendLists(userId, friendId int64) ([]List, error)
//	GetAvailableFriendsForList(userId, listId, page int64) ([]User, error)
//	GrantAccess(listId, friendId, userId int64) error
//	// Wish
//	GetWishes(listId, page int64) ([]Wish, error)
//	AddWish(userId, listId int64, title string) (int64, error)
//	GetWish(wishId int64) (Wish, error)
//	UpdateWish(wishId int64, title, url, description, price string) error
//	ReserveWish(wishId, friendId int64) error
//	UndoReservation(wishId int64) error
//	DeleteWish(wishId int64) error
//	// Friends
//	GetListOfFriends(userId, page int64) ([]User, error)
//	GetListOfFriendsSize(userId int64) (int64, error)
//	AddFriend(userId, friendId int64) error
//	DeleteFriend(userId, friendId int64) error
//}
