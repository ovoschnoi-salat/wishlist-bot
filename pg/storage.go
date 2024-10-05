package pg

//
//import (
//	"context"
//	"github.com/jackc/pgx/v5/pgxpool"
//	"log"
//	"wishlist_bot/repository"
//)
//
//func NewPgStorage(addr string) (repository.Storage, error) {
//	// url_example := "postgres://postgres:postgres@localhost:5432/postgres"
//	pool, err := pgxpool.New(context.Background(), addr)
//	if err != nil {
//		return Storage{}, err
//	}
//	conn, err := pool.Acquire(context.Background())
//	if err != nil {
//		return Storage{}, err
//	}
//	_ = createTables(conn)
//	conn.Release()
//	return Storage{pool: pool}, nil
//}
//
//type Storage struct {
//	pool *pgxpool.Pool
//}
//
//func (s Storage) GetAvailableFriendsForList(userId, listId, page int64) ([]repository.User, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Release()
//	return conn.GetAvailableFriendsForList(userId, listId, page)
//}
//
//func (s Storage) Acquire() (repository.Conn, error) {
//	conn, err := s.pool.Acquire(context.Background())
//	if err != nil {
//		return StorageConn{}, err
//	}
//	return StorageConn{conn: conn}, nil
//}
//
//func (s Storage) Close() {
//	log.Println("waiting releases:", s.pool.Stat().AcquiredConns())
//	s.pool.Close()
//}
//
//func (s Storage) AddUser(userId int64, username, lang string) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.AddUser(userId, username, lang)
//}
//
//func (s Storage) GetUserById(userId int64) (repository.User, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return repository.User{}, err
//	}
//	defer conn.Release()
//	return conn.GetUserById(userId)
//}
//
//func (s Storage) GetUserByUsername(username string) (repository.User, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return repository.User{}, err
//	}
//	defer conn.Release()
//	return conn.GetUserByUsername(username)
//}
//
//func (s Storage) DeleteUser(id int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.DeleteUser(id)
//}
//
//func (s Storage) AddList(userId int64, title string) (int64, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return 0, err
//	}
//	defer conn.Release()
//	return conn.AddList(userId, title)
//}
//
//func (s Storage) UpdateListTitle(listId int64, title string) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.UpdateListTitle(listId, title)
//}
//
//func (s Storage) DeleteList(listId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.DeleteList(listId)
//}
//
//func (s Storage) GetListSize(listId int64) (int64, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return 0, err
//	}
//	defer conn.Release()
//	return conn.GetListSize(listId)
//}
//
//func (s Storage) GetUserLists(userId int64) ([]repository.List, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Release()
//	return conn.GetUserLists(userId)
//}
//
//func (s Storage) GetFriendLists(userId, friendId int64) ([]repository.List, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Release()
//	return conn.GetFriendLists(userId, friendId)
//}
//
//func (s Storage) GrantAccess(listId, friendId, userId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.GrantAccess(listId, friendId, userId)
//}
//
//func (s Storage) AddWish(userId, listId int64, title string) (int64, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return 0, err
//	}
//	defer conn.Release()
//	return conn.AddWish(userId, listId, title)
//}
//
//func (s Storage) GetWishes(listId, page int64) ([]repository.Wish, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Release()
//	return conn.GetWishes(listId, page)
//}
//
//func (s Storage) GetWishlistSize(listId int64) (int64, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Storage) GetListOfFriends(id, page int64) ([]repository.User, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Release()
//	return conn.GetListOfFriends(id, page)
//}
//
//func (s Storage) GetListOfFriendsSize(userId int64) (int64, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return 0, err
//	}
//	defer conn.Release()
//	return conn.GetListOfFriendsSize(userId)
//}
//
//func (s Storage) GetWish(wishId int64) (repository.Wish, error) {
//	conn, err := s.Acquire()
//	if err != nil {
//		return repository.Wish{}, err
//	}
//	defer conn.Release()
//	return conn.GetWish(wishId)
//}
//
//func (s Storage) UpdateWish(wishId int64, title, url, description, price string) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.UpdateWish(wishId, title, url, description, price)
//}
//
//func (s Storage) ReserveWish(wishId, userId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.ReserveWish(wishId, userId)
//}
//
//func (s Storage) UndoReservation(wishId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.UndoReservation(wishId)
//}
//
//func (s Storage) DeleteWish(wishId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.DeleteWish(wishId)
//}
//
//func (s Storage) AddFriend(userId, friendId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.AddFriend(userId, friendId)
//}
//
//func (s Storage) DeleteFriend(userId, friendId int64) error {
//	conn, err := s.Acquire()
//	if err != nil {
//		return err
//	}
//	defer conn.Release()
//	return conn.DeleteFriend(userId, friendId)
//}
