package test

// file.top.comment

// Book
type Book struct {
	// Book.Name
	Name string
	// Book.Author
	Author struct {
		// Book.Author.Name
		Name string
		// Book.Author.Age
		Age int
	}
}

// value.document
var ThreeBodyProblem Book

// Interface.IBook
type IBook interface {
	// IBook.Name()
	Name()
}

// file.bottom.comment
