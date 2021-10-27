package transport

import (
	"context"

	"github.com/facily-tech/go-scaffold/pkg/domains/quote"
	pb "github.com/facily-tech/proto-examples/go-scaffold/build/go/quote"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/google/uuid"
)

type grpcServer struct {
	findByID, upsert, delete grpctransport.Handler
	pb.UnimplementedQuoteServiceServer
}

func NewGRPCServer(svc quote.ServiceI) pb.QuoteServiceServer {
	return &grpcServer{
		findByID: grpctransport.NewServer(
			quote.FindByID(svc),
			decodeGRPCFindByIDRequest,
			encodeGRPCFindByIDResponse,
		),
		upsert: grpctransport.NewServer(
			quote.Upsert(svc),
			decodeGRPCUpsertRequest,
			encodeGRPCUpsertResponse,
		),
		delete: grpctransport.NewServer(
			quote.DeleteByID(svc),
			decodeGRPCDeleteRequest,
			encodeGRPCDeleteResponse,
		),
	}
}

func (s *grpcServer) FindByID(ctx context.Context, req *pb.FindRequest) (*pb.FindResponse, error) {
	_, rep, err := s.findByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	findResp := rep.(pb.FindResponse) //nolint:govet
	resp := &pb.FindResponse{Id: findResp.Id, Content: findResp.Content, Error: findResp.Error}

	return resp, nil
}

func (s *grpcServer) Upsert(ctx context.Context, req *pb.UpsertRequest) (*pb.UpsertResponse, error) {
	_, rep, err := s.upsert.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	jsonResp := rep.(quote.JSONResponse)
	resp := &pb.UpsertResponse{Id: jsonResp.ID.String(), Content: jsonResp.Content, Error: jsonResp.Error}

	return resp, nil
}

func (s *grpcServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.Error, error) {
	_, _, err := s.delete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.Error{}, nil
}

func decodeGRPCFindByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.FindRequest)
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return quote.FindByIDRequest{ID: id}, nil
}

func encodeGRPCFindByIDResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(quote.JSONResponse)
	return pb.FindResponse{Id: resp.ID.String(), Content: resp.Content}, nil
}

func decodeGRPCUpsertRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpsertRequest)

	var id *uuid.UUID
	if len(req.Id) > 0 {
		pid, err := uuid.Parse(req.Id)
		if err != nil {
			return nil, err
		}
		id = &pid
	}

	return quote.UpsertRequest{JSONRequest: quote.JSONRequest{ID: id, Content: req.Content}}, nil
}

func encodeGRPCUpsertResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	resp := grpcResp.(quote.JSONResponse)
	return quote.JSONResponse{ID: resp.ID, Content: resp.Content, Error: resp.Error}, nil
}

func decodeGRPCDeleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteRequest)
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, err
	}

	return quote.DeleteRequest{ID: id}, nil
}

func encodeGRPCDeleteResponse(_ context.Context, grpcResp interface{}) (interface{}, error) {
	return nil, nil
}
