package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

var (
	Server        = "aanmeld.database.windows.net"
	Port          = 1433
	database      = "gegevens"
	User          = "Superthibo"
	Password      = "Superadmin!"
	db            *sql.DB
	driverName    = "sqlserver"
	serverConnStr string
)

func main() {
	router := gin.Default()

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", Server, User, Password, Port, database)
	db, err := sql.Open(driverName, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the database connection is successful
	err = checkDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	router.LoadHTMLGlob("templates/*")
	router.Static("/static/", "./static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/versturen", func(c *gin.Context) {
		Name := c.PostForm("name")
		Birthdate := c.PostForm("birthdate")
		Email := c.PostForm("email")
		Phone := c.PostForm("phone")
		log.Printf("Ontvangen gegevens: Naam=%s, Geboortedatum=%s, Email=%s, Telefoonnummer=%s", Name, Birthdate, Email, Phone)

		insertQuery := "INSERT INTO gegevens (Name, Birthdate, Email, Phone) VALUES (@p1, @p2, @p3, @p4)"

		_, err := db.Exec(insertQuery, Name, Birthdate, Email, Phone)

		if err != nil {
			log.Printf("Fout bij het uitvoeren van SQL-query: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Er is een fout opgetreden bij het opslaan van de gegevens"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Aanmeldgegevens succesvol verzonden"})
	})

	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			log.Fatal(err)
			return
		}
		router.Run(":" + strconv.Itoa(port))
	} else {
		router.Run(":8080")
	}
}

func checkDatabase(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	fmt.Println("Successfully connected to the database")
	return nil
}
