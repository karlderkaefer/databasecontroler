package main

import (
	"flag"
	"fmt"
	"github.com/karlderkaefer/databasemanager/database"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	flag.String("templates", "../config", "directory for *.tpl files contains sql scripts")
	port := flag.Int("port", 3000, "port server is listening to")
	flag.Parse()

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	router.Use(static.Serve("/static", static.LocalFile("./static", true)))

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run(fmt.Sprintf(":%d", *port))
}

func initializeRoutes() {
	router.GET("/", ShowIndexPage)
	router.GET("/api/list/:database", ListDatabaseApi)
	router.POST("/api/create", CreateDabaseApi)
	router.POST("/api/drop", DropDatabaseApi)
	router.POST("/api/recreate", RecreateDatabaseApi)
}

func ShowIndexPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":    "Database Manager",
			"database": "oracle11",
		},
	)
}

type FormData struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Database string `form:"database"`
}

func ListDatabaseApi(c *gin.Context) {
	databaseName := c.Param("database")
	db, err := database.GetDatabase(database.ParseVersion(databaseName))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	users, err := db.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func RecreateDatabaseApi(c *gin.Context) {
	var data FormData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := database.GetDatabase(database.ParseVersion(data.Database))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg, err := db.RecreateUser(data.Username, data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msg)
}

func CreateDabaseApi(c *gin.Context) {
	var data FormData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := database.GetDatabase(database.ParseVersion(data.Database))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg, err := db.CreateUser(data.Username, data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msg)
}

func DropDatabaseApi(c *gin.Context) {
	var data FormData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := database.GetDatabase(database.ParseVersion(data.Database))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg, err := db.DropUser(data.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msg)
}
