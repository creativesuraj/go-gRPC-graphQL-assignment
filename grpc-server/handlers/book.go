package handlers

import (
	"context"
	"errors"

	"github.com/aeon/grpc-server/models"
	pb "github.com/aeon/grpc-server/protos/book"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type BookHandlersServer struct {
	pb.UnimplementedBookHandlersServer
}

func (b *BookHandlersServer) GetBooks(ctx context.Context, id *wrapperspb.StringValue) (*pb.Books, error) {
	book := &models.Book{}
	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("Error connecting to db: %s", err.Error())
		return nil, err
	}

	defer dbConn.Client.Disconnect(ctx)

	filter := bson.M{}
	if id != nil && id.Value != "" {
		filter = bson.M{"uuid": id.GetValue()}
	}
	books := make([]*models.Book, 0)
	if err = dbConn.FindAllRecords(ctx, &books, &filter); err != nil {
		logrus.Errorf("GetBooks - Error getting books from the db: %s", err.Error())
		return nil, err
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
	book := &models.Book{}
	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("Error connecting to db: %s", err.Error())
		return nil, err
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
	book := &models.Book{}
	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("Error connecting to db: %s", err.Error())
		return nil, err
	}

	defer dbConn.Client.Disconnect(ctx)
	if err != nil {
		logrus.Error("UpdateBook - unable to parse string to uuid")
		return nil, errors.New("Unexpected error occcured")
	}
	filter := bson.M{"uuid": updateBookRequest.GetId()}
	updateQuery := bson.D{{Key: "$set", Value: bson.M{
		"author":  updateBookRequest.BookRequest.GetAuthor(),
		"title":   updateBookRequest.BookRequest.GetTitle(),
		"isbn":    updateBookRequest.BookRequest.GetIsbn(),
		"summary": updateBookRequest.BookRequest.GetSummary(),
	}}}

	if err = dbConn.UpdateRecord(ctx, filter, updateQuery); err != nil {
		logrus.Errorf("UpdateBook - Error updating book with id: %s, %s", updateBookRequest.GetId(), err.Error())
		return nil, err
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
	book := &models.Book{}
	dbConn, err := models.NewDBConnection(ctx, book.Collection())

	if err != nil {
		logrus.Errorf("Error connecting to db: %s", err.Error())
		return nil, err
	}

	defer dbConn.Client.Disconnect(ctx)
	if err = dbConn.DeleteRecord(ctx, id.GetValue()); err != nil {
		logrus.Errorf("DeleteBook - Error deleting a book: %s", err.Error())
		return nil, err
	}

	return &wrapperspb.BoolValue{
		Value: true,
	}, nil
}
