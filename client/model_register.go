package client

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/grpc"
)

type register struct {
	service *Gophkeeper

	textInput textinput.Model
	username  string
	password  string
	errorMsg  error
}

func initialModelRegister(service *Gophkeeper) register {
	ti := textinput.New()
	ti.Placeholder = "Ваш логин"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return register{
		textInput: ti,
		service:   service,
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
				return Initial(m.service), nil
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
				ctx, cancel := context.WithTimeout(context.Background(), waitTime)
				defer cancel()
				_, err := m.service.gClient.Register(ctx, &grpc.RegisterRequest{Login: m.username, Password: m.password})
				m.errorMsg = err
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m register) View() string {
	if m.errorMsg != nil {
		defer func() {
			m.errorMsg = nil
		}()
		return fmt.Sprintf(
			"Произошла ошибка \n\n%s\n\n%s",
			m.errorMsg,
			modelMessageEscOrQuit,
		) + "\n"
	}

	if m.username == "" {
		return fmt.Sprintf(
			"Введите желаемое имя пользователя \n\n%s\n\n%s",
			m.textInput.View(),
			modelMessageEscOrQuit,
		) + "\n"
	}
	if m.password == "" {
		return fmt.Sprintf(
			"Введите желаемый пароль (достаточной длины и с символами в верхнем и нижнем регистрах, цифрами и специальными символами) \n\n%s\n\n%s",
			m.textInput.View(),
			modelMessageEscOrQuit,
		) + "\n"
	}
	return fmt.Sprintf("Вы успешно создали аккаунт с именем \"%s\"\n\n%s\n", m.username, modelMessageEscOrQuit)
}
