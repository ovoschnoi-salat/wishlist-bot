package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"wishlist_bot/storage"
)

func NewPgStorage(addr string) (Storage, error) {
	// urlExample := "postgres://postgres:postgres@localhost:5432/postgres"
	pool, err := pgxpool.New(context.Background(), addr)
	if err != nil {
		return Storage{}, err
	}
	return Storage{pool: pool}, nil
}

type Storage struct {
	pool *pgxpool.Pool
}

func (s Storage) Acquire() (StorageConn, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return StorageConn{}, err
	}
	return StorageConn{conn: conn}, nil
}

func (s Storage) Close() {
	s.pool.Close()
}

func (s Storage) GetWishlist(id int64, page int64) ([]storage.Wish, error) {
	conn, err := s.Acquire()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	return conn.GetWishlist(id, page)
}

func (s Storage) GetWishlistSize(id int64) (size int64, err error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	row := conn.QueryRow(context.Background(), `SELECT count(id) FROM wishes WHERE user_id = $1`, id)
	err = row.Scan(&size)
	return
}

func (s Storage) GetFriendsList(id, page int64) ([]storage.User, error) {
	conn, err := s.Acquire()
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	return conn.GetFriendsList(id, page)
}

func (s Storage) AddUser(id int64, username string) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.AddUser(id, username)
}

func (s Storage) DeleteUser(id int64) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.DeleteUser(id)
}

func (s Storage) AddWish(userId int64, title, url string) (int64, error) {
	conn, err := s.Acquire()
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	return conn.AddWish(userId, title, url)
}

func (s Storage) GetWish(wishId int64) (storage.Wish, error) {
	conn, err := s.Acquire()
	if err != nil {
		return storage.Wish{}, err
	}
	defer conn.Release()
	return conn.GetWish(wishId)
}

func (s Storage) EditWish(wishId int64, title, url string) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.EditWish(wishId, title, url)
}

func (s Storage) EditWishTitle(wishId int64, title string) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.EditWishTitle(wishId, title)
}

func (s Storage) EditWishUrl(wishId int64, url string) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.EditWishUrl(wishId, url)
}

func (s Storage) ReserveWish(wishId, userId int64) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.ReserveWish(wishId, userId)
}

func (s Storage) UndoReservation(wishId int64) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.UndoReservation(wishId)
}

func (s Storage) DeleteWish(wishId int64) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.DeleteWish(wishId)
}

func (s Storage) AddFriend(id, friendId int64) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.AddFriend(id, friendId)
}

func (s Storage) DeleteFriend(id, friendId int64) error {
	conn, err := s.Acquire()
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.DeleteFriend(id, friendId)
}
