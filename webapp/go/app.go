package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/takashabe/go-isucon-exercise/webapp/go/session"
	_ "github.com/takashabe/go-isucon-exercise/webapp/go/session/memory"
)

var sessionManager session.Manager

type IndexContent struct {
	User      *UserModel
	Following int
	Followers int
	Tweets    []*Tweet
}

type Tweet struct {
	ID        int
	UserId    int
	UserName  string
	Content   string
	CreatedAt time.Time
}

type LoginContent struct {
	Message string
}

type UserModel struct {
	ID        int
	Name      string
	Email     string
	Salt      string
	Passhash  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "isucon@/isucon?parseTime=true")
	checkErr(errors.Wrap(err, "failed to open database"))
	return db
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) (*UserModel, error) {
	s, err := sessionManager.SessionStart(w, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to session start")
	}
	id := s.Get("id")
	if id == nil {
		return nil, errors.New("Not found user in session")
	}

	user := UserModel{}
	db := getDB()
	defer db.Close()
	stmt, err := db.Prepare("select id, name, email from user where id=?")
	defer stmt.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepared statement")
	}
	err = stmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		s.Delete("id")
		sessionManager.SessionDestroy(w, r)
		authError(w)
		return nil, errors.Wrapf(err, "Unregistered User(request id: %d)", id)
	}

	return &user, nil
}

func authenticate(email, password string) (UserModel, error) {
	user := UserModel{}

	db := getDB()
	defer db.Close()

	stmt, err := db.Prepare("select id from user where email=? and passhash=sha2(concat(salt, ?), 256)")
	defer stmt.Close()
	if err != nil {
		return user, errors.Wrap(err, "failed to prepared statement")
	}

	err = stmt.QueryRow(email, password).Scan(&user.ID)
	if err != nil {
		return user, errors.Wrap(err, "failed to query scan")
	}
	return user, nil
}

func authError(w http.ResponseWriter) {
	content := LoginContent{Message: "ログインに失敗しました"}
	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/login.tmpl"))
	w.WriteHeader(401)
	err := tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies template2"))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if "/" != r.URL.Path {
		log.Println("skipped /")
		return
	}
	log.Println("Called /")

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	content := IndexContent{User: user}
	tweets := make([]*Tweet, 100)

	db := getDB()
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, user_id, content, created_at " +
		"FROM tweet " +
		"WHERE USER_ID IN (SELECT follow_id FROM follow WHERE user_id=?) OR user_id = ? " +
		"ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to prepared statement"))
		http.NotFound(w, r)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.ID, user.ID)
	checkErr(errors.Wrap(err, "failed to query"))
	defer rows.Close()
	for i := 0; rows.Next(); i++ {
		t := Tweet{}
		err := rows.Scan(&t.ID, &t.UserId, &t.Content, &t.CreatedAt)
		checkErr(errors.Wrap(err, "failed to query scan"))
		tweets[i] = &t
	}
	content.Tweets = tweets

	followStmt, err := db.Prepare("SELECT count(*) FROM follow WHERE user_id = ?")
	checkErr(errors.Wrap(err, "failed to prepared statement"))
	defer followStmt.Close()
	followStmt.QueryRow(user.ID).Scan(&content.Following)

	followerStmt, err := db.Prepare("SELECT count(*) FROM follow WHERE follow_id = ?")
	checkErr(errors.Wrap(err, "failed to prepared statement"))
	defer followerStmt.Close()
	followerStmt.QueryRow(user.ID).Scan(&content.Followers)

	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/index.tmpl"))
	// pp.Println(content)
	err = tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies template1"))
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called loginHandler")
	// login
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			authError(w)
			return
		}
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		if email == "" || password == "" {
			authError(w)
			return
		}
		user, err := authenticate(email, password)
		if err != nil {
			authError(w)
			return
		}
		s, err := sessionManager.SessionStart(w, r)
		if err != nil {
			authError(w)
			return
		}
		s.Set("id", user.ID)
		http.Redirect(w, r, "/", 302)
		log.Println("redirect to /")
		return
	}

	// view login page
	content := LoginContent{Message: "Isutterへようこそ!!"}
	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/login.tmpl"))
	err := tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies template"))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called /logout")
	sessionManager.SessionDestroy(w, r)
	http.Redirect(w, r, "/login", 302)
	log.Println("redirect to /login")
	return
}

func initializeHandler(w http.ResponseWriter, r *http.Request) {
	// impossible to deploy a single binary
	exec.Command(os.Getenv("SHELL"), "-c", "../tools/init.sh").Output()
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/initialize", initializeHandler)
	http.ListenAndServe(":8080", nil)
}

func init() {
	manager, err := session.NewManager("memory", "gosess", 3600)
	checkErr(errors.Wrap(err, "failed to create session manager"))
	sessionManager = *manager
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
