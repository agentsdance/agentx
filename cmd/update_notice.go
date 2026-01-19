package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/agentsdance/agentx/internal/version"
)

const updateNoticeTimeout = 3 * time.Second

func maybeHandleUpdateNotice() {
	ctx, cancel := context.WithTimeout(context.Background(), updateNoticeTimeout)
	defer cancel()

	notice, available, err := version.CheckForUpdate(ctx)
	if err != nil || !available {
		return
	}

	if !isInteractiveTerminal() {
		fmt.Fprintf(os.Stderr, "new version %s, run %s to upgrade\n", notice.Latest, notice.Command)
		return
	}

	choice, err := promptUpdateChoice(notice)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new version %s, run %s to upgrade\n", notice.Latest, notice.Command)
		return
	}

	switch choice {
	case updateChoiceNow:
		_ = runUpgradeCommand(notice.Command)
	case updateChoiceSkipVersion:
		_ = version.SkipVersion(notice.Latest)
	}
}

func isInteractiveTerminal() bool {
	return isTerminal(os.Stdin) && isTerminal(os.Stdout)
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func displayVersion(value string) string {
	return strings.TrimPrefix(value, "v")
}

func runUpgradeCommand(command string) error {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

const (
	updateChoiceNow = iota
	updateChoiceSkip
	updateChoiceSkipVersion
)

type updatePromptModel struct {
	current string
	latest  string
	command string
	cursor  int
	choice  int
}

func newUpdatePromptModel(notice version.UpdateNotice) updatePromptModel {
	return updatePromptModel{
		current: displayVersion(version.Version),
		latest:  displayVersion(notice.Latest),
		command: notice.Command,
		cursor:  updateChoiceNow,
		choice:  -1,
	}
}

func (m updatePromptModel) Init() tea.Cmd {
	return nil
}

func (m updatePromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < updateChoiceSkipVersion {
				m.cursor++
			}
		case "enter":
			m.choice = m.cursor
			return m, tea.Quit
		case "1":
			m.choice = updateChoiceNow
			return m, tea.Quit
		case "2":
			m.choice = updateChoiceSkip
			return m, tea.Quit
		case "3":
			m.choice = updateChoiceSkipVersion
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.choice = updateChoiceSkip
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m updatePromptModel) View() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  ✨ Update available! %s -> %s\n\n", m.current, m.latest))
	b.WriteString(fmt.Sprintf("  Release notes: %s\n\n", version.ReleaseNotesURL()))

	choices := []string{
		fmt.Sprintf("1. Update now (runs `%s`)", m.command),
		"2. Skip",
		"3. Skip until next version",
	}
	for i, choice := range choices {
		cursor := " "
		if m.cursor == i {
			cursor = "›"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	b.WriteString("\nPress enter to continue")
	return b.String()
}

func promptUpdateChoice(notice version.UpdateNotice) (int, error) {
	model := newUpdatePromptModel(notice)
	program := tea.NewProgram(model)
	finalModel, err := program.Run()
	if err != nil {
		return updateChoiceSkip, err
	}

	final, ok := finalModel.(updatePromptModel)
	if !ok || final.choice < 0 {
		return updateChoiceSkip, nil
	}
	return final.choice, nil
}
