package repository

func NewRepository(repoName string) UserRepoI {
	switch repoName {
	case "user":
		return &UserRepository{}
	default:
		return nil
	}
}
