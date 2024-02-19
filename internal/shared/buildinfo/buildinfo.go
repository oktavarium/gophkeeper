package buildinfo

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var ErrorGettingWorkDir = errors.New("error acquiring work dir")
var ErrorOpeningRepo = errors.New("error opening repo")
var ErrorGettingGitHead = errors.New("error getting git head")
var ErrorGettingTags = errors.New("error getting tags")
var ErrorReadingTags = errors.New("error reading tags")

const (
	masterBranchName = "master"
	developPrefix    = "develop build"
	releasePrefix    = "release build"
	unknownVersion   = "unknown version"
)

var Version string
var BuildDate string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(ErrorGettingWorkDir)
		return
	}

	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Println(ErrorOpeningRepo)
		return
	}

	h, err := repo.Head()
	if err != nil {
		log.Println(ErrorGettingGitHead)
		return
	}

	iter, err := repo.Tags()
	if err != nil {
		log.Println(ErrorGettingTags)
		return
	}

	var tag *plumbing.Reference
	for {
		ref, err := iter.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Println(ErrorReadingTags)
			return
		}
		tag = ref
		break
	}

	versionPrefix := releasePrefix
	if h.Name().Short() != masterBranchName {
		versionPrefix = developPrefix
	}

	versionSuffix := unknownVersion
	if tag != nil {
		versionSuffix = tag.Name().Short()
	}

	Version = fmt.Sprintf("%s - %s", versionPrefix, versionSuffix)
}
