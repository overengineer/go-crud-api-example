package main

import (
	"os"
	"log"
	"fmt"
	"strconv"
	
	"net/http"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:artist`
	Price  float64 `json:"price"`
}

type Database struct {
	inst *sql.DB
}

type App struct {
	db Database
}

func (db Database) fetchAlbumByID(id int64) (Album, error) {
	var alb Album
	row := db.inst.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsByID %d: no suchalbum", id)
		}
		return alb, fmt.Errorf("albumsByID %d: %v", id, err)
	}
	return alb, nil
}

func (db Database) fetchAlbums() ([]Album, error) {
	var albums []Album

	rows, err := db.inst.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("fetchAlbums: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("fetchAlbums: %v", err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetchAlbums: %v", err)
	}

	return albums, nil
}

func getDatabase() (Database, error) {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		err = db.Ping()
	}
	if err != nil {
		return Database{nil}, err
	}
	return Database{db}, nil
}

func (app App) getAlbums(c *gin.Context) {
	var albums []Album
	var err error
	if albums, err = app.db.fetchAlbums(); err == nil {
		log.Printf("Albums: %v", albums)
		c.IndentedJSON(http.StatusOK, albums)
		return
	}
	log.Printf("Error: %v", err)
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Internal Error"})
}

func (app App) getAlbumByID(c *gin.Context) {
	var album Album
	var err error
	var id int64
	if id, err = strconv.ParseInt(c.Param("id"), 10, 64); err == nil {
		if album, err = app.db.fetchAlbumByID(id); err == nil {
			log.Printf("Album: %v", album)
			c.IndentedJSON(http.StatusOK, album)
			return
		}
	}
	log.Printf("Error: %v", err)
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func initApp() (App, error) {
	var db Database
	var err error
	db, err = getDatabase()
	return App{db}, err
}

// func postAlbums(c *gin.Context) {
// 	var newAlbum album

// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}

// 	albums = append(albums, newAlbum)
// 	c.IndentedJSON(http.StatusCreated, newAlbum)
// }

func main() {
	app, err := initApp()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to database")

	router := gin.Default()
	router.GET("/albums", app.getAlbums)
	router.GET("/albums/:id", app.getAlbumByID)
	// router.POST("/albums", postAlbums)
	
	router.Run("localhost:8080")
}
