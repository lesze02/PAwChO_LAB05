package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type App struct {
	Port string
}

func index(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	version := os.Getenv("VERSION")
	fmt.Fprintf(w, "Hostname: "+hostname)
	fmt.Fprintf(w, "\nVersion: "+version)
	ifaces, err := net.Interfaces()
	if err != nil {
		println(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				fmt.Fprintf(w, "\nIP: "+ip.String())
				return
			}
		}
	}
}

func (a *App) Start() {
	http.Handle("/", logreq(index))
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func logreq(f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("path: %s", r.URL.Path)
		f(w, r)
	})
}

func main() {
	server := App{
		Port: env("PORT", "8080"),
	}
	server.Start()
}
