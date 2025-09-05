package gen

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-safetemp"
	giturl "github.com/kubescape/go-git-url"
)

func FormatGitRepoUrl(repo string) (string, giturl.IGitURL, error) {
	repoURL, err := giturl.NewGitURL(repo)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf("git@%s:%s/%s.git", repoURL.GetHostName(), repoURL.GetOwnerName(), repoURL.GetRepoName()),
		repoURL, nil
}

func CloneGitRepo(repo string) (string, io.Closer, error) {
	src, repoURL, err := FormatGitRepoUrl(repo)
	var closer io.Closer = NopCloser{}
	if err != nil {
		return "", closer, err
	}

	tmpDir, closer, err := safetemp.Dir(os.TempDir(), repoURL.GetRepoName())
	if err != nil {
		return "", closer, err
	}
	opts := []getter.ClientOption{}

	client := getter.Client{
		Src:     src,
		Dst:     tmpDir,
		Mode:    getter.ClientModeDir,
		Options: opts,
	}
	err = client.Get()
	if err != nil {
		return "", closer, err
	}
	return tmpDir, closer, nil
}
