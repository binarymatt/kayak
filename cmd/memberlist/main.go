package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/binarymatt/kayak/internal/log"
	"github.com/hashicorp/memberlist"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	app := &cli.App{
		Action: run,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name: "peers",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: 7946,
			},
			&cli.StringFlag{
				Name:  "addr",
				Value: "localhost:8080",
			},
			&cli.StringFlag{
				Name:     "name",
				Required: true,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("error starting memberlist", "error", err)
	}

}

type server struct {
	Name          string
	BindPort      int
	AdvertisePort int
	Addr          string
	list          *memberlist.Memberlist
}

func (s *server) Discover(ctx context.Context) error {

	conf := memberlist.DefaultLocalConfig()
	conf.Name = s.Name
	conf.BindPort = s.BindPort
	conf.AdvertisePort = s.AdvertisePort

	list, err := memberlist.Create(conf)
	if err != nil {
		return err
	}
	s.list = list
	s.list.LocalNode().Meta = []byte(s.Addr)

	tick := time.NewTicker(30 * time.Second)
	for {
		slog.Info("starting loop")
		select {
		case <-ctx.Done():
			slog.Info("context done")
			return nil
		case <-tick.C:
			slog.Debug("tick finished, processing peers")

			for _, member := range list.Members() {
				fmt.Printf("Member: %s %s:%d %s\n", member.Name, member.Addr, member.Port, string(member.Meta))
			}
		}
	}
}

func (s *server) join(w http.ResponseWriter, r *http.Request) {
	peer := r.URL.Query().Get("peer")
	slog.Info("join request", "peer", peer)
	_, err := s.list.Join([]string{peer})
	if err != nil {
		http.Error(w, "no peer", http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte("ok"))
}
func (s *server) Serve(ctx context.Context) error {
	slog.Info("setting up http server")
	mux := http.NewServeMux()
	mux.HandleFunc("/join", s.join)
	server := http.Server{
		Addr:    s.Addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelTimeout()

		if err := server.Shutdown(timeoutCtx); err != nil {
			slog.Error("error shutting down", "error", err)
		}
	}()
	slog.Info("http listening", "addr", s.Addr)
	return server.ListenAndServe()
}
func setupLogging() {
	logger := slog.New(log.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}
func run(cctx *cli.Context) error {
	ctx := context.TODO()
	setupLogging()
	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	defer cancel()
	s := &server{
		Name:          cctx.String("name"),
		BindPort:      cctx.Int("port"),
		AdvertisePort: cctx.Int("port"),
		Addr:          cctx.String("addr"),
	}
	g := new(errgroup.Group)
	g.Go(func() error { return s.Discover(ctx) })
	g.Go(func() error { return s.Serve(ctx) })
	return g.Wait()
}
