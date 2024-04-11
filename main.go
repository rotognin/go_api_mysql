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

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum) // Inserir um novo album no banco
	router.PATCH("/albums", updateAlbum) // Atualizar um álbum no banco
	
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

func getAlbumByID(c *gin.Context) {
	connectDB(c)

	var registro album
	id := c.Param("id")

	err := dbx.conn.QueryRow("SELECT * FROM albums WHERE alb_id_album = ?", id).Scan(&registro.ID, &registro.Title, &registro.Artist, &registro.Price, &registro.IDBanco)

	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Álbum não encontrado."})
		return
	}

	c.IndentedJSON(http.StatusOK, registro)
}

// Obter todos os registros cadastrados no banco
func getAlbums(c *gin.Context) {
	//c.IndentedJSON(http.StatusOK, albums)
	connectDB(c)

	var registros album
	var albums []album

	rows, err := dbx.conn.Query("SELECT * FROM albums")
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Não foi possível realizar a consulta"})
		return
	}

	for rows.Next() {
		err := rows.Scan(&registros.ID, &registros.Title, &registros.Artist, &registros.Price, &registros.IDBanco)
		if err != nil {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Não foi possível obter os resultados do banco"})
			return
		}

		albums = append(albums, registros)
	}

	c.IndentedJSON(http.StatusOK, albums)

}

func postAlbum(c *gin.Context) {
	insertOrUpdateAlbum(c, "insert")
}

func updateAlbum(c *gin.Context) {
	insertOrUpdateAlbum(c, "update")
}

func insertOrUpdateAlbum(c *gin.Context, action string) {
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
		if action == "insert" {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "ID do álbum já cadastrado."})
			return
		}
	}

	if err != nil && action == "update" {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Álbum não cadastrado. Não será atualizado"})
		return
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

	var sql string

	switch action {
		case "insert":
			sql = "INSERT INTO albums (alb_title, alb_artist, alb_price, alb_id_album) VALUES (?, ?, ?, ?)"
		case "update":
			sql = "UPDATE albums SET alb_title = ?, alb_artist = ?, alb_price = ? WHERE alb_id_album = ?"
	}

	// Ao chegar até aqui, está OK para inserir as informações no banco
	stmt, err := dbx.conn.Prepare(sql)
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Erro ao executar a operação", "error": erro})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Operação realizada com sucesso", "msg": msg})
}
