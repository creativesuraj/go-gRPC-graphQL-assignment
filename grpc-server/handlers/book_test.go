package handlers

import (
	"context"
	"testing"

	pb "github.com/aeon/grpc-server/protos/book"
	"github.com/aeon/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_GetBooks(t *testing.T) {
	var ctx context.Context
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		require.Error(t, err)
	}
	defer conn.Close()

	client := pb.NewBookHandlersClient(conn)
	id := &wrapperspb.StringValue{}
	allBooks, err := client.GetBooks(ctx, id)
	books := allBooks.Book
	require.Len(t, books, 1)
	require.NoError(t, err)
	require.Equal(t, books[0].Author, "Salinger")
	require.Equal(t, books[0].Title, "Catcher")
	require.Equal(t, books[0].Isbn, "978-0-316-76948-4")
	require.Equal(t, books[0].Summary, "No summary")
	allBooks1, err := client.GetBooks(ctx, &wrapperspb.StringValue{Value: "123"})
	books1 := allBooks1.Book
	require.Len(t, books1, 0)
	require.NoError(t, err)
}

func Test_CreateBook(t *testing.T) {
	var ctx context.Context
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		require.Error(t, err)
	}
	defer conn.Close()

	client := pb.NewBookHandlersClient(conn)
	_, err = client.CreateBook(ctx, nil)
	require.NoError(t, err)
	bookRequest := &pb.BookRequest{
		Author:  "J.D. Salinger",
		Title:   "The Catcher in the Rye",
		Summary: "no summary",
		Isbn:    "978-0-316-76948-4",
	}
	bookResponse, err := client.CreateBook(ctx, bookRequest)
	require.NoError(t, err)
	require.Equal(t, bookResponse.Author, bookRequest.Author)
	require.Equal(t, bookResponse.Title, bookRequest.Title)
	require.Equal(t, bookResponse.Isbn, bookRequest.Isbn)
	require.Equal(t, bookResponse.Summary, bookRequest.Summary)

}

func Test_UpdateBook(t *testing.T) {
	var ctx context.Context
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		require.Error(t, err)
	}
	defer conn.Close()

	client := pb.NewBookHandlersClient(conn)
	_, err = client.UpdateBook(ctx, nil)
	require.NoError(t, err)
	updateBookRequest := &pb.UpdateBookRequest{
		Id: "123",
		BookRequest: &pb.BookRequest{
			Author:  "J.D. Salinger",
			Title:   "The Catcher in the Rye",
			Summary: "no summary",
			Isbn:    "978-0-316-76948-4",
		}}
	bookResponse, err := client.UpdateBook(ctx, updateBookRequest)
	require.NoError(t, err)
	require.Equal(t, bookResponse.Author, updateBookRequest.BookRequest.Author)
	require.Equal(t, bookResponse.Title, updateBookRequest.BookRequest.Title)
	require.Equal(t, bookResponse.Isbn, updateBookRequest.BookRequest.Isbn)
	require.Equal(t, bookResponse.Summary, updateBookRequest.BookRequest.Summary)

}

func Test_DeleteBook(t *testing.T) {
	var ctx context.Context
	conn, err := utils.NewGRPCConnection()

	if err != nil {
		require.Error(t, err)
	}
	defer conn.Close()

	client := pb.NewBookHandlersClient(conn)
	_, err = client.DeleteBook(ctx, nil)
	require.NoError(t, err)
	isDeleted, err := client.DeleteBook(ctx, &wrapperspb.StringValue{
		Value: "123",
	})
	require.NoError(t, err)
	require.Equal(t, isDeleted, true)
}
