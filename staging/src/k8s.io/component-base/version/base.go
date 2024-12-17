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
	gitVersion   = "v1.32.0-k3s3"
	gitCommit    = "8e970155a664406729918fbbb91c0339603366a7"
	gitTreeState = "clean"
	buildDate = "2024-12-17T20:44:19Z"
)
