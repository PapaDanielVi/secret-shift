// Package tui provides an interactive terminal UI for configuring secret sync.
package tui

import (
	"errors"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/PapaDanielVi/secret-shift/internal/config"
)

type state int

const (
	stateSourceType state = iota
	stateSourceURL
	stateSourceRepo
	stateSourceToken
	stateSourceProjectID
	stateSourceVaultAddr
	stateSourceVaultPath
	stateSourceEtcdEndpoints
	stateSourceEtcdPrefix
	stateSourceKubeNamespace
	stateSourceKubeSecretName
	stateSourceKubeLabel
	stateSourceDone

	stateProcessPrefix
	stateProcessSuffix
	stateProcessIncludeRegex
	stateProcessExcludeRegex
	stateProcessDone

	stateDestType
	stateDestURL
	stateDestRepo
	stateDestToken
	stateDestProjectID
	stateDestFilePath
	stateDestFileFormat
	stateDestEncrypt
	stateDestVaultAddr
	stateDestVaultPath
	stateDestEtcdEndpoints
	stateDestEtcdPrefix
	stateDestKubeNamespace
	stateDestKubeSecretName
	stateDestConflictStrategy
	stateDestDone

	stateConfirm
	stateRun
	stateExit
)

var sourceTypes = []string{"github", "gitlab", "vault", "etcd", "kubernetes"}

var destTypes = []string{"file", "github", "gitlab", "vault", "etcd", "kubernetes"}

type model struct {
	state    state
	cfg      *config.Config
	inputs   []textinput.Model
	cursor   int
	choices  []string
	result   string
	quitting bool
}

func initialModel() model {
	inputs := make([]textinput.Model, 20)
	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].CharLimit = 256
	}

	m := model{
		state:   stateSourceType,
		cfg:     &config.Config{},
		inputs:  inputs,
		choices: sourceTypes,
		cursor:  0,
	}
	return m
}

func Run() (*config.Config, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("run tui: %w", err)
	}

	final := m.(model) //nolint:errcheck
	if final.quitting && final.result == "" {
		return nil, errors.New("TUI cancelled")
	}

	if err := final.cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return final.cfg, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.Key()

		if key.Code == 'c' && key.Mod == tea.ModCtrl {
			m.quitting = true
			return m, tea.Quit
		}

		if key.Code == tea.KeyEscape {
			m.quitting = true
			return m, tea.Quit
		}

		// Bridge states auto-advance without user interaction
		if m.isBridgeState() {
			return m.handleEnter()
		}

		if m.isInputState() {
			if key.Code == tea.KeyEnter {
				return m.handleEnter()
			}
			m.inputs[0], _ = m.inputs[0].Update(msg)
			return m, nil
		}

		// Selection states: navigate and confirm
		switch key.Code {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			return m.handleEnter()
		}
	}

	return m, nil
}

// isInputState returns true when the current state requires text input.
func (m model) isInputState() bool {
	switch m.state {
	case stateSourceURL, stateSourceRepo, stateSourceToken,
		stateSourceProjectID, stateSourceVaultAddr, stateSourceVaultPath,
		stateSourceEtcdEndpoints, stateSourceEtcdPrefix,
		stateSourceKubeNamespace, stateSourceKubeSecretName, stateSourceKubeLabel,
		stateProcessPrefix, stateProcessSuffix,
		stateProcessIncludeRegex, stateProcessExcludeRegex,
		stateDestURL, stateDestRepo, stateDestToken,
		stateDestProjectID, stateDestFilePath, stateDestFileFormat,
		stateDestEncrypt, stateDestVaultAddr, stateDestVaultPath,
		stateDestEtcdEndpoints, stateDestEtcdPrefix,
		stateDestKubeNamespace, stateDestKubeSecretName,
		stateDestConflictStrategy:
		return true
	default:
	}
	return false
}

// isBridgeState returns true for states that auto-advance to the next phase.
func (m model) isBridgeState() bool {
	switch m.state {
	case stateSourceDone, stateProcessDone, stateDestDone:
		return true
	default:
	}
	return false
}

