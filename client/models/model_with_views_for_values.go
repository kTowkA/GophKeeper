// здесь описание модели где только просмотр данных возможен
package models

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

type viewMoidel struct {
	next      model
	prev      model
	ctx       context.Context
	values    []viewV
	header    string
	footer    string
	err       error
	getValues func() ([]viewV, error)
}
type viewV struct {
	title string
	value string
}

func (m *viewMoidel) WithContext(ctx context.Context) model {
	m.ctx = ctx
	m.updateItems()
	return m
}
func (m *viewMoidel) WithNext(next model) model {
	m.next = next
	return m
}
func (m *viewMoidel) WithPrev(prev model) model {
	m.prev = prev
	return m
}
func (m *viewMoidel) updateItems() *viewMoidel {
	values, err := m.getValues()
	if err != nil {
		m.err = err
		return m
	}
	m.values = values
	return m
}
func (m *viewMoidel) Init() tea.Cmd {
	return nil
}

func (m *viewMoidel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.ctx = context.WithValue(m.ctx, ctxValues, nil)
			if m.prev != nil {
				return m.prev.WithContext(m.ctx), nil
			}
		}
	}
	return m, nil
}

func (m *viewMoidel) View() string {
	view := m.header + "\n"
	if m.err != nil {
		view += "ошибка:" + m.err.Error() + "\n"
	}

	for _, v := range m.values {
		view += "\n" + " " + v.title + ": " + v.value
	}
	view += "\n\n" + m.footer
	return view
}
