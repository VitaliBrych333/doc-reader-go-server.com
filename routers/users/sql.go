package users

import (
    "database/sql"
    "github.com/google/uuid"
)

func sqlGetUserById(db *sql.DB, id string) FormInfo {
    row := db.QueryRow("select * from Users where id = ?", id)

    form := FormInfo{}
    err := row.Scan(&form.Id, &form.User_Id, &form.First_Name, &form.Last_Name, &form.Email, &form.Password, &form.Role, &form.Info)

    if err != nil {
        panic("SQL couldn't get user by Id! " + err.Error())
    }

    return form
}

func sqlAddUser(db *sql.DB, newUser User) sql.Result {
    userId := "user-" + uuid.New().String()
    result, err := db.Exec("insert into Users (User_Id, First_Name, Last_Name, Email, Password, Role, Info) values (?, ?, ?, ?, ?, ?, ?)", userId, newUser.First_Name, newUser.Last_Name, newUser.Email, newUser.Password, newUser.Role, newUser.Info)

    if err != nil {
        panic("SQL couldn't add user! " + err.Error())
    }

    return result
}

func sqlUpdateUser(db *sql.DB, id string, firstName string, lastName string, email string, password string, role string, info string) sql.Result {
    result, err := db.Exec("update Users set First_Name = ?, Last_Name = ?, Email = ?, Password = ?, Role = ?, Info = ? where User_Id = ?",
        firstName, lastName, email, password, role, info, id)

    if err != nil {
        panic("SQL couldn't update user! " + err.Error())
    }

    return result
}

func sqlDeleteUser(db *sql.DB, id string) sql.Result {
    result, err := db.Exec("delete from Users where User_Id = ?", id)

    if err != nil {
        panic("SQL couldn't delete user! " + err.Error())
    }

    return result
}