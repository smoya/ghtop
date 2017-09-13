package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"

	"time"

	"flag"

	"github.com/google/go-github/github"
	"github.com/smoya/ghtop/pkg/contributor"
	"github.com/smoya/ghtop/pkg/logx"
	"github.com/smoya/ghtop/pkg/server"
)

var port = flag.Int("port", 8080, "-port=<port> sets the server's listening port. 8080 by default.")
var env = flag.String("env", "prod", "-env=<environment> specifies the environment. prod by default.")
var ghToken = flag.String("gh-token", "", "-gh-token=<token> sets the token for Github API.")
var cacheTTL = flag.Int("ttl", 300, "-ttl=<cache-ttl> sets the ttl in seconds for the repository cache.")
var authUser = flag.String("auth-user", "", "-auth-user=<username> sets the username for basic authentication.")
var authPassword = flag.String("auth-password", "", "-auth-password=<password> sets the password for basic authentication.")

func init() {
	flag.Parse()

	if *ghToken == "" {
		log.Fatal("Missing Github token. -gh-token=<token>")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger, err := logx.NewStdOut("application", "ghtop", *env, logx.LevelInfo)
	if err != nil {
		log.Fatal(err)
	}

	ensureInterruptionsStopApplication(cancel, logger)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: *ghToken,
		},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	contributorRepo := contributor.WithCache(contributor.NewGihubRepository(client), time.Duration(*cacheTTL)*time.Second)

	config := server.NewConfig(*port, *env, *authUser, *authPassword)
	s := server.NewServer(config, logger, contributorRepo)

	err = s.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func ensureInterruptionsStopApplication(cancelFunc context.CancelFunc, logger logx.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Info(fmt.Sprintf("Got signal %s. Stopping server...", s))
		cancelFunc()

		os.Exit(1)
		return
	}()
}
