package version

const (
	// DefaultKubeBinaryVersion is the hard coded k8 binary version based on the latest K8s release.
	// It is supposed to be consistent with gitMajor and gitMinor, except for local tests, where gitMajor and gitMinor are "".
	// Should update for each minor release!
	DefaultKubeBinaryVersion = "1.31"
)

var (
	gitMajor = "1"
	gitMinor = "31"
	gitVersion   = "v1.31.1-k3s4"
	gitCommit    = "e5e28c2d0895e6fcf5c938c7100b734dd91a3017"
	gitTreeState = "clean"
	buildDate = "2024-09-25T20:20:34Z"
)
