package main

import (
	"github.com/ckinasch/godoku/sudokuboard"
	"github.com/gin-gonic/gin"
)

func serveAPI(board *sudokuboard.Board) {
	app := gin.Default()

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, board.Cells)
		// c.JSON(200, "Hello")
	})

	app.Run(":3399")
}

func main() {

	board := sudokuboard.GetBoard()
	board.PrintBoard()

	serveAPI(&board)

}

func init() {

}
