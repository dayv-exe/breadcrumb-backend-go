package utils

import "os"

type TableNames struct {
	Users         string
	Search        string
	Relationships string
}

func GetAllTableNames() *TableNames {
	users := os.Getenv("USERS_TABLE")
	if users == "" {
		panic("USERS_TABLE environment variable not set")
	}

	search := os.Getenv("SEARCH_TABLE")
	if search == "" {
		panic("SEARCH_TABLE environment variable not set")
	}

	relationships := os.Getenv("RELATIONSHIPS_TABLE")
	if relationships == "" {
		panic("USERS_TABLE environment variable not set")
	}

	return &TableNames{
		Users:         users,
		Search:        search,
		Relationships: relationships,
	}
}