func (m model) handleEnter() (tea.Model, tea.Cmd) { //nolint:gocyclo,cyclop,funlen
	switch m.state {
	case stateSourceType:
		m.cfg.Source.Type = sourceTypes[m.cursor]
		m.cursor = 0
		m.state = m.nextSourceState()
		if m.isInputState() {
			m.inputs[0].Focus()
		}

	case stateSourceURL:
		m.cfg.Source.URL = m.inputs[0].Value()
		m.resetInput()
		m.state = m.nextSourceState()
		if m.isInputState() {
			m.inputs[0].Focus()
		}

	case stateSourceRepo:
		m.cfg.Source.Repo = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceToken
		m.inputs[0].Focus()

	case stateSourceToken:
		m.cfg.Source.Token = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceDone

	case stateSourceProjectID:
		m.cfg.Source.ProjectID = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceToken
		m.inputs[0].Focus()

	case stateSourceVaultAddr:
		m.cfg.Source.VaultAddress = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceVaultPath
		m.inputs[0].Focus()

	case stateSourceVaultPath:
		m.cfg.Source.VaultPath = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceToken
		m.inputs[0].Focus()

	case stateSourceEtcdEndpoints:
		endpoints := strings.Split(m.inputs[0].Value(), ",")
		for i := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoints[i])
		}
		m.cfg.Source.EtcdEndpoints = endpoints
		m.resetInput()
		m.state = stateSourceEtcdPrefix
		m.inputs[0].Focus()

	case stateSourceEtcdPrefix:
		m.cfg.Source.EtcdPrefix = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceDone

	case stateSourceKubeNamespace:
		m.cfg.Source.KubeNamespace = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceKubeSecretName
		m.inputs[0].Focus()

	case stateSourceKubeSecretName:
		m.cfg.Source.KubeSecretName = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceKubeLabel
		m.inputs[0].Focus()

	case stateSourceKubeLabel:
		m.cfg.Source.KubeLabel = m.inputs[0].Value()
		m.resetInput()
		m.state = stateSourceDone

	case stateSourceDone:
		m.state = stateProcessPrefix
		m.cursor = 0
		m.inputs[0].Focus()

	case stateProcessPrefix:
		m.cfg.Process.AddPrefix = m.inputs[0].Value()
		m.resetInput()
		m.state = stateProcessSuffix
		m.inputs[0].Focus()

	case stateProcessSuffix:
		m.cfg.Process.AddSuffix = m.inputs[0].Value()
		m.resetInput()
		m.state = stateProcessIncludeRegex
		m.inputs[0].Focus()

	case stateProcessIncludeRegex:
		m.cfg.Process.IncludeRegex = m.inputs[0].Value()
		m.resetInput()
		m.state = stateProcessExcludeRegex
		m.inputs[0].Focus()

	case stateProcessExcludeRegex:
		m.cfg.Process.ExcludeRegex = m.inputs[0].Value()
		m.resetInput()
		m.state = stateProcessDone

	case stateProcessDone:
		m.state = stateDestType
		m.choices = destTypes
		m.cursor = 0

	case stateDestType:
		m.cfg.Destination.Type = destTypes[m.cursor]
		m.cursor = 0
		m.state = m.nextDestState()
		if m.isInputState() {
			m.inputs[0].Focus()
		}

	case stateDestURL:
		m.cfg.Destination.URL = m.inputs[0].Value()
		m.resetInput()
		m.state = m.nextDestState()
		if m.isInputState() {
			m.inputs[0].Focus()
		}

	case stateDestRepo:
		m.cfg.Destination.Repo = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestToken
		m.inputs[0].Focus()

	case stateDestToken:
		m.cfg.Destination.Token = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestDone

	case stateDestProjectID:
		m.cfg.Destination.ProjectID = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestToken
		m.inputs[0].Focus()

	case stateDestFilePath:
		m.cfg.Destination.Path = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestFileFormat
		m.inputs[0].Focus()

	case stateDestFileFormat:
		m.cfg.Destination.Format = m.inputs[0].Value()
		if m.cfg.Destination.Format == "" {
			m.cfg.Destination.Format = "json"
		}
		m.resetInput()
		m.state = stateDestEncrypt
		m.inputs[0].Focus()

	case stateDestEncrypt:
		val := strings.ToLower(m.inputs[0].Value())
		m.cfg.Destination.Encrypt = val == "y" || val == "yes"
		m.resetInput()
		m.state = stateDestDone

	case stateDestVaultAddr:
		m.cfg.Destination.VaultAddress = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestVaultPath
		m.inputs[0].Focus()

	case stateDestVaultPath:
		m.cfg.Destination.VaultPath = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestDone

	case stateDestEtcdEndpoints:
		endpoints := strings.Split(m.inputs[0].Value(), ",")
		for i := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoints[i])
		}
		m.cfg.Destination.EtcdEndpoints = endpoints
		m.resetInput()
		m.state = stateDestEtcdPrefix
		m.inputs[0].Focus()

	case stateDestEtcdPrefix:
		m.cfg.Destination.EtcdPrefix = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestDone

	case stateDestKubeNamespace:
		m.cfg.Destination.KubeNamespace = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestKubeSecretName
		m.inputs[0].Focus()

	case stateDestKubeSecretName:
		m.cfg.Destination.KubeSecretName = m.inputs[0].Value()
		m.resetInput()
		m.state = stateDestDone

	case stateDestDone:
		m.state = stateConfirm
		m.choices = []string{"Run sync", "Cancel"}
		m.cursor = 0

	case stateConfirm:
		if m.cursor == 0 {
			m.result = "run"
		} else {
			m.result = "cancel"
		}
		return m, tea.Quit
	default:
	}

	return m, nil
}

