package server

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

func StartHttpServer(srv *http.Server, addr string, connlimit int) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	defer l.Close()
	l = LimitListener(l, connlimit)

	if err := srv.Serve(l); err != nil {
		logrus.Error(err.Error())
	}
}

func StartHttpsServer(handle http.Handler, addr, certFile, keyFile string) {
	if err := http.ListenAndServeTLS(addr, certFile, keyFile, handle); err != nil {
		logrus.Error(err.Error())
	}
}
