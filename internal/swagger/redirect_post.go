package swagger

import "github.com/antonio-alexander/go-blog-cors/internal"

// swagger:route POST /redirect cors redirect_post
// And endpoint that will redirect to another endpoint.
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
//   200: RedirectPostResponseOK
//   500: RedirectPostResponseError

// swagger:response RedirectPostResponseOK
type RedirectPostResponseOK struct {
}

// swagger:response RedirectPostResponseError
type RedirectPostResponseError struct {
	// in:body
	Body internal.Error
}

// swagger:parameters cors redirect_post
type RedirectPostParams struct {
	// in: query
	Endpoint string `json:"endpoint"`
}
