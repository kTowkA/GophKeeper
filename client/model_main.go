package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type modelMain struct {
	service   *Gophkeeper
	choices   []choice
	cursor    int
	viewCount int
}

type choice struct {
	value            string
	viewOnlyLoggin   bool
	viewOnlyUnLoggin bool
	model            tea.Model
}

func Initial(service *Gophkeeper) modelMain {
	return modelMain{
		choices: []choice{
			{
				value:            "Авторизация",
				viewOnlyUnLoggin: true,
				model:            initialModelLogin(service),
			},
			{
				value:            "Регистрация",
				viewOnlyUnLoggin: true,
				model:            initialModelRegister(service),
			},
			{
				value: "Генерация одноразовых паролей",
				model: initialModelGeneratePassword(service),
			},
			{
				value:          "Сохранение данных",
				viewOnlyLoggin: true,
			},
			{
				value:          "Просмотр сохраненных данных",
				viewOnlyLoggin: true,
			},
		},
		service: service,
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
		if (choice.viewOnlyLoggin && m.service.IsLogged()) ||
			(choice.viewOnlyUnLoggin && !m.service.IsLogged()) ||
			(!choice.viewOnlyUnLoggin && !choice.viewOnlyLoggin) {
			s += fmt.Sprintf("%s %s\n", cursor, choice.value)
		}
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
