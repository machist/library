package mq

import (
	"context"
	"log"
	"os"

	"github.com/cenkalti/backoff/v4"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

type options struct {
	DataDirectory string
	DataClear     bool
	ServerOptions *server.Options
}

type Option func(*options)

func WithDataDirectory(dir string) Option {
	return func(o *options) {
		o.DataDirectory = dir
	}
}

func WithDataClear(clear bool) Option {
	return func(o *options) {
		o.DataClear = clear
	}
}

func WithServerOptions(serverOptions *server.Options) Option {
	return func(o *options) {
		o.ServerOptions = serverOptions
	}
}

type MQ struct {
	Server *server.Server
}

func New(ctx context.Context, opts ...Option) (*MQ, error) {
	options := &options{
		DataDirectory: "./data/mq",
	}

	for _, o := range opts {
		o(options)
	}

	if options.DataClear {
		if err := os.RemoveAll(options.DataDirectory); err != nil {
			return nil, err
		}
	}

	if options.ServerOptions == nil {
		options.ServerOptions = &server.Options{
			JetStream: true,
			StoreDir:  options.DataDirectory,
		}
	}

	ns, err := server.NewServer(options.ServerOptions)
	if err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		ns.Shutdown()
	}()

	ns.Start()

	return &MQ{
		Server: ns,
	}, nil
}

func (mq *MQ) Close() error {
	if mq.Server != nil && mq.Server.Running() {
		mq.Server.Shutdown()
	}
	return nil
}

func (mq *MQ) WaitForServer() {
	b := backoff.NewExponentialBackOff()

	for {
		d := b.NextBackOff()
		ready := mq.Server.ReadyForConnections(d)
		if ready {
			break
		}

		log.Printf("NATS server not ready, waited %s, retrying...", d)
	}
}

func (mq *MQ) Client() (*nats.Conn, error) {
	return nats.Connect(mq.Server.ClientURL())
}
