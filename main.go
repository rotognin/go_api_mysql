// Seguindo o tutorial em https://go.dev/doc/tutorial/web-service-gin

package main

import (
	"net/http"
	"database/sql"
	"time"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// "album" representa os dados de um álbum gravado
type album struct {
	ID     string  `json:"id"`       // ID enviado pelo cliente
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
	IDBanco int64  `json:"id_banco"` // ID autoincremental do banco
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums) // Inserir um novo album no banco
	router.GET("/info", infoService) // Informações sobre o serviço
	router.GET("/inserir", inserirInfo) // Testar a inserção de informações no banco
	
	router.Run("localhost:8180")
}

func connectDB(c *gin.Context) {
	db, err := sql.Open("mysql", "root:euaquinanet@tcp(127.0.0.1:3308)/albums_go")
	if err != nil {
		var erro string
		erro = err.Error()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Não foi aberto a conexão com o banco", "error": erro})
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	dbx.conn = db
}

func inserirInfo(c *gin.Context) {
	connectDB(c)

	stmt, err := dbx.conn.Prepare("INSERT INTO albums (alb_title, alb_artist, alb_price) VALUES (?, ?, ?)")
	if err != nil {
		var erro string
		erro = err.Error()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Erro no prepare do banco", "error": erro})
	}

	_, err = stmt.Exec("Segundo", "Jorge Fruit", 10.00)
	if err != nil {
		var erro string
		erro = err.Error()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Erro ao executar inserção", "error": erro})
	}

	var info string
	info = "Dados inseridos com sucesso!"
	c.IndentedJSON(http.StatusOK, gin.H{"message": info})
}

func infoService(c *gin.Context) {
	var info string
	info = "Serviço rodando em Go por webservice"
	c.IndentedJSON(http.StatusOK, gin.H{"message": info})
}

// getAlbums responds with the list of all albums as JSON. (pode ser dado qualquer nome, não precisa ser esse)
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	connectDB(c)

	var newAlbum album
	var msg string
	msg = ""

	err := c.BindJSON(&newAlbum)
	if err != nil {
		msg += " Erro ao processar JSON passado."
	}

	// Validar o ID enviado
	var id_msg string = checarID(newAlbum.ID)
	if id_msg != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": id_msg})
		return
	}

	// Verificar se esse ID já existe no banco (alb_id_album)
	var albums album

	row := dbx.conn.QueryRow("SELECT * FROM albums WHERE alb_id_album = ?", newAlbum.ID)
	if err := row.Scan(&albums.IDBanco, &albums.Title, &albums.Artist, &albums.Price, &albums.ID); err == nil {
		msg += " ID do álbum já cadastrado."
	} 

	// Validar o Título do Álbum
	newAlbum.Title = regexp.MustCompile(`[^a-zA-Z0-9 ,._-]+`).ReplaceAllString(newAlbum.Title, "")

	var title_msg string = checarTitle(newAlbum.Title)
	if title_msg != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": title_msg})
		return
	}

	// Validar o Artista
	newAlbum.Artist = regexp.MustCompile(`[^a-zA-Z0-9 ,._-]+`).ReplaceAllString(newAlbum.Artist, "")

	var artist_msg string = checarArtist(newAlbum.Artist)
	if artist_msg != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": artist_msg})
		return
	}

	// Validar o preço
	var price_msg string = checarPrice(newAlbum.Price)
	if price_msg != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": price_msg})
		return
	}

	if msg != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": msg})
		return
	}

	// Ao chegar até aqui, está OK para inserir as informações no banco
	stmt, err := dbx.conn.Prepare("INSERT INTO albums (alb_title, alb_artist, alb_price, alb_id_album) VALUES (?, ?, ?, ?)")
	if err != nil {
		var erro string
		erro = err.Error()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Erro no prepare do banco", "error": erro})
		return
	}

	_, err = stmt.Exec(newAlbum.Title, newAlbum.Artist, newAlbum.Price, newAlbum.ID)
	if err != nil {
		var erro string
		erro = err.Error()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Erro ao executar inserção", "error": erro})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Operação realizada com sucesso", "msg": msg})
}

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
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found", "error": "Sem informações"})
}
