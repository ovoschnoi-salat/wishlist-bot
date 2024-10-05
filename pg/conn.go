package pg

//
//import (
//	"context"
//	"fmt"
//	"github.com/jackc/pgx/v5"
//	"github.com/jackc/pgx/v5/pgconn"
//	"github.com/jackc/pgx/v5/pgxpool"
//	"strings"
//	"wishlist_bot/repository"
//)
//
//type StorageConn struct {
//	conn *pgxpool.Conn
//}
//
//func (s StorageConn) UpdateUserLanguage(userId int64, lang string) error {
//	_, err := s.exec("UPDATE users SET language = $2 WHERE user_id = $1", userId, lang)
//	return err
//}
//
//func (s StorageConn) GetAvailableFriendsForList(userId, listId, page int64) (users []repository.User, err error) {
//	rows, err := s.query(`select * from users natural join
//(select friend_id as user_id from friends where user_id = $1 and friend_id not in
//(select friend_id from list_access where list_id = $2) order by user_id LIMIT 6 OFFSET $3) f`, userId, listId, page*6)
//	if err != nil {
//		return
//	}
//	for rows.Next() {
//		u := repository.User{}
//		err = rows.Scan(&u.UserId, &u.Username)
//		if err != nil {
//			return
//		}
//		users = append(users, u)
//	}
//	return
//}
//
//func (s StorageConn) AddList(userId int64, title string) (listId int64, err error) {
//	row := s.queryRow("INSERT INTO lists (user_id, title) VALUES ($1, $2) RETURNING list_id", userId, title)
//	err = row.Scan(&listId)
//	return
//}
//
//func (s StorageConn) UpdateListTitle(listId int64, title string) error {
//	_, err := s.exec("UPDATE lists SET title = $2 WHERE list_id = $1", listId, title)
//	return err
//}
//
//func (s StorageConn) DeleteList(listId int64) error {
//	_, err := s.exec("DELETE FROM lists WHERE list_id = $1;", listId)
//	return err
//}
//
//func (s StorageConn) GetListSize(listId int64) (size int64, err error) {
//	err = s.queryRow("SELECT count(*) FROM wishes WHERE list_id = $1", listId).Scan(&size)
//	return
//}
//
//func (s StorageConn) GetUserLists(userId int64) (lists []repository.List, err error) {
//	rows, err := s.query("SELECT list_id, user_id, title FROM lists WHERE user_id = $1 ORDER BY list_id", userId)
//	if err != nil {
//		return
//	}
//	lists = make([]repository.List, 0, 6)
//	for rows.Next() {
//		var l repository.List
//		err = rows.Scan(&l.ListId, &l.UserId, &l.Title)
//		if err != nil {
//			return nil, err
//		}
//		lists = append(lists, l)
//	}
//	return
//}
//
//func (s StorageConn) GetFriendLists(userId, friendId int64) (lists []repository.List, err error) {
//	rows, err := s.query(`SELECT list_id, user_id, title
//FROM lists natural left join list_access a
//WHERE lists.user_id = $1 and (a.friend_id = $2 or lists.open) ORDER BY list_id`, friendId, userId)
//	if err != nil {
//		return
//	}
//	lists = make([]repository.List, 0, 6)
//	for rows.Next() {
//		var l repository.List
//		err = rows.Scan(&l.ListId, &l.UserId, &l.Title)
//		if err != nil {
//			return nil, err
//		}
//		lists = append(lists, l)
//	}
//	return
//}
//
//func (s StorageConn) GrantAccess(listId, friendId, userId int64) (err error) {
//	_, err = s.exec("INSERT INTO list_access (list_id, friend_id, user_id) VALUES ($1, $2, $3)", listId, friendId, userId)
//	return
//}
//
//func (s StorageConn) AddWish(userId, listId int64, title string) (wishId int64, err error) {
//	row := s.queryRow("INSERT INTO wishes (user_id, list_id, title) VALUES ($1, $2, $3) RETURNING wish_id", userId, listId, title)
//	err = row.Scan(&wishId)
//	return
//}
//
//func (s StorageConn) GetWishes(listId, page int64) (wishes []repository.Wish, err error) {
//	rows, err := s.query("SELECT wish_id, title, url, reserved_by FROM wishes WHERE list_id = $1 ORDER BY wish_id LIMIT 6 OFFSET $2",
//		listId, 6*page)
//	if err != nil {
//		return
//	}
//	for rows.Next() {
//		var w repository.Wish
//		err = rows.Scan(&w.UserId, &w.Title, &w.Url, &w.ReservedBy)
//		if err != nil {
//			return nil, err
//		}
//		wishes = append(wishes, w)
//	}
//	return
//}
//
//func (s StorageConn) AddUser(userId int64, username, lang string) error {
//	_, err := s.exec("INSERT INTO users (user_id, username, language) VALUES ($1, $2, $3)", userId, username, lang)
//	return err
//}
//
//func (s StorageConn) GetUserById(userId int64) (user repository.User, err error) {
//	row := s.queryRow("SELECT * FROM users WHERE user_id = $1", userId)
//	err = row.Scan(&user.UserId, &user.Username)
//	return
//}
//
//func (s StorageConn) GetUserByUsername(username string) (user repository.User, err error) {
//	row := s.queryRow("SELECT * FROM users WHERE username = $1", username)
//	err = row.Scan(&user.UserId, &user.Username)
//	return
//}
//
//func (s StorageConn) DeleteUser(userId int64) error {
//	_, err := s.exec("DELETE FROM users WHERE user_id = $1;", userId)
//	return err
//}
//
//func (s StorageConn) GetWish(wishId int64) (wish repository.Wish, err error) {
//	row := s.queryRow("SELECT * FROM wishes WHERE wish_id = $1", wishId)
//	err = row.Scan(&wish.WishId, &wish.UserId, &wish.Title, &wish.Url, &wish.ReservedBy)
//	return
//}
//
//func (s StorageConn) UpdateWish(wishId int64, title, url, description, price string) error {
//	sb := strings.Builder{}
//	sb.WriteString("UPDATE wishes SET")
//	addedFields := false
//	if title != "" {
//		sb.WriteString(" title = $1")
//		addedFields = true
//	}
//	if url != "" {
//		if addedFields {
//			sb.WriteString(",")
//		}
//		sb.WriteString(" url = $2")
//		addedFields = true
//	}
//	if description != "" {
//		if addedFields {
//			sb.WriteString(",")
//		}
//		sb.WriteString(" description = $3")
//		addedFields = true
//	}
//	if price != "" {
//		if addedFields {
//			sb.WriteString(",")
//		}
//		sb.WriteString(" price = $4")
//		addedFields = true
//	}
//	if !addedFields {
//		return fmt.Errorf("no fields to update")
//	}
//	sb.WriteString(" WHERE wish_id = $1")
//	_, err := s.exec(sb.String(), wishId, title, url, description, price)
//	return err
//}
//
//func (s StorageConn) ReserveWish(wishId, userId int64) error {
//	_, err := s.exec("UPDATE wishes SET reserved_by = $2 WHERE wish_id = $1", wishId, userId)
//	return err
//}
//
//func (s StorageConn) UndoReservation(wishId int64) error {
//	_, err := s.exec("UPDATE wishes SET reserved_by = null WHERE wish_id = $1", wishId)
//	return err
//}
//
//func (s StorageConn) DeleteWish(wishId int64) error {
//	_, err := s.exec("DELETE FROM wishes WHERE wish_id = $1", wishId)
//	return err
//}
//
//func (s StorageConn) GetWishlistSize(id int64) (wishlistSize int64, err error) {
//	err = s.queryRow("SELECT count(*) FROM wishes WHERE user_id = $1", id).Scan(&wishlistSize)
//	return
//}
//
//func (s StorageConn) AddFriend(id, friendId int64) error {
//	_, err := s.exec("INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)", id, friendId)
//	return err
//}
//
//func (s StorageConn) DeleteFriend(id, friendId int64) error {
//	_, err := s.exec("DELETE FROM friends WHERE user_id = $1 AND friend_id = $2 OR user_id = $2 AND friend_id = $1",
//		id, friendId)
//
//	return err
//}
//
//func (s StorageConn) GetListOfFriends(id, page int64) ([]repository.User, error) {
//	// EXPLAIN SELECT U.id, U.username FROM users U JOIN (SELECT friend_id from friends
//	// where user_id = 1 order by friend_id LIMIT 6) F ON U.id = F.friend_id
//	rows, err := s.query("SELECT users.user_id, username FROM friends natural JOIN users WHERE friends.friend_id = $1 ORDER BY users.user_id LIMIT 6 OFFSET $2", id, 6*page)
//	if err != nil {
//		return nil, err
//	}
//	var res []repository.User
//	for rows.Next() {
//		var u repository.User
//		err := rows.Scan(&u.UserId, &u.Username)
//		if err != nil {
//			return nil, err
//		}
//		res = append(res, u)
//	}
//	if rows.Err() != nil {
//		return nil, rows.Err()
//	}
//	return res, nil
//}
//
//func (s StorageConn) GetListOfFriendsSize(id int64) (wishlistSize int64, err error) {
//	err = s.queryRow("SELECT count(friend_id) FROM friends WHERE user_id = $1", id).Scan(&wishlistSize)
//	return
//}
//
//func (s StorageConn) exec(command string, args ...any) (pgconn.CommandTag, error) {
//	return s.conn.Exec(context.Background(), command, args...)
//}
//
//func (s StorageConn) query(command string, args ...any) (pgx.Rows, error) {
//	return s.conn.Query(context.Background(), command, args...)
//}
//
//func (s StorageConn) queryRow(command string, args ...any) pgx.Row {
//	return s.conn.QueryRow(context.Background(), command, args...)
//}
//
//func (s StorageConn) Release() {
//	s.conn.Release()
//}
