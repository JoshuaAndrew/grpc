package main

import (
	"fmt"
	"os"
	"sort"
	"net"
	"syscall"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/JoshuaAndrew/grpc/api"
	"github.com/JoshuaAndrew/grpc/service"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {
	app := cli.NewApp()
	app.Name = "ocid"
	app.Usage = "ocid server"
	app.Version = "0.0.1"

	//GLOBAL OPTIONS
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.BoolFlag{
			Name:  "tls",
			Usage: "Connection uses TLS if true, else plain TCP",
		},
		/*
                Alternate Names
                You can set alternate (or short) names for flags by providing a comma-delimited list for the Name
		*/
		cli.StringFlag{
			Name:  "socket, s",
			Usage: "socket path for containerd's GRPC server",
			Value: "/run/containerd/containerd.sock",
		},
		cli.StringFlag{
			Name:  "tls-private-key-file, key",
			Usage: "`FILE` containing x509 private key",
		},
		/*
		Placeholder Values
                Sometimes it's useful to specify a flag's value within the usage string itself.
                Such placeholders are indicated with back quotes
                Will result in help output like:
                --config FILE, -c FILE   Load configuration from FILE
                Note that only the first placeholder is used. Subsequent back-quoted words will be left as-is
		*/
		cli.StringFlag{
			Name:  "tls-cert-file, cert",
			Usage: "`FILE` containing x509 Certificate for HTTPS",
		},
		cli.StringFlag{
			Name:  "data, d",
			Usage: "`FILE` containing location and name",
		},
	}

	/* Flags Ordering
	  Flags for the application and commands are shown in the order they are defined.
	  However, it's possible to sort them from outside this library by using FlagsByName with sort
	*/
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Info("debug:", context.GlobalBool("debug"))
		}
		return nil
	}
	app.Action = func(context *cli.Context) error {
		//context.Args()  arguments
		//fmt.Printf("Hello %q", context.Args().Get(0))

		path := context.GlobalString("socket")
		if path == "" {
			return fmt.Errorf("--socket path cannot be empty")
		}
		//lis,_  := createUnixSocket(path)
		lis, _ := createTcpSocket(path)



		tls := context.GlobalBool("tls")
		var opts []grpc.ServerOption
		if tls {
			keyFile := context.GlobalString("tls-private-key-file")
			certFile := context.GlobalString("tls-cert-file")
			if keyFile == "" || certFile == "" {
				grpclog.Fatalf("Failed to load tls-private-key-file or tls-cert-file")
				logrus.Fatal("Failed to load tls-private-key-file or tls-cert-file")
				os.Exit(1)
			}
			credential, err := credentials.NewServerTLSFromFile(certFile, keyFile)
			if err != nil {
				grpclog.Fatalf("Failed to generate credentials %v", err)
			}
			opts = []grpc.ServerOption{grpc.Creds(credential)}
		}


		s := grpc.NewServer(opts...)
		greetingService, _ := service.NewGreetingService()
		data := context.GlobalString("data")
		routeService := service.NewRouteServer(data)

		//register service to grpc Server
		api.RegisterGreetingServiceServer(s, greetingService)
		api.RegisterRouteServer(s, routeService)

		if err := s.Serve(lis); err != nil {
			logrus.Fatal(err)
			return err
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "containerd: %s\n", err)
		os.Exit(1)
	}
}

func createUnixSocket(path string) (net.Listener, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0660); err != nil {
		return nil, err
	}
	if err := syscall.Unlink(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// The network net must be a stream-oriented network: "tcp", "tcp4", "tcp6", "unix" or "unixpacket"
	// For TCP and UDP, the syntax of laddr is "host:port", like "127.0.0.1:8080". If host is omitted, as in ":8080"
	return net.Listen("unix", path)
}

func createTcpSocket(path string) (net.Listener, error) {
	// The network net must be a stream-oriented network: "tcp", "tcp4", "tcp6", "unix" or "unixpacket"
	// For TCP and UDP, the syntax of laddr is "host:port", like "127.0.0.1:8080". If host is omitted, as in ":8080"
	return net.Listen("tcp", path)
}


