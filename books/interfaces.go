package books

type BookStorer interface {
	add(isbn int, title, author string) error
	remove(isbn int) (string, error)
	list(isbn int) (*Book, error)
	listavailible() ([]*Book, error)
}

type RecordStorer interface {
	borrow(int, int) (string, error)
	returnbook(int) (string, error)
}
