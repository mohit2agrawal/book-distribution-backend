package main

import (
	"backend/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
var chanOfChans []chan db.Transaction

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
}

func getUsers(c *gin.Context) {
	users, err := db.AllUsers()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Fetch failure!"})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	user, err := db.GetUser(int(id))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func addUser(c *gin.Context) {
	var newUser db.User

	if err := c.BindJSON(&newUser); err != nil {
		fmt.Println(err)
		return
	}

	_, err := db.CreateUser(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create a user"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func deleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = db.DeleteUser(int(id))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
	} else {
		c.Status(http.StatusAccepted)
	}
}

func getBooks(c *gin.Context) {
	books, err := db.AllBooks()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Fetch failure!"})
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	book, err := db.GetBook(int(id))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func addBook(c *gin.Context) {
	var newBook db.Book

	if err := c.BindJSON(&newBook); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("inserting book", newBook)

	_, err := db.CreateBook(&newBook)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create a book"})
		return
	}
	fmt.Println("inserted book", newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}

func deleteBook(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = db.DeleteBook(int(id))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
	} else {
		c.Status(http.StatusAccepted)
	}
}

func getTransactions(c *gin.Context) {
	trxs, err := db.AllTransactions()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Fetch failure!"})
		return
	}
	c.IndentedJSON(http.StatusOK, trxs)
}

func getTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	trx, err := db.GetTransaction(int(id))

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, trx)
}

func addTransaction(c *gin.Context) {
	var newTrx db.Transaction

	if err := c.BindJSON(&newTrx); err != nil {
		fmt.Println(err)
		return
	}

	_, err := db.CreateTransaction(&newTrx)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create a transaction"})
		return
	}

	for _, c := range chanOfChans {
		log.Println("sending to chan...")
		c <- newTrx
	}

	c.IndentedJSON(http.StatusCreated, newTrx)
}

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	c := make(chan db.Transaction)
	chanOfChans = append(chanOfChans, c)

	for tData := range c {

		log.Println("got from chan...")

		data, _ := json.Marshal(tData)

		err = ws.WriteMessage(1, data)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

	// for {
	// 	mt, message, err := ws.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)

	// 	data, _ := json.Marshal(db.User{1, "name", "phone"})

	// 	err = ws.WriteMessage(mt, data)
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }
}

func main() {
	dbPath, _ := filepath.Abs("store.db")
	must(db.Init(dbPath))

	// users = append(users, []User{{"1", "Mohit", "101"}, {"2", "Shreyank", "102"}}...)
	// books = append(books, []Book{{"3", "Krsna, The Supreme Personality of Godhead", English, 40}, {"1", "Bhagavad Gita", English, 250}, {"2", "Krsna", Hindi, 300}}...)
	// transactions = append(transactions, []Transaction{{"t1", time.Now(), "1", "1", 100}, {"t2", time.Now(), "1", "2", 450}}...)

	// fmt.Println(users)
	// fmt.Println(books)
	// fmt.Println(transactions)

	r := gin.Default()

	r.GET("/users", getUsers)
	r.GET("/users/:id", getUser)
	r.DELETE("/users/:id", deleteUser)
	r.POST("/users", addUser)

	r.GET("/books", getBooks)
	r.GET("/books/:id", getBook)
	r.DELETE("/books/:id", deleteBook)
	r.POST("/books", addBook)

	r.GET("/transactions", getTransactions)
	r.GET("/transactions/:id", getTransaction)
	r.POST("/transactions", addTransaction)
	r.GET("/", ping)

	r.GET("/echo", gin.WrapF(echo))

	r.Run("localhost:5050")
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
