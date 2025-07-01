package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"gopkg.in/yaml.v3"
	"doc-reader-go-server.com/routers/documents"
	"doc-reader-go-server.com/routers/users"
)

type Conf struct {
	DataBase struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Net      string `yaml:"net"`
		Addres   string `yaml:"addres"`
		DBName   string `yaml:"dbName"`
	}

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}

type User struct {
	Id         int    `json:"id"`
	User_Id    string `json:"userId"`
	First_Name string `json:"firstName"`
	Last_Name  string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	Info       string `json:"info"`
}

var c Conf

// Add a new global variable for the secret key
var secretKey = []byte("your-secret-key")
var loggedInUser string

func main() {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatal(err)
	}

	db := connectDB()
	defer db.Close()

	server := gin.Default()
	server.Use(setDB(db))
	server.Use(CORSMiddleware())

	// only for dev testing (it's not necessary)
	//---------------use HTML templates-----------------
	server.LoadHTMLGlob("templates/*")
	server.Static("/static", "./static")
	server.GET("/register", func(context *gin.Context) {
		context.HTML(http.StatusOK, "form.html", gin.H{
			"Title": "Register",
		})
	})
	//--------------------------------------------------

	server.POST("/login", handleLogin)
	server.GET("/logout", func(context *gin.Context) {
		loggedInUser = ""
		context.SetCookie("token", "", -1, "/", c.Server.Host, false, true)
		context.Redirect(http.StatusSeeOther, "/")
	})

	server.POST("/register", registerUser)

	documents.Routes(server, authenticateMiddleware)
	users.Routes(server, authenticateMiddleware)

	server.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/users")
	})

	server.Run(":" + c.Server.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func connectDB() *sql.DB {
	cfg := mysql.Config{
		User:                 c.DataBase.User,
		Passwd:               c.DataBase.Password,
		Net:                  c.DataBase.Net,
		Addr:                 c.DataBase.Addres,
		DBName:               c.DataBase.DBName,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		panic(err)
	}

	pingErr := db.Ping()

	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")

	return db
}

// middleware
func setDB(db *sql.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("DB", db)
		context.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := context.Request.Header.Get("Origin")

		context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Cache-Control, Content-Disposition")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "PUT, GET, POST, DELETE, OPTIONS")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}

		context.Next()
	}
}

func handleLogin(context *gin.Context) {
	logUser := User{}

	if err := context.BindJSON(&logUser); err != nil {
		return
	}

	allUsers := users.GetUsers(context)

	idx := slices.IndexFunc(allUsers, func(u users.User) bool { return u.Email == logUser.Email && u.Password == logUser.Password })

	if idx != -1 {
		tokenString, err := createToken(logUser.Email)

		if err != nil {
			context.String(http.StatusInternalServerError, "Error creating token")
			return
		}

		loggedInUser = logUser.Email

		context.SetSameSite(http.SameSiteNoneMode)                                    // for cross domain requests to get Cookie in Header - front (netlify.app), back - (leapcell.app)
		context.SetCookie("token", tokenString, 3600, "/", c.Server.Host, true, true) // c.Server.Host - domain - localhost or {your name}.leapcell.app

		context.JSON(http.StatusOK, gin.H{"userId": allUsers[idx].User_Id, "status": "Success"})

	} else {
		context.String(http.StatusUnauthorized, "Invalid credentials")
	}
}

func getRole(username string) string {
	if username == "senior" {
		return "senior"
	}
	return "employee"
}

func createToken(username string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "documents-reader-app",           // Issuer
		"aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time 1 hour
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

func authenticateMiddleware(context *gin.Context) {
	// Retrieve the token from the cookie
	tokenString, err := context.Cookie("token")
	if err != nil {
		fmt.Println("Token missing in cookie")
		context.String(http.StatusUnauthorized, "Token missing in cookie")
		context.Abort()
		return
	}

	// Verify the token
	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Printf("Token verification failed: %v\\n", err)

		context.String(http.StatusUnauthorized, "Token verification failed")
		context.Abort()
		return
	}

	// Print information about the verified token
	fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)

	// Continue with the next middleware or route handler
	context.Next()
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}

func registerUser(context *gin.Context) {
	users.AddUserInDB(context)
	context.Redirect(http.StatusMovedPermanently, "/login")
}
