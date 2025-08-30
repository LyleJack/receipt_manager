package user

type AccountStatus struct {
	LoggedIn bool
	UserID   int
	UserName string
	Settings Settings
}

type Settings struct {
	SaveReceipts bool
}

// OAuth login
