package cmd

import (
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/etcd"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/file"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/github"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/gitlab"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/kubernetes"
	_ "github.com/PapaDanielVi/secret-shift/internal/provider/vault"
)
