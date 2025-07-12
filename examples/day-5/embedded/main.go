package main

import "fmt"

type Author struct {
	ID    uint
	name  string
	email string
}

func (author Author) Signature() string {
	return fmt.Sprintf("%s signed!", author.name)
}

type Blog struct {
	Author
	ID      uint
	Upvotes int
}

func main() {
	blog := Blog{
		ID:      1,
		Upvotes: 99,
		Author: Author{
			name:  "Peter",
			email: "peter@email.com",
		},
	}

	fmt.Println("Signature: ", blog.Signature())
}
