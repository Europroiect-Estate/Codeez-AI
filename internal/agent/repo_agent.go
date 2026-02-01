package agent

// RepoAgent understands repo structure, conventions, build/test hints.
type RepoAgent struct {
	Cwd string
}

// Describe returns a short description of the repo (key files, languages, build).
func (r *RepoAgent) Describe() (string, error) {
	return "Repo: use repo.map() tool for key files and build hints.", nil
}
