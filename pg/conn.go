package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"wishlist_bot/storage"
)

type StorageConn struct {
	conn *pgxpool.Conn
}

func (s StorageConn) AddUser(id int64, username string) error {
	_, err := s.exec("INSERT INTO users (id, username) VALUES ($1, $2)", id, username)
	return err
}

func (s StorageConn) GetUserById(id int64) (user storage.User, err error) {
	row := s.queryRow("SELECT * FROM users WHERE id = $1", id)
	err = row.Scan(&user.Id, &user.Username)
	return
}

func (s StorageConn) GetUserByUsername(username string) (user storage.User, err error) {
	row := s.queryRow("SELECT * FROM users WHERE username = $1", username)
	err = row.Scan(&user.Id, &user.Username)
	return
}

func (s StorageConn) DeleteUser(id int64) error {
	_, err := s.exec("DELETE FROM users WHERE id = $1;", id)
	return err
}

func (s StorageConn) AddWish(userId int64, title, url string) (wishId int64, err error) {
	if url == "" {
		row := s.queryRow("INSERT INTO wishes (user_id, title) VALUES ($1, $2) RETURNING id", userId, title)
		err = row.Scan(&wishId)
	} else {
		row := s.queryRow("INSERT INTO wishes (user_id, title, url) VALUES ($1, $2, $3) RETURNING id", userId, title, url)
		err = row.Scan(&wishId)
	}
	return
}

func (s StorageConn) GetWish(wishId int64) (wish storage.Wish, err error) {
	row := s.queryRow("SELECT id, user_id, title, url, reserved_by FROM wishes WHERE id = $1", wishId)
	err = row.Scan(&wish.Id, &wish.Owner, &wish.Title, &wish.Url, &wish.ReservedBy)
	return
}

func (s StorageConn) EditWish(wishId int64, title, url string) error {
	if url == "" {
		return s.EditWishTitle(wishId, title)
	} else if title == "" {
		return s.EditWishUrl(wishId, url)
	}
	_, err := s.exec("UPDATE wishes SET title = $2, url = $3 WHERE id = $1", wishId, title, url)
	return err
}

func (s StorageConn) EditWishTitle(wishId int64, title string) error {
	_, err := s.exec("UPDATE wishes SET title = $2 WHERE id = $1", wishId, title)
	return err
}

func (s StorageConn) EditWishUrl(wishId int64, url string) error {
	if url == "" {
		_, err := s.exec("UPDATE wishes SET url = null WHERE id = $1", wishId)
		return err
	}
	_, err := s.exec("UPDATE wishes SET url = $2 WHERE id = $1", wishId, url)
	return err
}

func (s StorageConn) ReserveWish(wishId, userId int64) error {
	_, err := s.exec("UPDATE wishes SET reserved_by = $2 WHERE id = $1", wishId, userId)
	return err
}

func (s StorageConn) UndoReservation(wishId int64) error {
	_, err := s.exec("UPDATE wishes SET reserved_by = null WHERE id = $1", wishId)
	return err
}

func (s StorageConn) DeleteWish(wishId int64) error {
	_, err := s.exec("DELETE FROM wishes WHERE id = $1", wishId)
	return err
}

func (s StorageConn) GetWishlist(id int64, page int64) (wishes []storage.Wish, err error) {
	rows, err := s.query("SELECT id, title, url, reserved_by FROM wishes WHERE user_id = $1 ORDER BY id LIMIT 6 OFFSET $2",
		id, 6*page)
	if err != nil {
		return
	}
	for rows.Next() {
		var w storage.Wish
		err = rows.Scan(&w.Id, &w.Title, &w.Url, &w.ReservedBy)
		if err != nil {
			return nil, err
		}
		wishes = append(wishes, w)
	}
	return
}

func (s StorageConn) GetWishlistSize(id int64) (wishlistSize int64, err error) {
	err = s.queryRow("SELECT count(id) FROM wishes WHERE user_id = $1", id).Scan(&wishlistSize)
	return
}

func (s StorageConn) AddFriend(id int64, friendId int64) error {
	_, err := s.exec("INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)", id, friendId)
	return err
}

func (s StorageConn) DeleteFriend(id int64, friendId int64) error {
	_, err := s.exec("DELETE FROM friends WHERE user_id = $1 AND friend_id = $2 OR user_id = $2 AND friend_id = $1",
		id, friendId)

	return err
}

func (s StorageConn) GetFriendsList(id int64, page int64) ([]storage.User, error) {
	// EXPLAIN SELECT U.id, U.username FROM users U JOIN (SELECT friend_id from friends
	// where user_id = 1 order by friend_id LIMIT 6) F ON U.id = F.friend_id
	rows, err := s.query("SELECT id, username FROM friends JOIN users ON users.id = friends.friend_id WHERE friends.user_id = $1 ORDER BY users.id LIMIT 6 OFFSET $2", id, 6*page)
	if err != nil {
		return nil, err
	}
	var res []storage.User
	for rows.Next() {
		var u storage.User
		err := rows.Scan(&u.Id, &u.Username)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return res, nil
}

func (s StorageConn) GetFriendsListSize(id int64) (wishlistSize int64, err error) {
	err = s.queryRow("SELECT count(friend_id) FROM friends WHERE user_id = $1", id).Scan(&wishlistSize)
	return
}

func (s StorageConn) exec(command string, args ...any) (pgconn.CommandTag, error) {
	return s.conn.Exec(context.Background(), command, args...)
}

func (s StorageConn) query(command string, args ...any) (pgx.Rows, error) {
	return s.conn.Query(context.Background(), command, args...)
}

func (s StorageConn) queryRow(command string, args ...any) pgx.Row {
	return s.conn.QueryRow(context.Background(), command, args...)
}

func (s StorageConn) Release() {
	s.conn.Release()
}
