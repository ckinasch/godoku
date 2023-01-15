package main

import (
	"net/http"

	"github.com/ckinasch/godoku/sudokuboard"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func serveAPI(board *sudokuboard.Board) {
	app := gin.Default()

	app.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := app.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "helloo",
			})
		})
	}

	api.GET("/board", boardHandler)

	app.Run(":3399")
}

func boardHandler(c *gin.Context) {

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"board": "not yet serving the board",
	})

}

func main() {

	board := sudokuboard.GetBoard()
	board.PrintBoard()

	serveAPI(&board)

}

func init() {

}
