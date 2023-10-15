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
