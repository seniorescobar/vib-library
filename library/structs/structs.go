package structs

type Member struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Book struct {
	Id              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	Year            int    `json:"year"`
	AvailableAmount int    `json:"availableAmount"`
}
