package swagger

import (
	"github.com/antonio-alexander/go-blog-cors/internal"
)

// swagger:route POST /authorize cors authorize
// Approximates a basic authorization endpoint.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//       basic:
//
//     Schemes: http, https
//
// responses:
//   200: AuthorizePostResponseOK
//   401: AuthorizePostResponseNotAuthorized
//   500: AuthorizePostResponseError

// swagger:response AuthorizePostResponseOK
type AuthorizePostResponseOK struct {
	// in:body
	Body internal.Message
}

// swagger:response AuthorizePostResponseNotAuthorized
type AuthorizePostResponseNotAuthorized struct {
	// in:body
	Body internal.Error
}

// swagger:response AuthorizePostResponseError
type AuthorizePostResponseError struct {
	// in:body
	Body internal.Error
}

// swagger:parameters cors authorize
type AuthorizePostParams struct{}
