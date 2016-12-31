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

// DB table mapping
type UserModel struct {
	ID        int
	Name      string
	Email     string
	Salt      string
	Passhash  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DB table mapping
type Tweet struct {
	ID        int
	UserId    int
	UserName  string
	Content   string
	CreatedAt time.Time
}

// template content
type IndexContent struct {
	User      *UserModel
	Following int
	Followers int
	Tweets    []*Tweet
}

// template content
type LoginContent struct {
	Message string
}

// template content
type FollowingContent struct {
	FollowingList []*Following
}

// for FollowingContent table mapping struct
type Following struct {
	UserId    int
	FollowId  int
	UserName  string
	CreatedAt time.Time
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
		log.Println(errors.Wrap(err, "failed to applies login on authError template"))
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if "/" != r.URL.Path {
		return
	}

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
	err = tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies index template"))
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		return
	}

	// view login page
	content := LoginContent{Message: "Isutterへようこそ!!"}
	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/login.tmpl"))
	err := tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies login template"))
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionManager.SessionDestroy(w, r)
	http.Redirect(w, r, "/login", 302)
	return
}

func tweetHandler(w http.ResponseWriter, r *http.Request) {
	// POST
	if r.Method == "POST" {
		// require login
		user, err := getCurrentUser(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", 303)
			return
		}

		err = r.ParseForm()
		if err != nil {
			checkErr(errors.Wrap(err, "failed to parsed form on POST tweet"))
			http.NotFound(w, r)
			return
		}
		content := r.PostFormValue("content")
		if len(content) <= 0 {
			http.NotFound(w, r)
			return
		}

		db := getDB()
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO tweet (user_id, content) VALUES (?,?)")
		defer stmt.Close()
		checkErr(errors.Wrap(err, "failed to insert tweet prepared statement"))

		_, err = stmt.Exec(user.ID, content)
		checkErr(errors.Wrap(err, "failed to exec insert tweet"))

		http.Redirect(w, r, "/", 303)
		return
	}

	// require login
	_, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/tweet.tmpl"))
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies tweet template"))
	}
}

func followingHandler(w http.ResponseWriter, r *http.Request) {
	// require login
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	db := getDB()
	defer db.Close()

	followingStmt, err := db.Prepare("SELECT user_id, follow_id, created_at FROM follow WHERE user_id = ?")
	defer followingStmt.Close()
	checkErr(errors.Wrap(err, "failed to following prepared statement"))

	rows, err := followingStmt.Query(user.ID)
	checkErr(errors.Wrap(err, "failed to select following query"))

	fc := FollowingContent{
		FollowingList: make([]*Following, 0),
	}
	for i := 0; rows.Next(); i++ {
		f := Following{}

		// query from follow table
		err := rows.Scan(&f.UserId, &f.FollowId, &f.CreatedAt)
		checkErr(errors.Wrap(err, "failed to following query scan"))

		// query from user table
		err = db.QueryRow("SELECT name FROM user WHERE id = ?", f.FollowId).Scan(&f.UserName)
		checkErr(errors.Wrap(err, "failed to user query scan"))

		fc.FollowingList = append(fc.FollowingList, &f)
	}

	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/following.tmpl"))
	err = tmpl.Execute(w, fc)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies following template"))
	}
}

func initializeHandler(w http.ResponseWriter, r *http.Request) {
	// impossible to deploy a single binary
	exec.Command(os.Getenv("SHELL"), "-c", "../tools/init.sh").Output()
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler) // GET and POST
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/tweet", tweetHandler) // GET and POST
	// http.HandleFunc("/user", userHandler) // require user_id parameter -> "/user/101"
	http.HandleFunc("/following", followingHandler)
	// http.HandleFunc("/followers", followersHandler)
	// http.HandleFunc("/follow", followHandler) // POST. require user_id parameter -> "/follow/101"
	http.HandleFunc("/initialize", initializeHandler)

	log.Println("Started server...")
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
