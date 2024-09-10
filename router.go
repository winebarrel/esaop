package esaop

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/itchyny/timefmt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	dfe "github.com/newm4n/go-dfe"
	"github.com/robfig/cron/v3"
	"github.com/winebarrel/esaop/esa"
	goth_esa "github.com/winebarrel/goth-esa/esa"
)

type ContextKey string

const (
	sessionName               = "_esaop_session"
	sessionUserKey            = "user"
	contextUserKey ContextKey = "user"
)

func init() {
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "esa", nil
	}
}

func NewRouter(cfg *Config) http.Handler {
	initGoth(cfg)
	esaCli := esa.NewClient(cfg.Team)
	store := newCookieStore(cfg.SessionSecret, cfg.CookieSecure)
	router := mux.NewRouter()
	router.Use(authorizeMiddleware(store))
	reph := regexp.MustCompile(`\${([^}]+)}`)
	dfetr := dfe.NewPatternTranslation()

	router.Path("/auth/callback").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user, err := gothic.CompleteUserAuth(rw, r)

		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(rw, err)
			return
		}

		sess, _ := store.Get(r, sessionName)
		sess.Values[sessionUserKey] = user
		sess.Save(r, rw)

		rw.Header().Set("Location", "/")
		rw.WriteHeader(http.StatusTemporaryRedirect)
	})

	router.PathPrefix("/").Methods("GET").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := getUser(r)
		token := user.AccessToken
		reqPath := strings.TrimPrefix(r.URL.Path, "/")

		reqPath = reph.ReplaceAllStringFunc(reqPath, func(s string) string {
			ph := strings.TrimSuffix(strings.TrimPrefix(s, "${"), "}")
			var dtFmt string
			var tm time.Time

			if strings.Contains(ph, "|") {
				cronDate := strings.SplitN(ph, "|", 2)
				cronExp := strings.TrimSpace(cronDate[0])
				cronExp = strings.ReplaceAll(cronExp, ",", " ")
				cronExp = strings.ReplaceAll(cronExp, ";", ",")
				dtFmt = strings.TrimSpace(cronDate[1])
				skd, err := cron.ParseStandard(cronExp)

				if err != nil {
					log.Printf("failed to parse cron expression: %s: %s", err, cronExp)
					return s
				}

				tm = skd.Next(time.Now())
			} else {
				dtFmt = strings.TrimSpace(ph)
				tm = time.Now()
			}

			return tm.Format(dfetr.JavaToGoFormat(dtFmt))
		})

		reqPath = timefmt.Format(time.Now(), reqPath) // DEPRECATED:
		cat, name := splitPath(reqPath)

		if name == "" {
			var loc string

			if cat == "" {
				loc = fmt.Sprintf("http://%s.esa.io", cfg.Team)
			} else {
				loc = categoryURL(cfg.Team, cat)
			}

			rw.Header().Set("Location", loc)
			rw.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		post, err := esaCli.Get(token, name, cat)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(rw, err)
			return
		}

		var loc string

		if post != nil {
			loc = post.URL
		} else {
			loc = newPostURL(cfg.Team, cat, name)
		}

		rw.Header().Set("Location", loc)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	})

	return handlers.LoggingHandler(os.Stdout, router)
}

func initGoth(cfg *Config) {
	gothic.Store = newCookieStore(cfg.SessionSecret, cfg.CookieSecure)
	callback, _ := url.Parse(cfg.Oauth2.RedirectHost)
	callback.Path = path.Join(callback.Path, "auth/callback")

	goth.UseProviders(
		goth_esa.New(cfg.Oauth2.ClientID, cfg.Oauth2.ClientSecret, callback.String(), "read"),
	)
}

func authorizeMiddleware(store *sessions.CookieStore) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth") ||
				strings.HasPrefix(r.URL.Path, "/logout") {
				next.ServeHTTP(rw, r)
				return
			}

			sess, _ := store.Get(r, sessionName)
			user := sess.Values[sessionUserKey]

			if user == nil {
				gothic.BeginAuthHandler(rw, r)
			} else {
				ctx := context.WithValue(r.Context(), contextUserKey, user)
				next.ServeHTTP(rw, r.WithContext(ctx))
			}
		})
	}
}

func newCookieStore(secret string, secure bool) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options.HttpOnly = true
	store.Options.MaxAge = 86400 * 30
	store.Options.Secure = secure
	return store
}

func getUser(r *http.Request) goth.User {
	return r.Context().Value(contextUserKey).(goth.User)
}

func splitPath(path string) (string, string) {
	if path == "" {
		return "", ""
	}

	if strings.HasSuffix(path, "/") {
		return path, ""
	}

	names := strings.Split(path, "/")
	return strings.Join(names[0:len(names)-1], "/"), names[len(names)-1]
}

func categoryURL(team string, category string) string {
	cat := strings.TrimSuffix(category, "/")

	if !strings.HasPrefix(cat, "/") {
		cat = "/" + cat
	}

	cat = url.QueryEscape(cat)
	return fmt.Sprintf("https://%s.esa.io/#path=%s", team, cat)
}

func newPostURL(team string, category string, name string) string {
	catName := path.Join(category, name)

	if !strings.HasPrefix(catName, "/") {
		catName = "/" + catName
	}

	catName = url.QueryEscape(catName)
	return fmt.Sprintf("https://%s.esa.io/posts/new?category_path=%s", team, catName)
}