// resetInput clears the shared input field and resets echo mode.
func (m *model) resetInput() {
	m.inputs[0].Blur()
	m.inputs[0].Reset()
	m.inputs[0].EchoMode = textinput.EchoNormal
	m.inputs[0].Placeholder = ""
}

func (m model) nextSourceState() state {
	switch m.cfg.Source.Type {
	case "github":
		return stateSourceRepo
	case "gitlab":
		return stateSourceProjectID
	case "vault":
		return stateSourceVaultAddr
	case "etcd":
		return stateSourceEtcdEndpoints
	case "kubernetes":
		return stateSourceKubeNamespace
	default:
		return stateSourceToken
	}
}

func (m model) nextDestState() state {
	switch m.cfg.Destination.Type {
	case "file":
		return stateDestFilePath
	case "github":
		return stateDestRepo
	case "gitlab":
		return stateDestProjectID
	case "vault":
		return stateDestVaultAddr
	case "etcd":
		return stateDestEtcdEndpoints
	case "kubernetes":
		return stateDestKubeNamespace
	default:
		return stateDestToken
	}
}

func (m model) View() tea.View { //nolint:gocyclo,cyclop,funlen
	var b strings.Builder

	b.WriteString("╔══════════════════════════════════════╗\n")
	b.WriteString("║       SecretShift — TUI Setup        ║\n")
	b.WriteString("╚══════════════════════════════════════╝\n\n")

	switch m.state {
	case stateSourceType:
		b.WriteString("Select source type:\n\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString("  " + cursor + " " + choice + "\n")
		}

	case stateSourceRepo:
		b.WriteString("Enter source repository (owner/repo):\n\n")
		m.inputs[0].Placeholder = "owner/repo"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceToken:
		b.WriteString("Enter source token:\n\n")
		m.inputs[0].EchoMode = textinput.EchoPassword
		m.inputs[0].Placeholder = "token"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceProjectID:
		b.WriteString("Enter source project ID:\n\n")
		m.inputs[0].Placeholder = "12345"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceVaultAddr:
		b.WriteString("Enter Vault address:\n\n")
		m.inputs[0].Placeholder = "https://vault.example.com"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceVaultPath:
		b.WriteString("Enter Vault secret path:\n\n")
		m.inputs[0].Placeholder = "myapp/config"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceEtcdEndpoints:
		b.WriteString("Enter etcd endpoints (comma-separated):\n\n")
		m.inputs[0].Placeholder = "http://localhost:2379"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceEtcdPrefix:
		b.WriteString("Enter etcd prefix:\n\n")
		m.inputs[0].Placeholder = "/secrets/"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceKubeNamespace:
		b.WriteString("Enter Kubernetes namespace:\n\n")
		m.inputs[0].Placeholder = "default"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceKubeSecretName:
		b.WriteString("Enter Kubernetes secret name (leave empty for all):\n\n")
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateSourceKubeLabel:
		b.WriteString("Enter Kubernetes label selector (leave empty for none):\n\n")
		m.inputs[0].Placeholder = "app=myapp"
		b.WriteString("  " + m.inputs[0].View() + "\n")

	case stateProcessPrefix:
		b.WriteString("Add prefix to secret names? (leave empty to skip):\n\n")
		m.inputs[0].Placeholder = "PROD_"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateProcessSuffix:
		b.WriteString("Add suffix to secret names? (leave empty to skip):\n\n")
		m.inputs[0].Placeholder = "_ENV"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateProcessIncludeRegex:
		b.WriteString("Include regex filter (leave empty to skip):\n\n")
		m.inputs[0].Placeholder = "^DB_"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateProcessExcludeRegex:
		b.WriteString("Exclude regex filter (leave empty to skip):\n\n")
		m.inputs[0].Placeholder = "^DEBUG_"
		b.WriteString("  " + m.inputs[0].View() + "\n")

	case stateDestType:
		b.WriteString("Select destination type:\n\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString("  " + cursor + " " + choice + "\n")
		}

	case stateDestRepo:
		b.WriteString("Enter destination repository (owner/repo):\n\n")
		m.inputs[0].Placeholder = "owner/repo"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestToken:
		b.WriteString("Enter destination token:\n\n")
		m.inputs[0].EchoMode = textinput.EchoPassword
		m.inputs[0].Placeholder = "token"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestProjectID:
		b.WriteString("Enter destination project ID:\n\n")
		m.inputs[0].Placeholder = "12345"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestFilePath:
		b.WriteString("Enter output file path:\n\n")
		m.inputs[0].Placeholder = "./output/secrets.json"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestFileFormat:
		b.WriteString("Enter file format (json/yaml):\n\n")
		m.inputs[0].Placeholder = "json"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestEncrypt:
		b.WriteString("Encrypt file at rest? (y/n):\n\n")
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestVaultAddr:
		b.WriteString("Enter Vault address:\n\n")
		m.inputs[0].Placeholder = "https://vault.example.com"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestVaultPath:
		b.WriteString("Enter Vault secret path:\n\n")
		m.inputs[0].Placeholder = "myapp/config"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestEtcdEndpoints:
		b.WriteString("Enter etcd endpoints (comma-separated):\n\n")
		m.inputs[0].Placeholder = "http://localhost:2379"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestEtcdPrefix:
		b.WriteString("Enter etcd prefix:\n\n")
		m.inputs[0].Placeholder = "/secrets/"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestKubeNamespace:
		b.WriteString("Enter Kubernetes namespace:\n\n")
		m.inputs[0].Placeholder = "default"
		b.WriteString("  " + m.inputs[0].View() + "\n")
	case stateDestKubeSecretName:
		b.WriteString("Enter Kubernetes secret name:\n\n")
		m.inputs[0].Placeholder = "imported-secrets"
		b.WriteString("  " + m.inputs[0].View() + "\n")

	case stateConfirm:
		b.WriteString("Configuration summary:\n\n")
		b.WriteString("  Source:      " + m.cfg.Source.Type + "\n")
		if m.cfg.Source.Repo != "" {
			b.WriteString("  Source repo: " + m.cfg.Source.Repo + "\n")
		}
		if m.cfg.Source.ProjectID != "" {
			b.WriteString("  Source proj: " + m.cfg.Source.ProjectID + "\n")
		}
		if m.cfg.Source.VaultAddress != "" {
			b.WriteString("  Source vault: " + m.cfg.Source.VaultAddress + "\n")
		}
		b.WriteString("  Destination: " + m.cfg.Destination.Type + "\n")
		if m.cfg.Destination.Path != "" {
			b.WriteString("  Dest path:   " + m.cfg.Destination.Path + "\n")
		}
		if m.cfg.Destination.Repo != "" {
			b.WriteString("  Dest repo:   " + m.cfg.Destination.Repo + "\n")
		}
		if m.cfg.Destination.VaultAddress != "" {
			b.WriteString("  Dest vault:  " + m.cfg.Destination.VaultAddress + "\n")
		}
		if m.cfg.Process.AddPrefix != "" {
			b.WriteString("  Prefix:      " + m.cfg.Process.AddPrefix + "\n")
		}
		if m.cfg.Process.AddSuffix != "" {
			b.WriteString("  Suffix:      " + m.cfg.Process.AddSuffix + "\n")
		}
		b.WriteString("\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString("  " + cursor + " " + choice + "\n")
		}
	default:
	}

	b.WriteString("\n  ↑/↓: select  Enter: confirm  Esc: cancel\n")
	return tea.NewView(b.String())
}
