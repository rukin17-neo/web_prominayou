package models

type PageData struct {
	Title       string
	Content     string
	CurrentUser *User
}

type ContactsData struct {
	PageData
	Address string
	Phone   string
	Email   string
}

type ReviewsData struct {
	PageData
	Reviews []Review
}
