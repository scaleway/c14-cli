package version

var (
	// VERSION represents the version of c14
	VERSION = "v0.2"
	// GITCOMMIT is overlaoded by the Makefile
	GITCOMMIT = "commit"
	// UserAgent represents the user-agent used for the API calls
	UserAgent = "c14/" + VERSION
)
