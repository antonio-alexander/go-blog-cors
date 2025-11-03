package swagger

import "github.com/antonio-alexander/go-blog-cors/internal"

// swagger:route GET /redirect cors redirect_get
// And endpoint that will redirect to another endpoint.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
// responses:
//   200: RedirectGetResponseOK
//   500: RedirectGetResponseError

// swagger:response RedirectGetResponseOK
type RedirectGetResponseOK struct {
}

// swagger:response RedirectGetResponseError
type RedirectGetResponseError struct {
	// in:body
	Body internal.Error
}

// swagger:parameters cors redirect_get
type RedirectGetParams struct {
	// in: query
	Endpoint string `json:"endpoint"`
}
