package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/flosch/pongo2/v6"
	"net/http"
	"queue-manager/internal/providers"
	"queue-manager/internal/structures"
	"strings"
	"time"
)

type AdminLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminController struct {
	logger providers.Logger
	conf   *structures.Config
}

func (ac *AdminController) Index(w http.ResponseWriter, r *http.Request) {
	var tplMain = pongo2.Must(pongo2.FromFile(ac.conf.Template.TplDir + "dashboard.html"))
	str, err := tplMain.Execute(pongo2.Context{})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))

}

func (ac *AdminController) Static(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	parts := strings.Split(path, "/")
	fileName := parts[len(parts)-1]
	extParts := strings.Split(fileName, ".")
	ext := extParts[len(extParts)-1]
	fileNamePath := ""
	switch ext {
	case "css":
		fileNamePath = ac.conf.Template.PublicDir + "/css/" + fileName
		w.Header().Set("Content-Type", "text/css")
	case "js":
		fileNamePath = ac.conf.Template.PublicDir + "/js/" + fileName
		w.Header().Set("Content-Type", "text/javascript")
	case "png":
		fileNamePath = ac.conf.Template.PublicDir + "/img/" + fileName
		w.Header().Set("Content-Type", "image/png")
	case "jpg":
		fileNamePath = ac.conf.Template.PublicDir + "/img/" + fileName
		w.Header().Set("Content-Type", "image/jpeg")
	case "svg":
		fileNamePath = ac.conf.Template.PublicDir + "/img/" + fileName
		w.Header().Set("Content-Type", "image/svg+xml")
	case "json":
		fileNamePath = ac.conf.Template.PublicDir + "/css/" + fileName
		w.Header().Set("Content-Type", "application/json")

	}
	f, e := pongo2.FromFile(fileNamePath)
	if e != nil {
		fmt.Println(fileNamePath)
		w.WriteHeader(404)
		w.Write([]byte(e.Error()))
		return
	}

	var tplMain = pongo2.Must(f, e)
	str, err := tplMain.Execute(pongo2.Context{})

	if err != nil {
		fmt.Println(fileNamePath)
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))
}

func (ac *AdminController) Queues(w http.ResponseWriter, r *http.Request) {
	var tplMain = pongo2.Must(pongo2.FromFile(ac.conf.Template.TplDir + "queues.html"))
	str, err := tplMain.Execute(pongo2.Context{})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))
}

func (ac *AdminController) Queue(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("queue_name")

	var tplMain = pongo2.Must(pongo2.FromFile(ac.conf.Template.TplDir + "queue.html"))
	str, err := tplMain.Execute(pongo2.Context{"queue": name})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))
}

func (ac *AdminController) CreateQueue(w http.ResponseWriter, r *http.Request) {
	var tplMain = pongo2.Must(pongo2.FromFile(ac.conf.Template.TplDir + "add_queue_form.html"))
	str, err := tplMain.Execute(pongo2.Context{})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))
}

func (ac *AdminController) SendMessage(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("queue")
	var tplMain = pongo2.Must(pongo2.FromFile(ac.conf.Template.TplDir + "create_msg_form.html"))
	str, err := tplMain.Execute(pongo2.Context{"queue": name})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))
}

func (ac *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	var tplMain = pongo2.Must(pongo2.FromFile(ac.conf.Template.TplDir + "login.html"))
	str, err := tplMain.Execute(pongo2.Context{})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(str))
}

func (ac *AdminController) DoLogin(w http.ResponseWriter, r *http.Request) {
	var payload AdminLoginPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if payload.Username == ac.conf.DashboardAuthentication.Username && payload.Password == ac.conf.DashboardAuthentication.Password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": ac.conf.DashboardAuthentication.Username,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		secret := []byte(ac.conf.DashboardAuthentication.Secret)
		tokenString, err := token.SignedString(secret)
		if err != nil {
			ac.logger.Infof(providers.TypeApp, "Error creating token:", err)
			w.WriteHeader(401)
			return
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: true,
			Path:     "/",
		}

		http.SetCookie(w, &cookie)
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(401)
}

func (ac *AdminController) Logout(w http.ResponseWriter, r *http.Request) {
	deletedCookie := http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, &deletedCookie)
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}

func NewAdminController(logger providers.Logger, conf *structures.Config) *AdminController {
	return &AdminController{
		logger: logger,
		conf:   conf,
	}
}
