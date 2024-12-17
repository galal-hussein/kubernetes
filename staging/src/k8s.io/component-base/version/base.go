package version

const (
	// DefaultKubeBinaryVersion is the hard coded k8 binary version based on the latest K8s release.
	// It is supposed to be consistent with gitMajor and gitMinor, except for local tests, where gitMajor and gitMinor are "".
	// Should update for each minor release!
	DefaultKubeBinaryVersion = "1.31"
)

var (
	gitMajor = "1"
	gitMinor = "32"
	gitVersion   = "v1.32.0-k3s4"
	gitCommit    = "8443c6a93b10e0437ebba2a5e79c79b2e8baeca7"
	gitTreeState = "clean"
	buildDate = "2024-12-17T22:46:21Z"
)
