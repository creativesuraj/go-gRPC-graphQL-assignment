package resolvers

import (
	"context"

	"github.com/aeon/gql-server/models"
	pb "github.com/aeon/grpc-server/protos/book"
	"github.com/aeon/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (r *queryResolver) Books(ctx context.Context, where *models.BookFilters) ([]*models.Book, error) {

	conn, err := utils.NewGRPCConnection()

	if err != nil {
		logrus.Errorf("Error connecting to grpc server %s", err.Error())
	}
	defer conn.Close()

	client := pb.NewBookHandlersClient(conn)
	id := &wrapperspb.StringValue{}
	if where != nil && where.ID != nil {
		id = &wrapperspb.StringValue{
			Value: *where.ID,
		}
	}

	books, err := client.GetBooks(ctx, id)
	if err != nil {
		return nil, err
	}
	allBooks := make([]*models.Book, 0)
	for _, book := range books.Book {
		allBooks = append(allBooks, &models.Book{
			ID:      book.Id,
			Author:  book.Author,
			Title:   book.Title,
			Isbn:    book.Isbn,
			Summary: &book.Summary,
		})
	}
	return allBooks, nil
}

func (r *mutationResolver) CreateBook(ctx context.Context, input models.BookInput) (*models.Book, error) {
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		logrus.Errorf("Error connecting to grpc server %s", err.Error())
	}

	client := pb.NewBookHandlersClient(conn)
	bookRequest := &pb.BookRequest{
		Title:   input.Title,
		Author:  input.Author,
		Summary: *input.Summary,
		Isbn:    input.Isbn,
	}
	book, err := client.CreateBook(ctx, bookRequest)
	if err != nil {
		return nil, err
	}
	return &models.Book{
		ID:      book.Id,
		Author:  book.Author,
		Title:   book.Title,
		Summary: &book.Summary,
		Isbn:    book.Isbn,
	}, nil
}

func (r *mutationResolver) UpdateBook(ctx context.Context, id string, input models.BookInput) (*models.Book, error) {
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		logrus.Errorf("Error connecting to grpc server %s", err.Error())
	}

	client := pb.NewBookHandlersClient(conn)
	updateBookRequest := &pb.UpdateBookRequest{
		Id: id,
		BookRequest: &pb.BookRequest{
			Title:   input.Title,
			Author:  input.Author,
			Summary: *input.Summary,
			Isbn:    input.Isbn,
		},
	}

	book, err := client.UpdateBook(ctx, updateBookRequest)
	if err != nil {
		return nil, err
	}
	return &models.Book{
		ID:      book.Id,
		Author:  book.Author,
		Title:   book.Title,
		Summary: &book.Summary,
		Isbn:    book.Isbn,
	}, nil
}

func (r *mutationResolver) DeleteBook(ctx context.Context, id string) (bool, error) {
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		logrus.Errorf("Error connecting to grpc server %s", err.Error())
	}

	client := pb.NewBookHandlersClient(conn)

	boolValue, err := client.DeleteBook(ctx, &wrapperspb.StringValue{
		Value: id,
	})
	if err != nil {
		return false, err
	}
	return boolValue.Value, nil
}
