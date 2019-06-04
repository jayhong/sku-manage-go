package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"sku-manage/config"
	"sku-manage/log"
	"sku-manage/model"
	"sku-manage/server"
	"sku-manage/util"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:  %s [options] ...\n", os.Args[0])
	flag.PrintDefaults()
}

var (
	// Flags
	helpShort  = flag.Bool("h", false, "Show usage text (same as --help).")
	helpLong   = flag.Bool("help", false, "Show usage text (same as -h).")
	serverIp   = flag.String("ip", "127.0.0.1", "the server ip")
	serverPort = flag.Int("p", 8500, "the server port")
	configFile = flag.String("c", "config.json", "config file")
)

func main() {

	flag.Usage = usage
	flag.Parse()
	if *helpShort || *helpLong {
		flag.Usage()
		return
	}
	rand.Seed(time.Now().UnixNano())

	pearCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	if err := config.SrvConfig.Load(*configFile); err != nil {
		logrus.Error(err.Error())
		panic(err.Error())
	}
	log.InitLog("sku-manage")
	logrus.Debugf("%v", config.SrvConfig)

	model.DBInit()

	accountService := NewAccountService(
		new(util.DecodeAndValidator),
		PubJWT("inspect"),
	)
	router := mux.NewRouter()

	accountService.RegisterRoutes(router, "/v1/inspect")

	handler := negroni.New(negroni.NewRecovery(),
		server.NewLogger(),
		server.NewUberRatelimit(config.SrvConfig.Server.ReteLimit),
		server.NewParseForm())
	handler.Use(pearCors)
	handler.UseHandler(router)

	srv := &http.Server{
		Handler:      handler,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}
	go server.StartHttpServer(srv, fmt.Sprintf(":%d", *serverPort), config.SrvConfig.Server.ConnLimit)

	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("./file"))))
	go http.ListenAndServe(":8080", nil)

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan // wait for SIGINT
	logrus.Info("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	logrus.Info("Server gracefully stopped")
}
