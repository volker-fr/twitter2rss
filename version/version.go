package version

//git2go "github.com/libgit2/git2go"

const version = "0.3.0"

// TODO: get it on compile time only...
//func getGitCommit() string {
//	repo, err := git2go.OpenRepository(".")
//	if err != nil {
//		log.Printf("Can't open repository = %+v", err)
//		return "unkown-commit-id"
//	}
//
//	revspec, err := repo.RevparseSingle("HEAD")
//	if err != nil {
//		log.Printf("revparse err = %+v", err)
//		return "unknown-commit-id"
//	}
//
//	return revspec.Id().String()
//}

func GetVersion() string {
	return version //+ " (" + getGitCommit() + ")"
}
