package users

type User struct {
    Id           int     `json:"id"`
    User_Id      string  `json:"userId"`
    First_Name   string  `json:"firstName"`
    Last_Name    string  `json:"lastName"`
    Email        string  `json:"email"`
    Password     string  `json:"password"`
    Role         string  `json:"role"`
    Info         string  `json:"info"`
}

type FormInfo struct {
    Title        string  `json:"title"`
    User
}