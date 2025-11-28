package profile

import (
	"context"
	"path/filepath"

	"github.com/Silhouette-sophist/repo_profile/internal/service"
)

func AnalyzeRepoProfile(ctx context.Context, repoPath string) ([]*service.ModuleInfo, error) {
	curRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}
	modules, err := service.ParseRepo(curRepoPath)
	if err != nil {
		return nil, err
	}
	return modules, nil
}
