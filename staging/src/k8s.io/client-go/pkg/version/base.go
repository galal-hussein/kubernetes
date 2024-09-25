package version

const (
	// DefaultKubeBinaryVersion is the hard coded k8 binary version based on the latest K8s release.
	// It is supposed to be consistent with gitMajor and gitMinor, except for local tests, where gitMajor and gitMinor are "".
	// Should update for each minor release!
	DefaultKubeBinaryVersion = "1.31"
)

var (
	gitMajor = "1"
	gitMinor = "30"
	gitVersion   = "v1.30.4-k3s12"
	gitCommit    = "0856115a0464fe3572e1f7b6fe75c459d2fb4625"
	gitTreeState = "clean"
	buildDate = "2024-09-25T19:57:03Z"
)
