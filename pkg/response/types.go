package response

import pkgError "github.com/jacoobjake/einvoice-api/pkg/error"

type JSONApiResponse struct {
	Success          bool                       `json:"success"`
	Code             int                        `json:"code,omitempty"`
	Message          string                     `json:"message,omitempty"`
	Data             any                        `json:"data,omitempty"`
	ValidationErrors []pkgError.ValidationError `json:"validation_errors,omitempty"`
}
