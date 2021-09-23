package config

var gitCommitHash string

type Version struct {
	GitCommitHash string
}

func NewVersion() *Version {
	return &Version{
		GitCommitHash: gitCommitHash,
	}
}
