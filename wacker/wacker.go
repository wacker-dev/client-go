package wacker

import (
	"context"
	"fmt"
	"net"

	"github.com/mitchellh/go-homedir"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/wacker-dev/client-go/internal"
)

type (
	RunRequest     = internal.RunRequest
	ServeRequest   = internal.ServeRequest
	StopRequest    = internal.StopRequest
	RestartRequest = internal.RestartRequest
	DeleteRequest  = internal.DeleteRequest
	LogRequest     = internal.LogRequest

	ProgramResponse = internal.ProgramResponse
	Program         = internal.Program
	ListResponse    = internal.ListResponse
	LogResponse     = internal.LogResponse
)

const (
	PROGRAM_STATUS_RUNNING  uint32 = 0
	PROGRAM_STATUS_FINISHED uint32 = 1
	PROGRAM_STATUS_ERROR    uint32 = 2
	PROGRAM_STATUS_STOPPED  uint32 = 3

	PROGRAM_TYPE_WASI uint32 = 0
	PROGRAM_TYPE_HTTP uint32 = 1
)

type Options struct {
	ctx         context.Context
	sockPath    string
	dialOptions []grpc.DialOption
}

func WithContext(ctx context.Context) func(opts *Options) {
	return func(opts *Options) {
		opts.ctx = ctx
	}
}

func WithSockPath(path string) func(opts *Options) {
	return func(opts *Options) {
		opts.sockPath = path
	}
}

func WithGRPCDialOptions(dialOptions ...grpc.DialOption) func(opts *Options) {
	return func(opts *Options) {
		opts.dialOptions = dialOptions
	}
}

type Client struct {
	conn *grpc.ClientConn
	internal.WackerClient
}

func NewClient(options ...func(*Options)) (*Client, error) {
	opts := &Options{
		ctx: context.Background(),
	}
	for _, f := range options {
		f(opts)
	}

	if opts.sockPath == "" {
		p, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		opts.sockPath = fmt.Sprintf("%s/.wacker/wacker.sock", p)
	}

	opts.dialOptions = append(opts.dialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", addr)
		}),
		grpc.WithAuthority("wacker"),
	)
	conn, err := grpc.DialContext(opts.ctx, opts.sockPath, opts.dialOptions...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:         conn,
		WackerClient: internal.NewWackerClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
