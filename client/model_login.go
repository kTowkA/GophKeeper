package client

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type login struct {
	gophKeeperState *Gophkeeper
	textInput       textinput.Model
	username        string
	password        string
}

func initialModelLogin(state *Gophkeeper) login {
	ti := textinput.New()
	ti.Placeholder = "ваш логин"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return login{
		textInput: ti,
	}
}

func (m login) Init() tea.Cmd {
	return textinput.Blink
}

func (m login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.username == "" || (m.username != "" && m.password != "") {
				return Initial(m.gophKeeperState), nil
			}
			m.username = ""
			m.password = ""
			return m, nil
		case "enter":
			if m.username == "" {
				m.username = m.textInput.Value()
				m.textInput.SetValue("")
				m.textInput.Placeholder = "ваш пароль"
			} else if m.password == "" {
				m.password = m.textInput.Value()
			} else {

			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m login) View() string {
	if m.username == "" {
		return fmt.Sprintf(
			"Введите имя пользователя \n\n%s\n\n%s",
			m.textInput.View(),
			"(Нажмите Esc для возврата или ctrl+c для выхода)",
		) + "\n"
	}
	if m.password == "" {
		return fmt.Sprintf(
			"Введите пароль \n\n%s\n\n%s",
			m.textInput.View(),
			"(Нажмите Esc для возврата или ctrl+c для выхода)",
		) + "\n"
	}
	return fmt.Sprintf("Вы успешно создали аккаунт с именем \"%s\"\n\n(Нажмите Esc для возврата или ctrl+c для выхода)\n", m.username)
}
