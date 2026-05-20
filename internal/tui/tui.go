package tui

import (
	"fmt"
	"strings"

	"github.com/PapaDanielVi/secret-shift/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
var conflictStrategies = []string{"replace", "skip", "report"}

type model struct {
	state        state
	cfg          *config.Config
	inputs       []textinput.Model
	cursor       int
	choices      []string
	selected     int
	err          error
	result       string
	quitting     bool
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

	final := m.(model)
	if final.quitting && final.result == "" {
		return nil, fmt.Errorf("TUI cancelled")
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			return m.handleEnter()

		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}

		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case tea.KeyTab:
			if m.state >= stateSourceURL && m.state <= stateDestConflictStrategy {
				m.cursor = (m.cursor + 1) % 2
			}
		}
	}

	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateSourceType:
		m.cfg.Source.Type = sourceTypes[m.cursor]
		m.cursor = 0
		m.state = m.nextSourceState()

	case stateSourceURL:
		m.cfg.Source.URL = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = m.nextSourceState()

	case stateSourceRepo:
		m.cfg.Source.Repo = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceToken

	case stateSourceToken:
		m.cfg.Source.Token = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceDone

	case stateSourceProjectID:
		m.cfg.Source.ProjectID = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceToken

	case stateSourceVaultAddr:
		m.cfg.Source.VaultAddress = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceVaultPath

	case stateSourceVaultPath:
		m.cfg.Source.VaultPath = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceToken

	case stateSourceEtcdEndpoints:
		endpoints := strings.Split(m.inputs[0].Value(), ",")
		for i := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoints[i])
		}
		m.cfg.Source.EtcdEndpoints = endpoints
		m.inputs[0].Reset()
		m.state = stateSourceEtcdPrefix

	case stateSourceEtcdPrefix:
		m.cfg.Source.EtcdPrefix = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceDone

	case stateSourceKubeNamespace:
		m.cfg.Source.KubeNamespace = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceKubeSecretName

	case stateSourceKubeSecretName:
		m.cfg.Source.KubeSecretName = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceKubeLabel

	case stateSourceKubeLabel:
		m.cfg.Source.KubeLabel = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateSourceDone

	case stateSourceDone:
		m.state = stateProcessPrefix
		m.cursor = 0

	case stateProcessPrefix:
		m.cfg.Process.AddPrefix = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateProcessSuffix

	case stateProcessSuffix:
		m.cfg.Process.AddSuffix = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateProcessIncludeRegex

	case stateProcessIncludeRegex:
		m.cfg.Process.IncludeRegex = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateProcessExcludeRegex

	case stateProcessExcludeRegex:
		m.cfg.Process.ExcludeRegex = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateProcessDone

	case stateProcessDone:
		m.state = stateDestType
		m.choices = destTypes
		m.cursor = 0

	case stateDestType:
		m.cfg.Destination.Type = destTypes[m.cursor]
		m.cursor = 0
		m.state = m.nextDestState()

	case stateDestURL:
		m.cfg.Destination.URL = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = m.nextDestState()

	case stateDestRepo:
		m.cfg.Destination.Repo = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestToken

	case stateDestToken:
		m.cfg.Destination.Token = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestConflictStrategy

	case stateDestProjectID:
		m.cfg.Destination.ProjectID = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestToken

	case stateDestFilePath:
		m.cfg.Destination.Path = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestFileFormat

	case stateDestFileFormat:
		m.cfg.Destination.Format = m.inputs[0].Value()
		if m.cfg.Destination.Format == "" {
			m.cfg.Destination.Format = "json"
		}
		m.inputs[0].Reset()
		m.state = stateDestEncrypt

	case stateDestEncrypt:
		val := strings.ToLower(m.inputs[0].Value())
		m.cfg.Destination.Encrypt = val == "y" || val == "yes"
		m.inputs[0].Reset()
		m.state = stateDestDone

	case stateDestVaultAddr:
		m.cfg.Destination.VaultAddress = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestVaultPath

	case stateDestVaultPath:
		m.cfg.Destination.VaultPath = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestDone

	case stateDestEtcdEndpoints:
		endpoints := strings.Split(m.inputs[0].Value(), ",")
		for i := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoints[i])
		}
		m.cfg.Destination.EtcdEndpoints = endpoints
		m.inputs[0].Reset()
		m.state = stateDestEtcdPrefix

	case stateDestEtcdPrefix:
		m.cfg.Destination.EtcdPrefix = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestDone

	case stateDestKubeNamespace:
		m.cfg.Destination.KubeNamespace = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestKubeSecretName

	case stateDestKubeSecretName:
		m.cfg.Destination.KubeSecretName = m.inputs[0].Value()
		m.inputs[0].Reset()
		m.state = stateDestDone

	case stateDestConflictStrategy:
		m.cfg.Destination.ConflictStrategy = m.inputs[0].Value()
		if m.cfg.Destination.ConflictStrategy == "" {
			m.cfg.Destination.ConflictStrategy = "replace"
		}
		m.inputs[0].Reset()
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
	}

	return m, nil
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

