package pg

//
//import (
//	"context"
//	"github.com/jackc/pgx/v5/pgxpool"
//	"log"
//)
//
//func createTables(conn *pgxpool.Conn) error {
//	_, err := conn.Exec(context.Background(), `CREATE TABLE users (
//    user_id 	bigint primary key,
//    username 	varchar,
//	language 	varchar not null
//    )`)
//	if err != nil {
//		log.Println("error creating users table:", err)
//	}
//	_, err = conn.Exec(context.Background(), `CREATE TABLE lists (
//    list_id 	bigserial primary key,
//    user_id 	bigint not null,
//    title 		varchar not null,
//    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
//	)`)
//	if err != nil {
//		log.Println("error creating wishes table:", err)
//	}
//	_, err = conn.Exec(context.Background(), `CREATE TABLE list_access (
//    list_id 	bigint,
//    friend_id 	bigint,
//    user_id 	bigint not null,
//    open 		boolean not null default true,
//    FOREIGN KEY (list_id) REFERENCES lists(list_id) ON DELETE CASCADE,
//    FOREIGN KEY (friend_id) REFERENCES users(user_id) ON DELETE CASCADE,
//    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
//    CONSTRAINT "ID_PKEY" PRIMARY KEY (list_id,friend_id)
//	)`)
//	if err != nil {
//		log.Println("error creating wishes table:", err)
//	}
//	_, err = conn.Exec(context.Background(), `CREATE TABLE wishes (
//    wish_id 			bigserial primary key,
//    user_id 			bigint not null,
//    list_id 			bigint not null,
//    title 				varchar not null,
//    url 				varchar,
//    description     	varchar,
//	Price           	varchar,
//	reservation_free 	boolean not null default false,
//    reserved_by 		bigint,
//    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
//    FOREIGN KEY (list_id) REFERENCES lists(list_id) ON DELETE CASCADE,
//    FOREIGN KEY (reserved_by) REFERENCES users(user_id) ON DELETE SET NULL
//	)`)
//	if err != nil {
//		log.Println("error creating wishes table:", err)
//	}
//	_, err = conn.Exec(context.Background(), `CREATE TABLE friends (
//    user_id 	bigint primary key,
//    friend_id 	bigint not null,
//    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
//    FOREIGN KEY (friend_id) REFERENCES users(user_id) ON DELETE CASCADE
//    )`)
//	if err != nil {
//		log.Println("error creating wishes table:", err)
//	}
//	return nil
//}
