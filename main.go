package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(
			filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()

	// setup gomniauth
	gomniauth.SetSecurityKey("triathlete chat application")
	gomniauth.WithProviders(
		facebook.New("922065688652103", "cc72fb23014a133590e45c40e2979e10",
			"http://localhost:8080/auth/callback/facebook"),
		// github.New("key", "secret",
		// 	"http://localhost:8080/auth/callback/github"),
		google.New("684231278602-6f1cbk5tg4atdgukbmoj5h42n9fqao4c.apps.googleusercontent.com", "qCVkSPuN-_IuB_sI1m_3Az_x",
			"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/", http.RedirectHandler("/chat", http.StatusTemporaryRedirect))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	//start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
