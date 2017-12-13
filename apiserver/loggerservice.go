package apiserver

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/nyaxt/otaru/logger"
	"github.com/nyaxt/otaru/logger/logbuf"
	"github.com/nyaxt/otaru/pb"
)

type loggerService struct {
	lbuf *logbuf.LogBuf
}

func (*loggerService) GetCategories(ctx context.Context, req *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	cs := logger.Registry().Categories()

	pcs := make([]*pb.LoggerCategory, 0, len(cs))
	for _, c := range cs {
		pcs = append(pcs, &pb.LoggerCategory{
			Category: c.Category,
			Level:    uint32(c.Level),
		})
	}

	return &pb.GetCategoriesResponse{Category: pcs}, nil
}

func (*loggerService) SetCategory(ctx context.Context, req *pb.SetCategoryRequest) (*pb.SetCategoryResponse, error) {
	c := logger.Registry().CategoryIfExist(req.Category)
	if c == nil {
		return nil, grpc.Errorf(codes.NotFound, "Specified category not found")
	}

	c.Level = logger.Level(req.Level)

	return &pb.SetCategoryResponse{}, nil
}

func (s *loggerService) QueryLogs(ctx context.Context, req *pb.QueryLogsRequest) (*pb.QueryLogsResponse, error) {
	es := s.lbuf.Query(int(req.MinId), req.Category, int(req.Limit))
	pes := make([]*pb.QueryLogsResponse_Entry, 0, len(es))
	for _, e := range es {
		pes = append(pes, &pb.QueryLogsResponse_Entry{
			Id:       uint32(e.Id),
			Log:      e.Log,
			Category: e.Category,
			Level:    uint32(e.Level),
			Time:     e.Time.Unix(),
			Location: e.Location,
		})
	}

	return &pb.QueryLogsResponse{Entry: pes}, nil
}

func (s *loggerService) GetLatestLogEntryId(ctx context.Context, req *pb.GetLatestLogEntryIdRequest) (*pb.GetLatestLogEntryIdResponse, error) {
	id := s.lbuf.LatestEntryId()
	return &pb.GetLatestLogEntryIdResponse{Id: uint32(id)}, nil
}

const MaxEntries = 10000

func InstallLoggerService() Option {
	lbuf := logbuf.New(MaxEntries)
	logger.Registry().AddOutput(lbuf)

	return func(o *options) {
		o.serviceRegistry = append(o.serviceRegistry, serviceRegistryEntry{
			registerServiceServer: func(s *grpc.Server) {
				pb.RegisterLoggerServiceServer(s, &loggerService{lbuf})
			},
			registerProxy: pb.RegisterLoggerServiceHandlerFromEndpoint,
		})
	}
}
