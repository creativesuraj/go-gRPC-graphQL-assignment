package models

type Book struct {
	ID      string `bson:"uuid"`
	Title   string
	Author  string
	Isbn    string `bson:"isbn"`
	Summary string
}

func (Book) Collection() string {
	return "books"
}

func (b Book) GetId() string {
	return b.ID
}

func (b Book) GetTitle() string {
	return b.Title
}

func (b Book) GetAuthor() string {
	return b.Author
}

func (b Book) GetSummary() string {
	return b.Summary
}
func (b Book) GetIsbn() string {
	return b.Isbn
}
