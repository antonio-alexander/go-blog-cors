package swagger

import "github.com/antonio-alexander/go-blog-cors/internal"

// swagger:route POST /proxy cors proxy_post
// And endpoint that will proxy to another endpoint.
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
//   200: ProxyPostResponseOK
//   500: ProxyPostResponseError

// swagger:response ProxyPostResponseOK
type ProxyPostResponseOK struct {
}

// swagger:response ProxyPostResponseError
type ProxyPostResponseError struct {
	// in:body
	Body internal.Error
}

// swagger:parameters cors proxy_post
type ProxyPostParams struct {
	// in: query
	Endpoint string `json:"endpoint"`
}
