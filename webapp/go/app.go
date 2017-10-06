package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
	"time"

	"golang.org/x/net/netutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/takashabe/go-router"
	"github.com/takashabe/go-session"
	_ "github.com/takashabe/go-session/memory"
)

var (
	ErrUnregisteredUser = errors.New("unregistered user")
	ErrAuthentication   = errors.New("failed authentication")

	sess *session.Manager
)

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
	UserID    int
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
type UserContent struct {
	Myself     *UserModel
	User       *UserModel
	Tweets     []*Tweet
	Followable bool
}

// template content
type FollowingContent struct {
	FollowingList []*Following
}

// template content
type FollowersContent struct {
	UserList []*UserModel
}

// for FollowingContent table mapping struct
type Following struct {
	UserId    int
	FollowId  int
	UserName  string
	CreatedAt time.Time
}

var gdb *sql.DB

func setup() {
	db, err := sql.Open("mysql", "isucon@/isucon?parseTime=true")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(32)
	db.SetMaxIdleConns(32)
	gdb = db
}

func getDB() *sql.DB {
	return gdb
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) (*UserModel, error) {
	s, err := sess.SessionStart(w, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to session start")
	}
	id := s.Get("id")
	if id == nil {
		return nil, errors.New("Not found user in session")
	}

	user, err := getUser(id.(int))
	if err != nil {
		if errors.Cause(err) == ErrUnregisteredUser {
			s.Delete(id)
			sess.SessionDestroy(w, r)
		}
		return nil, err
	}

	return &user, nil
}

func getUser(id int) (UserModel, error) {
	db := getDB()

	user := UserModel{}
	stmt, err := db.Prepare("select id, name, email from user where id=?")
	if err != nil {
		return user, errors.Wrap(err, "failed to prepared statement")
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return user, errors.Errorf("%v: %v", ErrUnregisteredUser, err)
	}

	return user, nil
}

