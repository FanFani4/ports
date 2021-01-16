package server

import (
	"context"

	"github.com/FanFani4/ports/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewServer(db *mongo.Collection) *Server {
	return &Server{db: db}
}

type Server struct {
	ports.UnimplementedPortDomainServiceServer
	db *mongo.Collection
}

func (s *Server) Insert(ctx context.Context, port *ports.Port) (*ports.InsertResponse, error) {
	opts := options.Update().SetUpsert(true)
	_, err := s.db.UpdateOne(ctx, bson.M{"_id": port.Id}, bson.M{"$set": port}, opts)

	return &ports.InsertResponse{Success: err == nil}, err
}

func (s *Server) Get(ctx context.Context, args *ports.GetArgs) (*ports.Port, error) {
	var res *ports.Port

	err := s.db.FindOne(ctx, bson.M{"_id": args.Id}).Decode(&res)

	return res, err
}

func (s *Server) List(ctx context.Context, args *ports.ListArgs) (*ports.ListResponse, error) {
	opts := options.Find().SetSkip(args.Skip).SetLimit(args.Limit)

	resp := &ports.ListResponse{}

	count, err := s.db.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, err
	}

	resp.Count = count

	cursor, err := s.db.Find(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &resp.Ports)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
