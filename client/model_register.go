package client

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type register struct {
	gophKeeperState *Gophkeeper

	textInput textinput.Model
	username  string
	password  string
}

func initialModelRegister(state *Gophkeeper) register {
	ti := textinput.New()
	ti.Placeholder = "Ваш логин"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return register{
		textInput:       ti,
		gophKeeperState: state,
	}
}

func (m register) Init() tea.Cmd {
	return textinput.Blink
}

func (m register) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
			} else {
				m.password = m.textInput.Value()
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m register) View() string {
	if m.username == "" {
		return fmt.Sprintf(
			"Введите желаемое имя пользователя \n\n%s\n\n%s",
			m.textInput.View(),
			"(Нажмите Esc для возврата или ctrl+c для выхода)",
		) + "\n"
	}
	if m.password == "" {
		return fmt.Sprintf(
			"Введите желаемый пароль (достаточной длины и с символами в верхнем и нижнем регистрах, цифрами и специальными символами) \n\n%s\n\n%s",
			m.textInput.View(),
			"(Нажмите Esc для возврата или ctrl+c для выхода)",
		) + "\n"
	}
	return fmt.Sprintf("Вы успешно создали аккаунт с именем \"%s\"\n\n(Нажмите Esc для возврата или ctrl+c для выхода)\n", m.username)
}
