package routes

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"html"
	"html/template"
	"net/http"
	"time"
	. "webpanel/components/provider"
)

type CookieUser struct {
	Username string
	Password string
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

/* src */
func Home(w http.ResponseWriter, r *http.Request)  {
	cookieUser := GetUserByCookie(w, r)

	if cookieUser == "none" {
		LoginPage(w, r)
		return
	}

	// todo GetUser MongoDB
	user := GetUser(cookieUser)

	t, _ := template.ParseFiles("./views/index.html")

	t.Execute(w, user)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if GetUserByCookie(w, r) != "none" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}
	t, _ := template.ParseFiles("./views/login.html")

	t.Execute(w, nil)
}



func LoginHandler(w http.ResponseWriter, r *http.Request)  {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	r.ParseForm()
	username := html.EscapeString(r.FormValue("username"))
	password := GetMD5Hash(html.EscapeString(r.FormValue("password")))

	var res CookieUser
	e := Coll.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&res)

	if e != nil {
		//w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	expectedPassword := res.Password

	if expectedPassword != password {
		//w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	sessionToken := uuid.New().String()

	_, err := Cache.Do("SETEX", "sessions/" + sessionToken, "10000", username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	cookie := http.Cookie{
		Name: "SESSION_ID",
		Value: sessionToken,
		Expires: time.Now().Add(1 * time.Hour),
		Path: "/",
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}


