package main
import (
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
  "os"
  "time"
  "database/sql"
  "log"
  _ "github.com/go-sql-driver/mysql"
)

type db_creds struct{
  name string
  port string
  host string
  user string
  password string
}

var creds db_creds;
var db *sql.DB
var err error;

func setCredentials(){
  if creds.name != ""{
    return
  }
  err := godotenv.Load(".env")
  if err != nil {
    fmt.Println(err)
  }
   
  creds.name= os.Getenv("DB_NAME")
  creds.port= os.Getenv("DB_PORT")
  creds.host= os.Getenv("DB_HOST")
  creds.user= os.Getenv("DB_USER")
  creds.password= os.Getenv("DB_PASSWORD")
}

func InitialiseDBInstance(){
  if db != nil{
    return 
  }
  db, err = sql.Open("mysql", "root:root@(127.0.0.1:3306)/root?parseTime=true") 
  if err != nil{
    log.Fatal(err)
  }
}

func CreateDB(http.ResponseWriter, *http.Request){
      if db == nil{
        InitialiseDBInstance()
      }
      query := `
            CREATE TABLE IF NOT EXISTS users(
                id BINARY(16) PRIMARY KEY,
                username TEXT NOT NULL,
                email TEXT NOT NULL,
                created_at DATETIME
            );`
      fmt.Println(query)
      if _, err := db.Exec(query); err != nil{
        log.Fatal(err);
      }
}

func InsertUser(http.ResponseWriter, *http.Request){
        if db == nil{
        InitialiseDBInstance()
      }
      // insert a new user
      username := "johndoe"
      email := "test@yopmail.com"
      created_at := time.Now()
      fmt.Println(username, email)
      _, err := db.Exec(`
          INSERT INTO users (id, username, email, created_at) VALUES(UUID_TO_BIN(UUID()), ?, ?, ?)`, username, email, created_at)
      if err != nil{
        fmt.Println("ERR: ",err)
      }
      fmt.Println("Inserted")
    }

    func GetUser(http.ResponseWriter, *http.Request){ // Query a single user
       if db == nil{
        InitialiseDBInstance()
      }
        var (
            id       []byte 
            username  string
            email  string
            createdAt time.Time
        )
        query := "SELECT id, username, email, created_at FROM users"
        fmt.Println(query)
        if err := db.QueryRow(query).Scan(&id, &username, &email, &createdAt); err != nil{
            log.Fatal(err)
        }
        fmt.Println(id, username, email, createdAt)
    }

    func GetUsers(http.ResponseWriter, *http.Request){ // Query all users
        
      if db == nil{
        InitialiseDBInstance()
      }
      type user struct {
            id       []byte 
            username  string
            email  string
            createdAt time.Time
        }
        fmt.Println("query all users")
        rows, err := db.Query(`SELECT id, username, email, created_at FROM users`)
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()
        var users []user
        for rows.Next() {
            var u user
            err := rows.Scan(&u.id, &u.username, &u.email, &u.createdAt)
            if err != nil {
                log.Fatal(err)
            }
            users = append(users, u)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%#v", users)
    }

    func DeleteUser(http.ResponseWriter, *http.Request){
      if db == nil{
        InitialiseDBInstance()
      }
      fmt.Println("delete user")
        _, err := db.Exec(`DELETE FROM users`) 
        if err != nil {
            log.Fatal(err)
        }
    }

func main(){
  fmt.Println("Hello World")
  setCredentials()
  r := mux.NewRouter()
  r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
    fmt.Printf("Server is now running at %s\n", r.URL.Path)
    InitialiseDBInstance()
    if err := db.Ping(); err!=nil{
      fmt.Println(err)
    }
    fmt.Println("It works!")
  })
  r.HandleFunc("/db", CreateDB).Methods("POST")
  r.HandleFunc("/users", GetUsers).Methods("GET")
  r.HandleFunc("/user", GetUser).Methods("GET")
  r.HandleFunc("/user", InsertUser).Methods("POST")
  r.HandleFunc("/user", DeleteUser).Methods("DELETE")
  
  http.ListenAndServe(":8080", r)
}