func authenticate(email, password string) (UserModel, error) {
	user := UserModel{}

	db := getDB()

	stmt, err := db.Prepare("select id from user where email=? and passhash=sha2(concat(salt, ?), 256)")
	if err != nil {
		return UserModel{}, errors.Wrap(err, "failed to prepared statement")
	}
	defer stmt.Close()

	err = stmt.QueryRow(email, password).Scan(&user.ID)
	if err != nil {
		return UserModel{}, errors.Wrap(err, "failed to query scan")
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

func getIndex(w http.ResponseWriter, r *http.Request) {
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	content := IndexContent{User: user}
	tweets := []*Tweet{}

	db := getDB()

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
		t := &Tweet{}
		err := rows.Scan(&t.ID, &t.UserID, &t.Content, &t.CreatedAt)
		checkErr(errors.Wrap(err, "failed to query scan"))
		tweets = append(tweets, t)
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

func getLogin(w http.ResponseWriter, r *http.Request) {
	content := LoginContent{Message: "Isutterへようこそ!!"}
	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/login.tmpl"))
	err := tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies login template"))
	}
}

func postLogin(w http.ResponseWriter, r *http.Request) {
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

	s, err := sess.SessionStart(w, r)
	if err != nil {
		authError(w)
		return
	}

	s.Set("id", user.ID)
	http.Redirect(w, r, "/", 302)
}

func getLogout(w http.ResponseWriter, r *http.Request) {
	sess.SessionDestroy(w, r)
	http.Redirect(w, r, "/login", 302)
	return
}

func getTweet(w http.ResponseWriter, r *http.Request) {
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

func postTweet(w http.ResponseWriter, r *http.Request) {
	// require login
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
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

	stmt, err := db.Prepare("INSERT INTO tweet (user_id, content) VALUES (?,?)")
	defer stmt.Close()
	checkErr(errors.Wrap(err, "failed to insert tweet prepared statement"))

	_, err = stmt.Exec(user.ID, content)
	checkErr(errors.Wrap(err, "failed to exec insert tweet"))

	http.Redirect(w, r, "/", 302)
}

func userHandler(w http.ResponseWriter, r *http.Request, userID int) {
	// require login
	myself, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	content := UserContent{Myself: myself}

	db := getDB()

	stmt, err := db.Prepare("SELECT t.id,  t.user_id,  u.name,  t.content,  t.created_at " +
		"FROM tweet as t JOIN user as u " +
		"WHERE t.user_id=u.id AND user_id = ? ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to prepared statement"))
		http.NotFound(w, r)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to prepared statement"))
		http.NotFound(w, r)
		return
	}
	defer rows.Close()

	tweets := []*Tweet{}
	for i := 0; rows.Next(); i++ {
		t := &Tweet{}
		err := rows.Scan(&t.ID, &t.UserID, &t.UserName, &t.Content, &t.CreatedAt)
		checkErr(errors.Wrap(err, "failed to query scan"))
		tweets = append(tweets, t)
	}
	content.Tweets = tweets

	targetUser, err := getUser(userID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	content.User = &targetUser
	content.Followable = followable(myself.ID, targetUser.ID)

	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/user.tmpl"))
	err = tmpl.Execute(w, content)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies tweet template"))
	}
}

func followable(srcID int, dstID int) bool {
	if srcID == dstID {
		return false
	}

	db := getDB()
	stmt, err := db.Prepare("select count(*) from follow where user_id=? and follow_id=?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	var cnt int
	err = stmt.QueryRow(srcID, dstID).Scan(&cnt)
	if err != nil {
		return false
	}
	return cnt == 0
}

func getFollowing(w http.ResponseWriter, r *http.Request) {
	// require login
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	db := getDB()

	followingStmt, err := db.Prepare("SELECT user_id, follow_id, created_at FROM follow WHERE user_id = ?")
	defer followingStmt.Close()
	checkErr(errors.Wrap(err, "failed to following prepared statement"))

	rows, err := followingStmt.Query(user.ID)
	checkErr(errors.Wrap(err, "failed to select following query"))
	defer rows.Close()

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

func getFollowers(w http.ResponseWriter, r *http.Request) {
	// require login
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	db := getDB()
	stmt, err := db.Prepare("SELECT id, name, created_at " +
		"FROM user WHERE id IN (SELECT user_id FROM follow WHERE follow_id=?)")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to prepared statement"))
		http.NotFound(w, r)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.ID)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to query"))
		http.NotFound(w, r)
		return
	}
	defer rows.Close()

	fc := FollowersContent{
		UserList: []*UserModel{},
	}
	for i := 0; rows.Next(); i++ {
		u := UserModel{}
		err := rows.Scan(&u.ID, &u.Name, &u.CreatedAt)
		checkErr(errors.Wrap(err, "failed to followers query scan"))
		fc.UserList = append(fc.UserList, &u)
	}

	tmpl := template.Must(template.ParseFiles("views/layout.tmpl", "views/followers.tmpl"))
	err = tmpl.Execute(w, fc)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to applies following template"))
	}
}

func postFollow(w http.ResponseWriter, r *http.Request, id int) {
	// require login
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}

	db := getDB()
	stmt, err := db.Prepare("INSERT INTO follow (user_id, follow_id) VALUES (?, ?)")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to prepared statement"))
		http.NotFound(w, r)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, id)
	checkErr(errors.Wrap(err, "failed to exec insert follow"))

	http.Redirect(w, r, "/", 302)
}

func getInitialize(w http.ResponseWriter, r *http.Request) {
	// impossible to deploy a single binary
	_, err := exec.Command(os.Getenv("SHELL"), "-c", "../tools/init.sh").Output()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to initialize"))
	}
	_, err = exec.Command(os.Getenv("SHELL"), "-c", "../tools/index.sh").Output()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to initialize"))
	}
}

func main() {
	setup()
	defer gdb.Close()

	r := router.NewRouter()
	r.Get("/", getIndex)
	r.Get("/login", getLogin)
	r.Get("/logout", getLogout)
	r.Get("/tweet", getTweet)
	r.Get("/user/:id", userHandler)
	r.Get("/following", getFollowing)
	r.Get("/followers", getFollowers)
	r.Get("/initialize", getInitialize)

	r.Post("/login", postLogin)
	r.Post("/tweet", postTweet)
	r.Post("/follow/:id", postFollow)

	r.ServeFile("/css/*filepath", http.Dir("static/css"))

	log.Println("running server...")
	r.PrintRoutes(os.Stdout)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	listener := netutil.LimitListener(l, 10)
	log.Fatal(http.Serve(listener, r))
}

func init() {
	manager, err := session.NewManager("memory", "gosess", 3600)
	if err != nil {
		panic(err.Error())
	}
	sess = manager
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
