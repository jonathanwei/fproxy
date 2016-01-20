package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

var beCommand = cli.Command{
	Name:      "backend",
	Aliases:   []string{"be"},
	Usage:     "Run the backend.",
	ArgsUsage: "[path/to/config/proto]",
	Description: strings.TrimSpace(`
The default path for the config proto is "~/.config/fproxy/backend.textproto".

Supported fields in the config are listed below.

server: this field is used to store all configuration for the RPC server that
we run. This includes the address it's located at, the certificate it presents
to clients (and the associated private key), and the CA it uses to verify
client certificates. A sample clause:
  server {
    addr: ":10000"
    cert_file: "/certs/backend/backend.pem"
    key_file: "/certs/backend/backend-key.pem"
    client_ca_file: "/certs/frontend/ca.pem"
  }

serve_path: this required field lets you set the root path to serve files from.
A sample clause:
  serve_path: "/home/foo/bar/"
`),
	Action: func(c *cli.Context) {
		runBe(defaultConfigPath(c, "backend.textproto"))
	},
}

func runBe(configPath string) {
	var config pb.BackendConfig
	readConfig(configPath, &config)

	if config.ServePath == "" {
		glog.Fatal("No serving path was given in config.serve_path.")
	}

	l, err := net.Listen("tcp", config.Server.Addr)
	if err != nil {
		glog.Fatalf("Failed to listen on %v: %v", config.Server.Addr, err)
	}
	defer l.Close()

	glog.Infof("Listening for requests on %v", l.Addr())

	s := grpc.NewServer(grpc.Creds(getBeTLS(&config)))
	pb.RegisterBackendServer(s, &backendServer{
		fs: http.Dir(config.ServePath),
	})
	err = s.Serve(l)
	if err != nil {
		glog.Fatalf("Failed to serve on %v: %v", config.Server.Addr, err)
	}
}

type backendServer struct {
	fs http.FileSystem
}

func (b *backendServer) Open(stream pb.Backend_OpenServer) error {
	// TODO: use this with ACLs.
	authCookie := GetAuthCookie(stream.Context())
	glog.Infof("Got auth cookie: %v", authCookie)

	firstReq, err := stream.Recv()
	if err != nil {
		glog.Warningf("Failed to read initial request: %v", err)
		return err
	}

	if firstReq.Action != pb.OpenRequest_INITIAL {
		err := grpc.Errorf(codes.InvalidArgument, "Unexpected initial action: %v", firstReq)
		glog.Warning(err)
		return err
	}

	path := filepath.Clean("/" + firstReq.Path)
	file, err := b.fs.Open(path)
	if err != nil {
		glog.Warningf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	finfo, err := file.Stat()
	if err != nil {
		glog.Warningf("Error stating file: %v", err)
		return err
	}

	var firstResp pb.OpenResponse
	firstResp.Info = fileInfoToProto(finfo)

	if finfo.IsDir() {
		chinfos, err := file.Readdir(-1)
		if err != nil {
			glog.Warningf("Error reading directory: %v", err)
			return err
		}

		for _, chinfo := range chinfos {
			firstResp.Child = append(firstResp.Child, fileInfoToProto(chinfo))
		}
	}

	err = stream.Send(&firstResp)
	if err != nil {
		glog.Warningf("Error sending first response: %v", err)
		return err
	}

	// If this is a directory, then we've already yielded all metainfo above;
	// return immediately.
	if finfo.IsDir() {
		return nil
	}

	// Otherwise, this is a file so we should wait for further requests for data.
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			glog.Warningf("Unexpected stream error: %v", err)
			return err
		}

		var resp pb.OpenResponse

		// Handle a seek request.
		switch req.Action {
		case pb.OpenRequest_SEEK:
			newOffset, err := file.Seek(req.SeekOffsetBytes, protoWhenceToOSWhence(req.SeekType))
			if err != nil {
				glog.Warningf("Got seek error: %v", err)
				return err
			}
			resp.NewOffset = newOffset
		case pb.OpenRequest_CONTENT:
			// Bound the request to 1MB. The conversion to int is safe as it's int32
			// to int, which will always fit.
			const maxSize = 32000
			size := int(req.RequestedBytes)
			if size > maxSize {
				size = maxSize
			}

			// Perform the read.
			buf := getBuf(int(size))
			n, err := file.Read(buf)

			if n > 0 && err == io.EOF {
				err = nil
			}

			if err == io.EOF {
				return grpc.Errorf(codes.OutOfRange, "EOF")
			}

			if err != nil {
				glog.Warningf("Failed to read file: %v", err)
				return err
			}
			resp.ResponseBytes = buf[:n]
		default:
			err := grpc.Errorf(codes.InvalidArgument, "Unknown request action: %v", req)
			glog.Warning(err)
			return err
		}

		err = stream.Send(&resp)
		if err != nil {
			glog.Warningf("Error sending response: %v", err)
			return err
		}
	}

	return nil
}

func protoWhenceToOSWhence(seekType pb.OpenRequest_SeekType) int {
	switch seekType {
	case pb.OpenRequest_CUR:
		return os.SEEK_CUR
	case pb.OpenRequest_START_OF_FILE:
		return os.SEEK_SET
	case pb.OpenRequest_END_OF_FILE:
		return os.SEEK_END
	default:
		panic(fmt.Sprintf("Got unknown seekType: %v", seekType))
	}
}

func fileInfoToProto(finfo os.FileInfo) *pb.FileInfo {
	return &pb.FileInfo{
		Name:             finfo.Name(),
		SizeBytes:        finfo.Size(),
		Mode:             uint32(finfo.Mode()),
		LastModTimeNanos: finfo.ModTime().UnixNano(),
	}
}

func getBuf(size int) []byte {
	// TODO: use a sync.Pool.
	return make([]byte, size)
}

func putBuf(b []byte) {
	// TODO: put into a sync.Pool.
}

func getBeTLS(config *pb.BackendConfig) credentials.TransportAuthenticator {
	server := config.Server

	cert, err := tls.LoadX509KeyPair(server.CertFile, server.KeyFile)
	if err != nil {
		glog.Fatalf("Couldn't load backend certificate and key: %v", err)
	}

	clientCAPemBytes, err := ioutil.ReadFile(server.ClientCaFile)
	if err != nil {
		glog.Fatalf("Couldn't read frontend CA certificate: %v", err)
	}
	clientCAs := x509.NewCertPool()
	ok := clientCAs.AppendCertsFromPEM(clientCAPemBytes)
	if !ok {
		glog.Fatal("Unable to append frontend CA certificate to client CA pool.")
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
		MinVersion:   tls.VersionTLS12,
	})
}
