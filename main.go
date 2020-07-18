package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/csrf"
)

const cookieName = "session-id"

var templates = template.Must(template.ParseGlob("templates/*.html"))

// var sess *Session

var store = NewMemoryStore()

func main() {
	CSRF := csrf.Protect(
		[]byte("---super-long-and-secret-key---"),
		csrf.HttpOnly(true),
		csrf.Secure(false), // change to true for production to enforce HTTPS
	)

	mux := http.NewServeMux()
	mux.Handle("/", index())
	mux.Handle("/login", login())
	mux.Handle("/logout", logout())
	mux.Handle("/content", content())

	// Finally, start the HTTP server on port 8080.
	// If anything goes wrong, the log.Fatal call will output
	// the error to the console and exit the application.
	log.Fatal(http.ListenAndServe(":8080", CSRF(mux)))
}

func index() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "index.html", nil)
	}
	return http.HandlerFunc(fn)
}

func login() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]interface{})

		switch r.Method {
		case http.MethodGet:
			data[csrf.TemplateTag] = csrf.TemplateField(r)
			templates.ExecuteTemplate(w, "login.html", data)

		case http.MethodPost:
			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := authenticate(username, password)
			if err != nil {
				errors := make(map[string]string)
				errors["username"] = "Invalid credentials"
				data["errors"] = errors
				templates.ExecuteTemplate(w, "login.html", data)
				return
			}

			sessionID, err := uuid.NewV4()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			sess := Session{}
			sess.Values = make(map[interface{}]interface{})
			sess.Values["user"] = user
			store.Set(sessionID.String(), sess)

			http.SetCookie(w, &http.Cookie{
				HttpOnly: true,
				MaxAge:   24 * 60 * 60,
				Name:     cookieName,
				Path:     "/",
				Value:    sessionID.String(),
			})

			http.Redirect(w, r, "/content", http.StatusSeeOther)

		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}
	return http.HandlerFunc(fn)
}

func logout() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value
		sess, err := store.Get(sessionID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   cookieName,
			MaxAge: -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	return http.HandlerFunc(fn)
}

func content() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value
		sess, err := store.Get(sessionID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, ok := sess.Values["user"]
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		templates.ExecuteTemplate(w, "content.html", user)
	}
	return http.HandlerFunc(fn)
}
