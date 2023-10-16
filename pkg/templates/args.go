package templates

import "github.com/hrvadl/coursework_db/pkg/models"

type SignUpArgs struct {
	Success bool
	Error   string
}

type GenericErrorArgs struct {
	Error string
}

type DealsArgs struct {
	Logined    bool
	Securities []models.Security
	Deals      []models.Deal
	Role       string
}

type DealArgs struct {
	Logined    bool
	Deal       *models.Deal
	Securities []models.Security
	Role       string
	AmountHas  int
}

type ToastArgs struct {
	Message string
	Error   string
}

type DealListArgs struct {
	Deals []models.Deal
}

type ProfileArgs struct {
	User         *models.User
	Transactions []models.Transaction
	Logined      bool
}
