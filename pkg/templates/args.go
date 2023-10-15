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
	Securities []models.Security
	Deals      []models.Deal
	Role       string
}

type DealArgs struct {
	Deal *models.Deal
}

type ToastArgs struct {
	Message string
	Error   string
}

type DealListArgs struct {
	Deals []models.Deal
}
