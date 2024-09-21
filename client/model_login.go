package client

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/grpc"
)

type login struct {
	service   *Gophkeeper
	textInput textinput.Model
	username  string
	password  string
	errorMsg  error
}

func initialModelLogin(service *Gophkeeper) login {
	ti := textinput.New()
	ti.Placeholder = "ваш логин"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 30

	return login{
		textInput: ti,
		service:   service,
	}
}

func (m login) Init() tea.Cmd {
	return textinput.Blink
}

func (m login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.textInput.Placeholder = "ваш пароль"
			} else {
				m.password = m.textInput.Value()
				ctx, cancel := context.WithTimeout(context.Background(), waitTime)
				defer cancel()
				resp, err := m.service.gClient.Login(ctx, &grpc.LoginRequest{Login: m.username, Password: m.password})
				m.errorMsg = err
				if err == nil {
					m.service.token = resp.Token
				}
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m login) View() string {
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
			"Введите имя пользователя \n\n%s\n\n%s",
			m.textInput.View(),
			modelMessageEscOrQuit,
		) + "\n"
	}
	if m.password == "" {
		return fmt.Sprintf(
			"Введите пароль \n\n%s\n\n%s",
			m.textInput.View(),
			modelMessageEscOrQuit,
		) + "\n"
	}
	return fmt.Sprintf("Вы успешно вошли под именем \"%s\"\n\n%s\n", m.username, modelMessageEscOrQuit)
}
