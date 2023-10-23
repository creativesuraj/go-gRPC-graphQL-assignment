package handlers

import (
	"context"

	"github.com/aeon/grpc-server/models"
	pb "github.com/aeon/grpc-server/protos/book"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var internalServerError = "An unexpected error occurred"

type BookHandlersServer struct {
	pb.UnimplementedBookHandlersServer
}

func (b *BookHandlersServer) GetBooks(ctx context.Context, id *wrapperspb.StringValue) (*pb.Books, error) {
	book := &models.Book{}
	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("Error connecting to db: %s", err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}

	defer dbConn.Client.Disconnect(ctx)

	filter := bson.M{}
	if id != nil && id.Value != "" {
		logrus.Infof("GetBooks - query to get books by id :%s", id.GetValue())
		filter = bson.M{"uuid": id.GetValue()}
	} else {
		logrus.Infof("GetBooks - query to get all books")
	}
	books := make([]*models.Book, 0)
	if err = dbConn.FindAllRecords(ctx, &books, &filter); err != nil {
		logrus.Errorf("GetBooks - Error getting books from the db: %s", err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}
	bookResponses := make([]*pb.BookResponse, 0)
	for _, book := range books {
		bookResponses = append(bookResponses, &pb.BookResponse{
			Id:      book.GetId(),
			Title:   book.GetTitle(),
			Author:  book.GetAuthor(),
			Summary: book.GetSummary(),
			Isbn:    book.GetIsbn(),
		})
	}
	return &pb.Books{
		Book: bookResponses,
	}, nil
}

func (b *BookHandlersServer) CreateBook(ctx context.Context, bookRequest *pb.BookRequest) (*pb.BookResponse, error) {

	if bookRequest == nil {
		logrus.Error("BookRequest message is nil")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid BookRequest Message")
	}

	book := &models.Book{}

	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("CreateBook - Error connecting to db: %s", err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}

	defer dbConn.Client.Disconnect(ctx)

	book = &models.Book{
		ID:      uuid.New().String(),
		Author:  bookRequest.GetAuthor(),
		Summary: bookRequest.GetSummary(),
		Title:   bookRequest.GetTitle(),
		Isbn:    bookRequest.GetIsbn(),
	}

	if err = dbConn.InsertRecord(ctx, book); err != nil {
		logrus.Errorf("CreateBook - Error inserting a book: %s", err.Error())
		return nil, err
	}

	logrus.Infof("Inserted new book, id: %s", book.GetId())
	return &pb.BookResponse{
		Id:      book.GetId(),
		Title:   book.GetTitle(),
		Author:  book.GetAuthor(),
		Summary: book.GetSummary(),
		Isbn:    book.GetIsbn(),
	}, nil
}

func (b *BookHandlersServer) UpdateBook(ctx context.Context, updateBookRequest *pb.UpdateBookRequest) (*pb.BookResponse, error) {
	if updateBookRequest == nil {
		logrus.Error("UpdateBookRequest message is nil")
		return nil, status.Errorf(codes.InvalidArgument, "UpdateBookRequest message is required")
	}
	if updateBookRequest.Id == "" {
		logrus.Error("UpdateBook - Book Id is required to update the book")
		return nil, status.Errorf(codes.InvalidArgument, "book id is required")
	}

	if updateBookRequest.BookRequest == nil {
		logrus.Error("UpdateBookRequest message is nil")
		return nil, status.Errorf(codes.InvalidArgument, "BookRequest message is required")
	}

	book := &models.Book{}

	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("UpdateBook - Error connecting to db: %s", err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}

	defer dbConn.Client.Disconnect(ctx)

	filter := bson.M{"uuid": updateBookRequest.GetId()}
	updateQuery := bson.D{{Key: "$set", Value: bson.M{
		"author":  updateBookRequest.BookRequest.GetAuthor(),
		"title":   updateBookRequest.BookRequest.GetTitle(),
		"isbn":    updateBookRequest.BookRequest.GetIsbn(),
		"summary": updateBookRequest.BookRequest.GetSummary(),
	}}}

	if err = dbConn.UpdateRecord(ctx, filter, updateQuery); err != nil {
		logrus.Errorf("UpdateBook - Error updating book with id: %s, %s", updateBookRequest.GetId(), err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}
	logrus.Infof("Updated book with id: %s", updateBookRequest.GetId())
	return &pb.BookResponse{
		Id:      updateBookRequest.GetId(),
		Title:   updateBookRequest.BookRequest.GetTitle(),
		Author:  updateBookRequest.BookRequest.GetAuthor(),
		Summary: updateBookRequest.BookRequest.GetSummary(),
		Isbn:    updateBookRequest.BookRequest.GetIsbn(),
	}, nil
}

func (b *BookHandlersServer) DeleteBook(ctx context.Context, id *wrapperspb.StringValue) (*wrapperspb.BoolValue, error) {
	if id == nil || id.Value == "" {
		logrus.Error("DeleteBook - id is not provided")
		return nil, status.Errorf(codes.InvalidArgument, "Book id is required")
	}

	book := &models.Book{}

	dbConn, err := models.NewDBConnection(ctx, book.Collection())
	if err != nil {
		logrus.Errorf("Error connecting to db: %s", err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}

	defer dbConn.Client.Disconnect(ctx)
	if err = dbConn.DeleteRecord(ctx, id.GetValue()); err != nil {
		logrus.Errorf("DeleteBook - Error deleting a book: %s", err.Error())
		return nil, status.Errorf(codes.Internal, internalServerError)
	}

	return &wrapperspb.BoolValue{
		Value: true,
	}, nil
}
