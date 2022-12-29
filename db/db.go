package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var userBucket = []byte("users")
var bookBucket = []byte("books")
var txBucket = []byte("transactions")
var db *bolt.DB

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
	Price    int    `json:"price"`
}

type Transaction struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    int       `json:"userId"`
	BookID    int       `json:"bookId"`
	Amount    int       `json:"amount"`
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(bookBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(txBucket)
		return err
	})
}

func CreateUser(user *User) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)

		user.ID = id
		val, _ := json.Marshal(user)
		return b.Put(key, val)
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllUsers() ([]User, error) {
	var users []User
	var user User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			_ = json.Unmarshal(v, &user)
			users = append(users, user)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(id int) (User, error) {
	var user User

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		data := b.Get(itob(id))
		err := json.Unmarshal(data, &user)
		return err
	})
	return user, err
}

func DeleteUser(id int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(userBucket)
		return b.Delete(itob(id))
	})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func CreateBook(book *Book) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bookBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)

		book.ID = id
		fmt.Println("inserting book...", book)
		val, _ := json.Marshal(book)
		return b.Put(key, val)
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllBooks() ([]Book, error) {
	var books []Book
	var book Book
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bookBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			_ = json.Unmarshal(v, &book)
			books = append(books, book)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return books, nil
}

func GetBook(id int) (Book, error) {
	var book Book

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bookBucket)
		data := b.Get(itob(id))
		err := json.Unmarshal(data, &book)
		return err
	})
	return book, err
}

func DeleteBook(id int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bookBucket)
		return b.Delete(itob(id))
	})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func CreateTransaction(trasact *Transaction) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(txBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)

		trasact.ID = id
		val, _ := json.Marshal(trasact)
		return b.Put(key, val)
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllTransactions() ([]Transaction, error) {
	var trxs []Transaction
	var trx Transaction
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(txBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			_ = json.Unmarshal(v, &trx)
			trxs = append(trxs, trx)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return trxs, nil
}

func GetTransaction(id int) (Transaction, error) {
	var trx Transaction

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(txBucket)
		data := b.Get(itob(id))
		err := json.Unmarshal(data, &trx)
		return err
	})
	return trx, err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
