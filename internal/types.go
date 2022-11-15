package internal

//These variables are populated at build time
// to find where the variables are...use  go tool nm ./app | grep app
//REFERENCE: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
var (
	Version   string
	GitCommit string
	GitBranch string
)

//Message is used to transport strings in json format for a given message
type Message struct {
	//a message
	// example: Hello, World!
	Message string `json:"message"`
}

//Error is used to communciate an error in json format
type Error struct {
	// the status code of the error
	// example: 500
	StatusCode int `json:"status_code"`

	//an error
	// example: Unspecified error has occurred
	Error string `json:"error"`
}
