package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Premium@007"
	dbname   = "postgres"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type urlShortener struct {
	Id      string `json:"id"`
	Url     string `json:"url"`
	UrlCode string `json:"url_code"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var urls = []urlShortener{}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	insertStmt := `insert into jovin.albums("id", "title","artist","price") values($1, $2,$3,$4)`
	_, e := db.Exec(insertStmt, newAlbum.ID, newAlbum.Title, newAlbum.Artist, newAlbum.Price)
	CheckError(e)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func postUrl(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	var newUrl urlShortener

	// Call BindJSON to bind the received JSON to
	// newAlbum
	if err := c.BindJSON(&newUrl); err != nil {
		return
	}

	// Add the new album to the slice.
	urls = append(urls, newUrl)
	insertStmt := `insert into jovin.url_details( "url") values($1)`
	_, e := db.Exec(insertStmt, newUrl.Url)
	CheckError(e)
	c.IndentedJSON(http.StatusCreated, newUrl)
}
func getUrlByID(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	var id = c.Param("id")
	var url string
	var urlCode string
	selectStmt := `select url,url_code from jovin.url_details where id=$1`
	row := db.QueryRow(selectStmt, id)
	switch err := row.Scan(&url, &urlCode); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		urls = nil
		fmt.Println(id)
		newUrl := urlShortener{Id: id, Url: url, UrlCode: urlCode}
		urls = append(urls, newUrl)
		c.IndentedJSON(http.StatusOK, urls)
		return
	default:
		panic(err)
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Url code not found!!"})
}

func main() {

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.POST("/url", postUrl)
	router.GET("/url/:id", getUrlByID)
	router.Run("localhost:8888")
}
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
