package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/libsql/libsql-client-go/libsql"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var db *sql.DB

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	token := os.Getenv("DB_TOKEN")

	var dbUrl = fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, token)

	var err error
	db, err = sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal("failed to open db ", dbUrl, err)
		os.Exit(1)
	}
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

func getAlbums(c *gin.Context) {
	var albums []album
	var (
		Album album
	)
	rows, err := db.Query("SELECT id, title, artist, price FROM albums;")
	if err != nil {
		log.Printf("Cannot fetch from the Albums table. Error: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to query albums table"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&Album.ID, &Album.Title, &Album.Artist, &Album.Price)
		if err != nil {
			log.Printf("Error reading row from Albums table; Error: %s", err)
		}
		albums = append(albums, Album)
	}
	if len(albums) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No albums found"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id should be a numeric value"})
		return
	}
	fetchOneQuery := fmt.Sprintf("SELECT id, title, artist, price FROM albums WHERE id == %d", id)
	var (
		Album album
	)
	if err := db.QueryRow(fetchOneQuery).Scan(&Album.ID, &Album.Title, &Album.Artist, &Album.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to fetch Album"})
		log.Printf("Unable to Fetch Album; Error: %s", err)
		return
	}

	c.IndentedJSON(http.StatusOK, Album)
}

func postAlbums(c *gin.Context) {

	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	query := "INSERT INTO albums (id, title, artist, price) VALUES (?, ?, ?, ?)"
	_, err := db.ExecContext(context.Background(), query, newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	if err != nil {
		log.Printf("Unable to insert; Error: %s", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to insert to albums table."})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
