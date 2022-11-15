package swagger

import "github.com/antonio-alexander/go-blog-cors/internal"

// swagger:route GET /proxy cors proxy_get
// And endpoint that will proxy to another endpoint.
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
//   200: ProxyGetResponseOK
//   500: ProxyGetResponseError

// swagger:response ProxyGetResponseOK
type ProxyGetResponseOK struct {
}

// swagger:response ProxyGetResponseError
type ProxyGetResponseError struct {
	// in:body
	Body internal.Error
}

// swagger:parameters cors proxy_get
type ProxyGetParams struct {
	// in: query
	Endpoint string `json:"endpoint"`
}
