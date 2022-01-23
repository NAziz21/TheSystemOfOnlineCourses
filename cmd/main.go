package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/NAziz21/TheSystemOfOnlineCourses/cmd/app"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/managers"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/users"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/ocsdb"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pool, err := pgxpool.Connect(connectCtx, dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer pool.Close()

	mux := mux.NewRouter()
	userSvc := users.NewService(pool)
	managerSvc := managers.NewService(pool)
	server := app.NewServer(mux, userSvc, managerSvc)
	server.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}

	return srv.ListenAndServe()

}
