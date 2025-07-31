package models

type PageData struct {
	Title   string
	Content string
}

type ContactsData struct {
	PageData
	Address string
	Phone   string
	Email   string
}
