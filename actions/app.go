package actions

import (
	"strings"

	"github.com/gorilla/websocket"

	"github.com/gobuffalo/buffalo"
	popmw "github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	i18n "github.com/gobuffalo/mw-i18n"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/unrolled/secure"

	"github.com/Filip/blog/models"
	"github.com/gobuffalo/packr"

	"time"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var broadcast = make(chan Message)
var clients = make(map[*websocket.Conn]bool)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_blog_session",
		})
		// Automatically redirect to SSL
		app.Use(forceSSL())

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		app.Use(func(next buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				c.Set("year", time.Now().Year())
				c.Set("stripTag", strip.StripTags)
				c.Set("replace", strings.Replace)
				return next(c)
			}
		})

		app.GET("/", HomeHandler)
		app.GET("/home", VisitorHandler)
		app.GET("/article-detail/{article_id}", VisitorArticleShowHandler)
		app.POST("/async/search-article", VisitorSearchArticleHandler)
		app.POST("/async/create-message", CreateMessage)
		app.GET("/ws", MassageHandler)

		go handleMessages()

		app.Use(SetCurrentUser)
		app.Use(Authorize)
		app.GET("/users/new", UsersNew)
		app.POST("/users", UsersCreate)
		app.GET("/signin", AuthNew)
		app.POST("/signin", AuthCreate)
		app.DELETE("/signout", AuthDestroy)
		app.Middleware.Skip(Authorize, HomeHandler, VisitorHandler, VisitorArticleShowHandler, VisitorSearchArticleHandler, UsersNew, UsersCreate, AuthNew, AuthCreate)
		app.Resource("/articles", ArticlesResource{})
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
