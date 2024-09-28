// здесь описание модели где возможен ввод данных
package models

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textinputWithReq struct {
	textinput.Model
	isReq bool
}
type modelWithInputs struct {
	focusIndex    int
	inputs        []textinputWithReq
	err           error
	status        bool
	button        string
	focusButton   string
	header        string
	successHeader string
	footer        string
	result        string
	execFunc      func() error
	prev          model
	next          model
	ctx           context.Context
}

func (m *modelWithInputs) WithContext(ctx context.Context) model {
	m.ctx = ctx
	return m
}
func (m *modelWithInputs) WithNext(next model) model {
	m.next = next
	return m
}
func (m *modelWithInputs) WithPrev(prev model) model {
	m.prev = prev
	return m
}
func (m *modelWithInputs) Init() tea.Cmd {
	return textinput.Blink
}
func (m *modelWithInputs) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.setInputsEmpty()
			m.setErrorEmpty()
			m.setStatus(false)
			return m.prev.WithContext(m.ctx), nil
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				if m.err != nil {
					m.err = nil
				}
				if m.reqIsEmpty() {
					break
				}
				m.err = m.execFunc()
				if m.err == nil && m.next != nil {
					return m.next.WithContext(m.ctx), nil
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styleRed
					m.inputs[i].TextStyle = styleRed
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}
func (m *modelWithInputs) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i].Model, cmds[i] = m.inputs[i].Model.Update(msg)
	}

	return tea.Batch(cmds...)
}
func (m *modelWithInputs) View() string {
	view := m.header + "\n\n"
	if m.status {
		view = m.successHeader
		if m.result != "" {
			view += "\n\n" + m.result
		}
		view += "\n\nДля отмены нажмите esc. Для выхода нажмите ctrl+c"
		return view
	}
	if m.err != nil {
		view += styleError.Render("ошибка: ", m.err.Error())
		view += "\n\n"
	}
	for _, v := range m.inputs {
		view += v.View() + "\n"
	}
	view += "\n"
	button := m.button
	if m.focusIndex == len(m.inputs) {
		button = m.focusButton
	}
	view += button

	view += "\n\n" + m.footer
	return view
}
func (m *modelWithInputs) setInputsEmpty() {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
}
func (m *modelWithInputs) setErrorEmpty() {
	m.err = nil
}
func (m *modelWithInputs) setStatus(st bool) {
	m.status = st
}
func (m *modelWithInputs) reqIsEmpty() bool {
	for _, v := range m.inputs {
		if v.Value() == "" && v.isReq {
			return true
		}
	}
	return false
}
