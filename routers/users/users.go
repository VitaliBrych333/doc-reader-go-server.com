package users

import (
    "fmt"
    "log"
    "net/http"
    "database/sql"
    "github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine, authenticateMiddleware gin.HandlerFunc) {
    // only for dev testing (it's not necessary)
    //---------------use HTML templates-----------------
    route.GET("/users", authenticateMiddleware, getUsersPage)

    route.GET("/newUser", authenticateMiddleware, func(context *gin.Context) {
        context.HTML(http.StatusOK, "form.html", gin.H{
            "Title": "New User",
        })
    })
    //--------------------------------------------------

    user := route.Group("user")
    {
        user.GET("/:id", getUserById)
        user.POST("", addUser)

        user.PUT("/:id", updateUserById)
        user.POST("/:id", updateUserById)

        user.DELETE("/delete/:id", authenticateMiddleware, deleteUserById)
        user.GET("/delete/:id", authenticateMiddleware, deleteUserById)
    }
}

func GetUsers(context *gin.Context) []User {
    db := context.MustGet("DB").(*sql.DB)
    rows, err := db.Query("select * from Users")

    if err != nil {
        panic(err)
    }

    users := []User{}
     
    for rows.Next() {
        user := User{}
        err := rows.Scan(&user.Id, &user.User_Id, &user.First_Name, &user.Last_Name, &user.Email, &user.Password, &user.Role, &user.Info)

        if err != nil {
            log.Fatalf("impossible to scan rows of query: %s", err)
            fmt.Println("error", err)
            continue
        }

        users = append(users, user)
    }

   return users
}

func getUsersPage(context *gin.Context) {
    users := GetUsers(context)
    context.HTML(http.StatusOK, "users.html", users)
}

func getUserById(context *gin.Context) {
    id := context.Param("id")
    db := context.MustGet("DB").(*sql.DB)

    form := sqlGetUserById(db, id)
    form.Title = "Edit"

    context.HTML(http.StatusOK, "form.html", form)
}

func AddUserInDB(context *gin.Context) {
    db := context.MustGet("DB").(*sql.DB)
    newUser := User{}

    if err := context.BindJSON(&newUser); err != nil {
        return
    }

    result := sqlAddUser(db, newUser)

    id, err := result.LastInsertId()
    if err != nil {
        fmt.Printf("Add User: %v", err)
    }

    newUser.Id = int(id)

    fmt.Println(result.LastInsertId())  // id added
    fmt.Println(result.RowsAffected())  // count affected rows

    context.JSON(http.StatusCreated, newUser)
}

func addUser(context *gin.Context) {
    AddUserInDB(context)
}

func updateUserById(context *gin.Context) {
    id := context.Param("id")
    db := context.MustGet("DB").(*sql.DB)

    firstName := context.PostForm("first_name")
    lastName := context.PostForm("last_name")
    email := context.PostForm("email")
    password := context.PostForm("password")
    role := context.PostForm("role")
    info := context.PostForm("info") 

    result := sqlUpdateUser(db, id, firstName, lastName, email, password, role, info)

    fmt.Println(result.LastInsertId())  // id updated
    fmt.Println(result.RowsAffected())  // count affected rows

    context.Redirect(http.StatusMovedPermanently, "/users")
}

func deleteUserById(context *gin.Context) {
    id := context.Param("id")
    db := context.MustGet("DB").(*sql.DB)

    result := sqlDeleteUser(db, id)

    fmt.Println(result.LastInsertId())  // id deleted
    fmt.Println(result.RowsAffected())  // count affected rows

    context.Redirect(http.StatusMovedPermanently, "/users")
}