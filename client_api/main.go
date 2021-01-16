package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/FanFani4/ports/client_api/api"
	"github.com/FanFani4/ports/client_api/reader"
	"github.com/FanFani4/ports/client_api/sender"
	"github.com/FanFani4/ports/ports"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
)

func main() {
	log := logrus.New()
	ctx, cancel := context.WithCancel(context.Background())

	err := godotenv.Load()
	if err != nil {
		log.Info("no log file")
	}

	jsonPath := os.Getenv("PORTS_JSON_PATH")

	reader, err := reader.NewJSONReader(ctx, log, jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	portsResponse := reader.GetPorts()

	conn, err := grpc.Dial(os.Getenv("SERVER_ADDR"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	cli := ports.NewPortDomainServiceClient(conn)

	go runSync(log, portsResponse, cli)

	port := os.Getenv("PORT")
	if _, err = strconv.Atoi(port); err != nil {
		log.Fatal("PORT must be an int")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("HOST is required")
	}

	s := &fasthttp.Server{
		Logger:  log,
		Handler: api.NewAPI(log, cli).HandleFasthttp,
	}

	go func() {
		log.Info("Server running on: http://" + host + ":" + port)
		log.Info(s.ListenAndServe(host + ":" + port))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)

	<-c

	log.Info("sync shutdown")
	cancel()

	log.Info("Start server shutting down")

	err = s.Shutdown()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("\nShutting down the server...")
}

func runSync(log *logrus.Logger, portsResponse <-chan *ports.Port, cli ports.PortDomainServiceClient) {

	sender := sender.NewSender(context.Background(), log, cli)

	sender.Send(portsResponse)
}
