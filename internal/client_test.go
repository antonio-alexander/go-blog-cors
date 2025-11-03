package internal_test

import (
	"fmt"
	"testing"

	"github.com/antonio-alexander/go-blog-cors/internal"
)

func TestXxx(t *testing.T) {
	// s, err := internal.ClientHelloWorld()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(s)

	username, password := "username", "password"
	s, err := internal.ClientAuthorize(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}
