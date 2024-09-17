package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/grpc"
)

type generatePassword struct {
	gophKeeperState *Gophkeeper
	textInput       textinput.Model
	password        string
	err             string
}

func initialModelGeneratePassword(state *Gophkeeper) generatePassword {
	ti := textinput.New()
	ti.Placeholder = "10"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.Validate = func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}

	return generatePassword{
		textInput:       ti,
		gophKeeperState: state,
	}
}

func (m generatePassword) Init() tea.Cmd {
	return textinput.Blink
}

func (m generatePassword) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.password == "" {
				return Initial(m.gophKeeperState), nil
			}
			m.password = ""
			return m, nil
		case "enter":
			ctx, cancel := context.WithTimeout(context.Background(), waitTime)
			defer cancel()
			n, err := strconv.Atoi(m.textInput.Value())
			if err != nil {
				m.err = err.Error()
				return m, nil
			}
			resp, err := m.gophKeeperState.gClient.GeneratePassword(ctx, &grpc.GeneratePasswordRequest{
				Length: int32(n),
			})
			if err != nil {
				m.err = err.Error()
				return m, nil
			}
			m.password = resp.Password
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m generatePassword) View() string {
	defer func() {
		m.err = ""
	}()
	if m.err != "" {
		m.err = "\n\n" + m.err
	}
	if m.password == "" {
		return fmt.Sprintf(
			"Введите желаемую длину пароля (не менее 4 символов) %s\n\n%s\n\n%s",
			m.err,
			m.textInput.View(),
			"(Нажмите Esc для возврата или ctrl+c для выхода)",
		) + "\n"
	}
	return fmt.Sprintf("Ваш сгенерированный пароль %s\n\n%s\n\n(Нажмите Esc для возврата или ctrl+c для выхода)\n", m.err, m.password)
}
