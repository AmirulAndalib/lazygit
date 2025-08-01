package sync

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FetchAndAutoForwardBranchesAllBranches = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Fetch from remote and auto-forward branches with config set to 'allBranches'",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AutoForwardBranches = "allBranches"
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
		shell.NewBranch("feature")
		shell.NewBranch("diverged")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
		shell.SetBranchUpstream("feature", "origin/feature")
		shell.SetBranchUpstream("diverged", "origin/diverged")
		shell.Checkout("master")
		shell.HardReset("HEAD^")
		shell.Checkout("feature")
		shell.HardReset("HEAD~2")
		shell.Checkout("diverged")
		shell.HardReset("HEAD~2")
		shell.EmptyCommit("local")
		shell.NewBranch("checked-out")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("checked-out").IsSelected(),
				Contains("diverged ↓2↑1"),
				Contains("feature ↓2").DoesNotContain("↑"),
				Contains("master ↓1").DoesNotContain("↑"),
			)

		t.Views().Files().
			IsFocused().
			Press(keys.Files.Fetch)

		// AutoForwardBranches is "allBranches": both master and feature get forwarded
		t.Views().Branches().
			Lines(
				Contains("checked-out").IsSelected(),
				Contains("diverged ↓2↑1"),
				Contains("feature ✓"),
				Contains("master ✓"),
			)
	},
})
