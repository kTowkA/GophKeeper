// здесь описание модели где пролисходит отображение списка с данными
package models

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14
const defaultWidth = 20

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type modelWithList struct {
	list      list.Model
	items     []list.Item
	choice    string
	err       error
	title     string
	next      model
	prev      model
	ctx       context.Context
	getValues func() ([]string, error)
}

func (m *modelWithList) WithContext(ctx context.Context) model {
	m.ctx = ctx
	m = m.updateItems().updateItems()
	return m
}
func (m *modelWithList) WithNext(next model) model {
	m.next = next
	return m
}
func (m *modelWithList) WithPrev(prev model) model {
	m.prev = prev
	return m
}
func (m *modelWithList) updateItems() *modelWithList {
	values, err := m.getValues()
	if err != nil {
		m.err = err
		return m
	}
	m.items = make([]list.Item, len(values))
	for i := range values {
		m.items[i] = item(values[i])
	}
	m.list = list.New(m.items, itemDelegate{}, defaultWidth, listHeight)
	return m
}

func (m *modelWithList) Init() tea.Cmd {
	return nil
}

func (m *modelWithList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.ctx = context.WithValue(m.ctx, ctxFolder, nil)
			m.choice = ""
			m.list.ResetSelected()
			return m.prev.WithContext(m.ctx), nil
		case "enter":
			m.choice = ""
			defer m.list.ResetSelected()
			i, ok := m.list.SelectedItem().(item)
			if ok {
				// m.choice = string(i)
				ctx := m.ctx
				if val, ok := ctx.Value(ctxFolder).(string); !ok || val == "" {
					ctx = context.WithValue(ctx, ctxFolder, string(i))
				} else {
					ctx = context.WithValue(ctx, ctxValues, string(i))
				}
				return m.next.WithContext(ctx), nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *modelWithList) View() string {
	return "\n" + m.list.View()
}
