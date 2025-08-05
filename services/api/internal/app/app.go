package app

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"net/http"
	"sync"
)

type App struct {
	config  *AppConfig
	router  *chi.Mux
	jwtAuth *jwtauth.JWTAuth

	dbOnce sync.Once
	db     *bun.DB
}

func (app *App) initJwtAuth() {
	app.jwtAuth = jwtauth.New("HS256", app.config.jwtSecret, nil)
}

func (app *App) init() {
	app.loadAppConfig()
	app.initJwtAuth()
	app.initRouter()
}

func (app *App) Router() *chi.Mux { return app.router }

func (app *App) Config() *AppConfig { return app.config }

func (app *App) JwtAuth() *jwtauth.JWTAuth { return app.jwtAuth }

func (app *App) Db() *bun.DB {
	app.dbOnce.Do(func() {
		log.Println("Connecting to database")

		pgConn := pgdriver.NewConnector(
			pgdriver.WithAddr(fmt.Sprintf("%s:%d", app.config.PostgresHost, app.config.PostgresPort)),
			pgdriver.WithUser(app.config.PostgresUser),
			pgdriver.WithPassword(app.config.PostgresPassword),
			pgdriver.WithDatabase(app.config.PostgresDb),
			pgdriver.WithInsecure(true),
		)
		pgDb := sql.OpenDB(pgConn)
		db := bun.NewDB(pgDb, pgdialect.New())

		app.db = db
	})

	return app.db
}

func (app *App) Serve(addRoutes func(*App)) {
	addRoutes(app)

	log.Printf("Starting API at %d", app.config.Port)

	serverErr := http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), app.router)

	if nil != serverErr {
		log.Fatal(serverErr)
	}
}

func Init() *App {
	app := &App{}

	app.init()

	return app
}
