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
	gitVersion   = "v1.32.0-k3s1"
	gitCommit    = "f06cc46f703aaa03b4a0e36645c0aaf7b79ac65f"
	gitTreeState = "clean"
	buildDate = "2024-12-12T00:01:17Z"
)
