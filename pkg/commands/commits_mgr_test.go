package commands_test

import (
	"github.com/go-errors/errors"
	. "github.com/jesseduffield/lazygit/pkg/commands"
	. "github.com/jesseduffield/lazygit/pkg/commands/commandsfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommitsMgr", func() {
	var (
		commander  *FakeICommander
		config     *FakeIGitConfigMgr
		commitsMgr *CommitsMgr
		statusMgr  *FakeIStatusMgr
		mgrCtx     *MgrCtx
	)

	BeforeEach(func() {
		commander = NewFakeCommander()
		config = &FakeIGitConfigMgr{}
		config.ColorArgCalls(func() string { return "always" })

		mgrCtx = NewFakeMgrCtx(commander, config, nil)
		statusMgr = &FakeIStatusMgr{}

		commitsMgr = NewCommitsMgr(mgrCtx, statusMgr)
	})

	Describe("RewordHead", func() {
		It("runs expected command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git commit --allow-empty --amend --only -m \"newName\""),
			}, func() {
				commitsMgr.RewordHead("newName")
			})
		})
	})

	DescribeTable("CommitCmdObj",
		func(message string, flags string, expected string) {
			Expect(commitsMgr.CommitCmdObj(message, flags).ToString()).To(Equal(expected))
		},
		Entry(
			"with message",
			"my message",
			"",
			"git commit -m \"my message\"",
		),
		Entry(
			"with additional flags",
			"my message",
			"--flag",
			"git commit --flag -m \"my message\"",
		),
		Entry(
			"with multiline message",
			"line one\nline two",
			"--flag",
			"git commit --flag -m \"line one\" -m \"line two\"",
		),
	)

	Describe("GetHeadMessage", func() {
		It("runs expected command and trims output", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				{"git log -1 --pretty=%s", "blah blah\n", nil},
			}, func() {
				message, err := commitsMgr.GetHeadMessage()

				Expect(message).To(Equal("blah blah"))
				Expect(err).To(BeNil())
			})
		})

		It("returns error if one occurs", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				{"git log -1 --pretty=%s", "", errors.New("my error")},
			}, func() {
				message, err := commitsMgr.GetHeadMessage()

				Expect(message).To(Equal(""))
				Expect(err).To(MatchError("my error"))
			})
		})
	})

	Describe("GetMessageFirstLine", func() {
		It("returns first line", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				{"git show --no-patch --pretty=format:%s abc123", "firstline", nil},
			}, func() {
				message, err := commitsMgr.GetMessageFirstLine("abc123")

				Expect(message).To(Equal("firstline"))
				Expect(err).To(BeNil())
			})
		})

		It("bubbles up error", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				{"git show --no-patch --pretty=format:%s abc123", "", errors.New("my error")},
			}, func() {
				message, err := commitsMgr.GetMessageFirstLine("abc123")

				Expect(message).To(Equal(""))
				Expect(err).To(MatchError("my error"))
			})
		})
	})

	Describe("AmendHead", func() {
		It("runs command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git commit --amend --no-edit --allow-empty"),
			}, func() {
				err := commitsMgr.AmendHead()

				Expect(err).To(BeNil())
			})
		})
	})

	Describe("AmendHeadCmdObj", func() {
		It("returns command object", func() {
			obj := commitsMgr.AmendHeadCmdObj()
			Expect(obj.ToString()).To(Equal("git commit --amend --no-edit --allow-empty"))
		})
	})

	Describe("ShowCmdObj", func() {
		It("returns command object", func() {
			obj := commitsMgr.ShowCmdObj("abc123", "path")
			Expect(obj.ToString()).To(Equal("git show --submodule --color=always --no-renames --stat -p abc123 -- \"path\""))
		})

		It("handles lack of a path", func() {
			obj := commitsMgr.ShowCmdObj("abc123", "")
			Expect(obj.ToString()).To(Equal("git show --submodule --color=always --no-renames --stat -p abc123"))
		})
	})

	Describe("Revert", func() {
		It("runs command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git revert abc123"),
			}, func() {
				err := commitsMgr.Revert("abc123")
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("RevertMerge", func() {
		It("runs command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git revert abc123 -m 1"),
			}, func() {
				err := commitsMgr.RevertMerge("abc123", 1)
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("CreateFixupCommit", func() {
		It("runs command", func() {
			WithRunCalls(commander, []ExpectedRunCall{
				SuccessCall("git commit --fixup=abc123"),
			}, func() {
				err := commitsMgr.CreateFixupCommit("abc123")
				Expect(err).To(BeNil())
			})
		})
	})
})
