package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
	"golang.org/x/net/context"
)

type grpcFs struct {
	ctx    context.Context
	client pb.BackendClient
}

func (g grpcFs) Open(path string) (http.File, error) {
	// TODO: Add timeout to initial request.
	childCtx, childCancel := context.WithCancel(g.ctx)
	stream, err := g.client.Open(childCtx)
	if err != nil {
		glog.Warningf("Couldn't open client stream: %v", err)
		return nil, err
	}

	err = stream.Send(&pb.OpenRequest{
		Action: pb.OpenRequest_INITIAL,
		Path:   path,
	})
	if err != nil {
		glog.Warningf("Couldn't send initial request: %v", err)
		return nil, err
	}

	resp, err := stream.Recv()
	if err != nil {
		glog.Warningf("Couldn't recieve initial response: %v", err)
		return nil, err
	}

	reader := grpcFileReader{stream}
	return &grpcFile{
		ctx:        childCtx,
		cancelFunc: childCancel,
		stream:     stream,
		firstResp:  resp,

		reader:    reader,
		bufReader: bufio.NewReaderSize(reader, 32000),
	}, nil
}

type grpcFile struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	stream pb.Backend_OpenClient

	firstResp *pb.OpenResponse

	reader    grpcFileReader
	bufReader *bufio.Reader
}

func (g *grpcFile) Read(buf []byte) (int, error) {
	return g.bufReader.Read(buf)
}

func (g *grpcFile) Stat() (os.FileInfo, error) {
	return fileInfoFromProto{g.firstResp.Info}, nil
}

func (g *grpcFile) Readdir(count int) ([]os.FileInfo, error) {
	if len(g.firstResp.Child) == 0 {
		if count <= 0 {
			return nil, nil
		} else {
			return nil, io.EOF
		}
	}

	if count <= 0 || count > len(g.firstResp.Child) {
		count = len(g.firstResp.Child)
	}

	var ret []os.FileInfo
	for i := 0; i < count; i++ {
		ret = append(ret, fileInfoFromProto{g.firstResp.Child[i]})
	}
	g.firstResp.Child = g.firstResp.Child[count:]

	return ret, nil
}

func (g *grpcFile) Seek(offset int64, whence int) (int64, error) {
	err := g.stream.Send(&pb.OpenRequest{
		Action:          pb.OpenRequest_SEEK,
		SeekType:        osWhenceToProtoWhence(whence),
		SeekOffsetBytes: offset,
	})
	if err != nil {
		glog.Warningf("Got seek error: %v", err)
		return 0, err
	}

	resp, err := g.stream.Recv()
	if err != nil {
		glog.Warningf("Got seek error: %v", err)
		return 0, err
	}

	g.bufReader.Reset(g.reader)

	return resp.NewOffset, nil
}

func (g *grpcFile) Close() error {
	g.cancelFunc()
	return nil
}

type grpcFileReader struct {
	stream pb.Backend_OpenClient
}

func (g grpcFileReader) Read(buf []byte) (int, error) {
	err := g.stream.Send(&pb.OpenRequest{
		Action:         pb.OpenRequest_CONTENT,
		RequestedBytes: int32(len(buf)),
	})
	if err != nil {
		glog.Warningf("Got read error: %v", err)
		return 0, err
	}

	resp, err := g.stream.Recv()
	if err != nil {
		if grpc.Code(err) == codes.OutOfRange {
			return 0, io.EOF
		}

		glog.Warningf("Got read error: %v", err)
		return 0, err
	}

	got := resp.ResponseBytes
	if len(buf) < len(got) {
		got = got[:len(buf)]
	}

	copy(buf, got)
	return len(got), nil
}

type fileInfoFromProto struct {
	p *pb.FileInfo
}

func (f fileInfoFromProto) Mode() os.FileMode {
	return os.FileMode(f.p.Mode)
}

func (f fileInfoFromProto) IsDir() bool {
	return f.Mode().IsDir()
}

func (f fileInfoFromProto) ModTime() time.Time {
	return time.Unix(0, f.p.LastModTimeNanos)
}

func (f fileInfoFromProto) Name() string {
	return f.p.Name
}

func (f fileInfoFromProto) Size() int64 {
	return f.p.SizeBytes
}

func (f fileInfoFromProto) Sys() interface{} {
	return nil
}

func osWhenceToProtoWhence(whence int) pb.OpenRequest_SeekType {
	switch whence {
	case os.SEEK_CUR:
		return pb.OpenRequest_CUR
	case os.SEEK_SET:
		return pb.OpenRequest_START_OF_FILE
	case os.SEEK_END:
		return pb.OpenRequest_END_OF_FILE
	default:
		panic(fmt.Sprintf("Got unknown whence: %v", whence))
	}
}
