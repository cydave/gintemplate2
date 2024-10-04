package server

import (
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/cydave/gintemplate2/internal/assets"
	"github.com/cydave/gintemplate2/internal/config"
	"github.com/cydave/gintemplate2/internal/middlewares"
)

func configureStaticFS(r *gin.Engine) error {
	hdlr := ServeStaticFS("/static", assets.Static)
	r.GET("/static/*filepath", hdlr)
	r.HEAD("/static/*filepath", hdlr)
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
