package grpc

import (
	"context"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/diyliv/tages/internal/db"
	"github.com/diyliv/tages/internal/models"
	"github.com/diyliv/tages/pkg/file"
	storagepb "github.com/diyliv/tages/proto/storage"
)

type grpcservice struct {
	logger *zap.Logger
	db     db.Database
	files  []models.File
	mu     sync.Mutex
}

func NewgRPCService(logger *zap.Logger, db db.Database) *grpcservice {
	return &grpcservice{
		logger: logger,
		db:     db,
		files:  make([]models.File, 0),
	}
}

func (gs *grpcservice) Upload(stream storagepb.StorageService_UploadServer) error {
	fb := file.NewFileBuffer()

	firstReq, err := stream.Recv()
	if err != nil {
		gs.logger.Error("Error while receiving: " + err.Error())
		return err
	}

	if err := fb.Write(firstReq.Chunk); err != nil {
		gs.logger.Error("Error while writing to buffer: " + err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.files = append(gs.files, models.File{
		Name:       strings.ReplaceAll(firstReq.GetName(), "client_files/", ""),
		CreateadAt: time.Now().Local(),
		UpdatedAt:  time.Now().Local(),
	})

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				if err := gs.db.Store(strings.ReplaceAll(firstReq.GetName(), "client_files/", ""), fb); err != nil {
					gs.logger.Error("Error while storing file: " + err.Error())
					return status.Error(codes.Internal, "Error while storing your file: "+err.Error())
				}
				return stream.SendAndClose(&storagepb.UploadResp{Status: "success"})
			}
			return status.Error(codes.Internal, "something went wrong :(")
		}

		if err := fb.Write(req.GetChunk()); err != nil {
			gs.logger.Error("Error while writing to buffer: " + err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (gs *grpcservice) GetAllFiles(ctx context.Context, e *empty.Empty) (*storagepb.GetAllFilesResp, error) {
	var files []*storagepb.File

	gs.mu.Lock()
	defer gs.mu.Unlock()
	if len(gs.files) == 0 {
		return nil, status.Error(codes.NotFound, "there are no files at the moment")
	}

	for _, val := range gs.files {
		files = append(files, &storagepb.File{
			Name:      val.Name,
			CreatedAt: timestamppb.New(val.CreateadAt),
			UpdatedAt: timestamppb.New(val.UpdatedAt),
		})
	}

	return &storagepb.GetAllFilesResp{File: files}, status.Error(codes.OK, "")
}

func (gs *grpcservice) SendFiles(e *empty.Empty, stream storagepb.StorageService_SendFilesServer) error {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		gs.logger.Error("Error while reading files in current directory: " + err.Error())
		return status.Error(codes.Internal, "something went wrong :(")
	}

	for _, file := range files {
		read, err := ioutil.ReadFile(file.Name())
		if err != nil {
			gs.logger.Error("Error while reading file: " + err.Error())
			return status.Error(codes.Internal, err.Error())
		}
		if err := stream.Send(&storagepb.SendFilesResp{
			Name:  file.Name(),
			Chunk: read,
		}); err != nil {
			gs.logger.Error("Error while sending file: " + err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}
	return status.Error(codes.OK, "")
}