func (m model) View() string {
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
			b.WriteString(fmt.Sprintf("  %s %s\n", cursor, choice))
		}

	case stateSourceRepo:
		b.WriteString("Enter source repository (owner/repo):\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceToken:
		b.WriteString("Enter source token:\n\n")
		m.inputs[0].EchoMode = textinput.EchoPassword
		m.inputs[0].Placeholder = "token"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceProjectID:
		b.WriteString("Enter source project ID:\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceVaultAddr:
		b.WriteString("Enter Vault address:\n\n")
		m.inputs[0].Placeholder = "https://vault.example.com"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceVaultPath:
		b.WriteString("Enter Vault secret path:\n\n")
		m.inputs[0].Placeholder = "myapp/config"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceEtcdEndpoints:
		b.WriteString("Enter etcd endpoints (comma-separated):\n\n")
		m.inputs[0].Placeholder = "http://localhost:2379"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceEtcdPrefix:
		b.WriteString("Enter etcd prefix:\n\n")
		m.inputs[0].Placeholder = "/secrets/"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceKubeNamespace:
		b.WriteString("Enter Kubernetes namespace:\n\n")
		m.inputs[0].Placeholder = "default"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceKubeSecretName:
		b.WriteString("Enter Kubernetes secret name (leave empty for all):\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateSourceKubeLabel:
		b.WriteString("Enter Kubernetes label selector (leave empty for none):\n\n")
		m.inputs[0].Placeholder = "app=myapp"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))

	case stateProcessPrefix:
		b.WriteString("Add prefix to secret names? (leave empty to skip):\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateProcessSuffix:
		b.WriteString("Add suffix to secret names? (leave empty to skip):\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateProcessIncludeRegex:
		b.WriteString("Include regex filter (leave empty to skip):\n\n")
		m.inputs[0].Placeholder = "^DB_"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateProcessExcludeRegex:
		b.WriteString("Exclude regex filter (leave empty to skip):\n\n")
		m.inputs[0].Placeholder = "^DEBUG_"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))

	case stateDestType:
		b.WriteString("Select destination type:\n\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", cursor, choice))
		}

	case stateDestRepo:
		b.WriteString("Enter destination repository (owner/repo):\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestToken:
		b.WriteString("Enter destination token:\n\n")
		m.inputs[0].EchoMode = textinput.EchoPassword
		m.inputs[0].Placeholder = "token"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestProjectID:
		b.WriteString("Enter destination project ID:\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestFilePath:
		b.WriteString("Enter output file path:\n\n")
		m.inputs[0].Placeholder = "./output/secrets.json"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestFileFormat:
		b.WriteString("Enter file format (json/yaml):\n\n")
		m.inputs[0].Placeholder = "json"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestEncrypt:
		b.WriteString("Encrypt file at rest? (y/n):\n\n")
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestVaultAddr:
		b.WriteString("Enter Vault address:\n\n")
		m.inputs[0].Placeholder = "https://vault.example.com"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestVaultPath:
		b.WriteString("Enter Vault secret path:\n\n")
		m.inputs[0].Placeholder = "myapp/config"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestEtcdEndpoints:
		b.WriteString("Enter etcd endpoints (comma-separated):\n\n")
		m.inputs[0].Placeholder = "http://localhost:2379"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestEtcdPrefix:
		b.WriteString("Enter etcd prefix:\n\n")
		m.inputs[0].Placeholder = "/secrets/"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestKubeNamespace:
		b.WriteString("Enter Kubernetes namespace:\n\n")
		m.inputs[0].Placeholder = "default"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestKubeSecretName:
		b.WriteString("Enter Kubernetes secret name:\n\n")
		m.inputs[0].Placeholder = "imported-secrets"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))
	case stateDestConflictStrategy:
		b.WriteString("Conflict strategy (replace/skip/report):\n\n")
		m.inputs[0].Placeholder = "replace"
		b.WriteString(fmt.Sprintf("  %s\n", m.inputs[0].View()))

	case stateConfirm:
		b.WriteString("Configuration summary:\n\n")
		b.WriteString(fmt.Sprintf("  Source:      %s\n", m.cfg.Source.Type))
		if m.cfg.Source.Repo != "" {
			b.WriteString(fmt.Sprintf("  Source repo: %s\n", m.cfg.Source.Repo))
		}
		if m.cfg.Source.ProjectID != "" {
			b.WriteString(fmt.Sprintf("  Source proj: %s\n", m.cfg.Source.ProjectID))
		}
		b.WriteString(fmt.Sprintf("  Destination: %s\n", m.cfg.Destination.Type))
		if m.cfg.Destination.Path != "" {
			b.WriteString(fmt.Sprintf("  Dest path:   %s\n", m.cfg.Destination.Path))
		}
		if m.cfg.Destination.Repo != "" {
			b.WriteString(fmt.Sprintf("  Dest repo:   %s\n", m.cfg.Destination.Repo))
		}
		if m.cfg.Process.AddPrefix != "" {
			b.WriteString(fmt.Sprintf("  Prefix:      %s\n", m.cfg.Process.AddPrefix))
		}
		if m.cfg.Process.AddSuffix != "" {
			b.WriteString(fmt.Sprintf("  Suffix:      %s\n", m.cfg.Process.AddSuffix))
		}
		b.WriteString("\n")
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", cursor, choice))
		}
	}

	b.WriteString("\n  ↑/↓: select  Enter: confirm  Esc: cancel\n")
	return b.String()
}
