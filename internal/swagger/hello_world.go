package swagger

import "github.com/antonio-alexander/go-blog-cors/internal"

// swagger:route GET / cors hello_world
// Returns "Hello, World!".
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
//   200: HelloWorldGetResponseOk

// swagger:response HelloWorldGetResponseOk
type HelloWorldGetResponseOk struct {
	// in:body
	Body internal.Message
}

// swagger:parameters hello_world
type HelloWorldGetParams struct{}
