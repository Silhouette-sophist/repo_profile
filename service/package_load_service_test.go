package service

import (
	"context"
	"testing"
)

func TestLoadPackages(t *testing.T) {
	ctx := context.Background()
	repoPath := "/Users/silhouette/work-practice/repo_profile"
	pkgPath := "github.com/Silhouette-sophist/repo_profile"
	// LoadPackages(ctx, &LoadConfig{
	// 	RepoPath: repoPath,
	// 	PkgPath:  pkgPath,
	// 	LoadEnum: LoadAllPkg,
	// })
	LoadPackages(ctx, &LoadConfig{
		RepoPath: repoPath,
		PkgPath:  pkgPath,
		LoadEnum: LoadSpecificPkgWithChild,
	})
}

func TestLoadAllPackages(t *testing.T) {
	ctx := context.Background()
	repoPath := "/Users/silhouette/work-practice/repo_profile"
	pkgPath := "github.com/Silhouette-sophist/repo_profile"
	LoadPackages(ctx, &LoadConfig{
		RepoPath: repoPath,
		PkgPath:  pkgPath,
		LoadEnum: LoadCurrentRepo,
	})
}
