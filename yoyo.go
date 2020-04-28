package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const YOYO_VER string = "0.7"

func main() {
	yoyoGraceHttp()
}


func yoyoGraceHttp(){
	fmt.Println("+++ START YOYO GRACE HTTP +++")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Version : "+YOYO_VER)
		rHostname, _ := os.Hostname()
		fmt.Fprintln(w, "HostName : "+rHostname)
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()
	log.Println("server started")

	stopC := make(chan os.Signal,1)
	signal.Notify(stopC,syscall.SIGTERM,syscall.SIGINT,syscall.SIGKILL)

	sgn := <- stopC
	fmt.Println("signal: ", sgn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)


	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown Error : ",err.Error())
	}
	// 5초의 타임아웃으로 ctx.Done()을 캐치합니다.
	select {
	case <-ctx.Done():
		fmt.Println("timeout of 5 seconds.")
	}

	fmt.Println("Server exiting")
}


