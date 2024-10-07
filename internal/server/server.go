package server

import (
	"html/template"
	"net/http"

	"github.com/cydave/staticfs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/cydave/gintemplate2/internal/assets"
	"github.com/cydave/gintemplate2/internal/config"
	"github.com/cydave/gintemplate2/internal/middlewares"
)

func configureStaticFS(r *gin.Engine) error {
	s := staticfs.New(assets.Static)
	handler := s.Serve("/static")

	alias := func(to string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Request.URL.Path = "/static" + to
			r.HandleContext(c)
		}
	}
	for _, a := range getRootAssets() {
		r.GET(a, alias(a))
		r.HEAD(a, alias(a))
	}

	// Non top-level assets are mapped as expected.
	r.GET("/static/*filepath", handler)
	r.HEAD("/static/*filepath", handler)
	return nil
}

func configureTemplating(r *gin.Engine) error {
	funcMaps := template.FuncMap{}
	templ := template.New("").Funcs(funcMaps)
	templ, err := templ.ParseFS(assets.Templates, "templates/*.tmpl")
	if err != nil {
		return err
	}
	r.SetHTMLTemplate(templ)
	return nil
}

func configureSessions(r *gin.Engine) error {
	cfg := config.Get()
	store := cookie.NewStore(config.GetSessionSecret())
	store.Options(sessions.Options{
		Path:     cfg.GetString("cookie.path"),
		Domain:   cfg.GetString("cookie.domain"),
		MaxAge:   cfg.GetInt("cookie.max_age"),
		Secure:   cfg.GetBool("cookie.secure"),
		HttpOnly: cfg.GetBool("cookie.http_only"),
		SameSite: http.SameSiteStrictMode,
	})
	r.Use(sessions.Sessions("session", store))
	return nil
}

func Init() (*gin.Engine, error) {
	cfg := config.Get()
	if env := cfg.GetString("environment"); env == "" || env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	if err := configureStaticFS(r); err != nil {
		return nil, err
	}
	if err := configureTemplating(r); err != nil {
		return nil, err
	}
	if err := configureSessions(r); err != nil {
		return nil, err
	}

	// Register controllers / routes here.
	r.GET("/", middlewares.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title": "Hello World",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title": "Login",
		})
	})

	return r, nil
}
