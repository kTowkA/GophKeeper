package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type modelMain struct {
	gophKeeperState *Gophkeeper
	choices         []choice
	cursor          int
}

type choice struct {
	value    string
	all      bool
	logged   bool
	unLogged bool
	model    tea.Model
}

func Initial(state *Gophkeeper) modelMain {
	return modelMain{
		choices: []choice{
			{
				value:    "Авторизация",
				unLogged: true,
				model:    initialModelLogin(state),
			},
			{
				value:    "Регистрация",
				unLogged: true,
				model:    initialModelRegister(state),
			},
			{
				value: "Генерация одноразовых паролей",
				all:   true,
				model: initialModelGeneratePassword(state),
			},
			{
				value:  "Сохранить данные",
				logged: true,
			},
		},
		gophKeeperState: state,
	}
}
func (m modelMain) Init() tea.Cmd {
	return nil
}
func (m modelMain) View() string {
	s := "Выберите действие\n\n"

	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice.value)
	}

	s += "\nНажмите q или ctrl+c для выхода.\n"

	return s
}

func (m modelMain) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m.choices[m.cursor].model, nil
		}
	}
	return m, nil
}
