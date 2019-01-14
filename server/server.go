package server

import (
	"context"
	"gitex.labbs.com.br/labbsr0x/proxy/go-horse/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"

	"gitex.labbs.com.br/labbsr0x/proxy/go-horse/config"
	"gitex.labbs.com.br/labbsr0x/proxy/go-horse/handlers"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GoHorse GoHorse
func GoHorse() *iris.Application {

	logSetup()

	app := iris.New()
	app.Use(recover.New())
	app.Use(prometheus.GetMetrics().ServeHTTP)

	app.Get("/active-filters", handlers.ActiveFiltersHandler)
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	//TODO mapear rota para receber token ou nao
	authToken := app.Party("/token/{token:string}/")
	authToken.Post("/{version:string}/containers/{containerId:string}/attach", handlers.AttachHandler)
	authToken.Get("/{version:string}/containers/{id:string}/logs", handlers.LogsHandler).Name = "container-logs"
	authToken.Get("/{version:string}/services/{id:string}/logs", handlers.LogsHandler).Name = "service-logs"
	authToken.Post("/{version:string}/containers/{containerId:string}/wait", handlers.WaitHandler)
	authToken.Post("/{version:string}/exec/{execInstanceId:string}/start", handlers.ExecHandler)
	authToken.Get("/{version:string}/containers/{containerId:string}/stats", handlers.StatsHandler)
	authToken.Get("/{version:string}/events", handlers.EventsHandler)

	app.Post("/{version:string}/containers/{containerId:string}/attach", handlers.AttachHandler)
	app.Get("/{version:string}/containers/{id:string}/logs", handlers.LogsHandler).Name = "container-logs"
	app.Get("/{version:string}/services/{id:string}/logs", handlers.LogsHandler).Name = "service-logs"
	app.Post("/{version:string}/containers/{containerId:string}/wait", handlers.WaitHandler)
	app.Post("/{version:string}/exec/{execInstanceId:string}/start", handlers.ExecHandler)
	app.Get("/{version:string}/containers/{containerId:string}/stats", handlers.StatsHandler)
	app.Get("/{version:string}/events", handlers.EventsHandler)
	app.Any("*", handlers.ProxyHandler)

	app.Run(iris.Addr(config.Port), iris.WithoutStartupLog)
	return app
}

// Restart Restart
func Restart(app *iris.Application) *iris.Application {
	app.Shutdown(context.Background())
	return GoHorse()
}

func logSetup() {
	zerolog.SetGlobalLevel(config.LogLevel)
	if config.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
