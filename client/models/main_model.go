package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/client/models/options"
)

type pageType int

const (
	// основная страница
	pageMain = pageType(iota)
	// страница с выбором меню "создать новую папку и сохранить" или "выбрать уже существующую папку и сохранить"
	pageSaveAction
	// pageViewData
)

type choice struct {
	value string
	// отобразить только для авторизованных
	viewOnlyLoggin bool
	// отобразить только для неавторизованных
	viewOnlyUnLoggin bool
}
type page struct {
	choises []choice
	cursor  int
	header  string
	footer  string
}

type modelMain struct {
	opt         *options.Options
	pages       map[pageType]page
	currentPage pageType
}

func (ctrl *Controller) WithContext(ctx context.Context) model {
	ctrl.ctx = ctx
	return ctrl
}
func (ctrl *Controller) initMM() {
	mm := &modelMain{
		opt:         ctrl.opt,
		currentPage: pageMain,
		pages: map[pageType]page{
			pageMain: {
				cursor: 0,
				header: "Добро пожаловать в GophKeeper! Выберите доступное действие",
				footer: "Нажмите ctrl+c для выхода",
				choises: []choice{
					{
						value:            "РЕГИСТРАЦИЯ",
						viewOnlyUnLoggin: true,
					},
					{
						value:            "АВТОРИЗАЦИЯ",
						viewOnlyUnLoggin: true,
					},
					{
						value:          "СОХРАНИТЬ ДАННЫЕ",
						viewOnlyLoggin: true,
					},
					{
						value:          "ПРОСМОТРЕТЬ ДАННЫЕ",
						viewOnlyLoggin: true,
					},
					{
						value: "ГЕНЕРАЦИЯ ПАРОЛЕЙ",
					},
				},
			},
			pageSaveAction: {
				cursor: 0,
				header: "Выберите доступное действие",
				footer: "Нажмите esc для возврата или ctrl+c для выхода",
				choises: []choice{
					{
						value:          "ВЫБРАТЬ СУЩЕСТВУЮЩУЮ ПАПКУ",
						viewOnlyLoggin: true,
					},
					{
						value:          "СОЗДАТЬ НОВУЮ ПАПКУ",
						viewOnlyLoggin: true,
					},
				},
			},
		},
	}
	ctrl.main = mm
}

func (ctrl *Controller) Init() tea.Cmd {
	return nil
}
func (ctrl *Controller) View() string {
	view := ctrl.main.pages[ctrl.main.currentPage].header + "\n\n"

	if !ctrl.isShow(ctrl.main.pages[ctrl.main.currentPage].cursor) {
		for i := range ctrl.main.pages[ctrl.main.currentPage].choises {
			if ctrl.isShow(i) {
				cp := ctrl.main.pages[ctrl.main.currentPage]
				cp.cursor = i
				ctrl.main.pages[ctrl.main.currentPage] = cp
				break
			}
		}
	}

	for i, v := range ctrl.main.pages[ctrl.main.currentPage].choises {
		if !ctrl.isShow(i) {
			continue
		}
		if ctrl.main.pages[ctrl.main.currentPage].cursor == i {
			view += styleBlue.Render(fmt.Sprintf("%s %s", "->", v.value))
			view += "\n"
			continue
		}
		view += fmt.Sprintf("%s %s\n", "  ", v.value)
	}
	view += "\n" + ctrl.main.pages[ctrl.main.currentPage].footer + "\n"
	return view
}

func (ctrl *Controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentPage := ctrl.main.pages[ctrl.main.currentPage]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ctrl, tea.Quit
		case "up":
			if currentPage.cursor > 0 {
				for i := currentPage.cursor - 1; i >= 0; i-- {
					if ctrl.isShow(i) {
						currentPage.cursor = i
						ctrl.main.pages[ctrl.main.currentPage] = currentPage
						break
					}
				}
			}
		case "down":
			if currentPage.cursor != len(currentPage.choises)-1 {
				for i := currentPage.cursor + 1; i < len(currentPage.choises); i++ {
					if ctrl.isShow(i) {
						currentPage.cursor = i
						ctrl.main.pages[ctrl.main.currentPage] = currentPage
						break
					}
				}
			}
		case "esc":
			switch ctrl.main.currentPage {
			case pageSaveAction:
				ctrl.main.currentPage = pageMain
			}
		case "enter", " ":
			switch ctrl.main.currentPage {
			case pageMain:
				switch currentPage.cursor {
				case 0:
					return ctrl.mRegistration.WithPrev(ctrl).WithContext(ctrl.ctx), nil
				case 1:
					return ctrl.mLogin.WithPrev(ctrl).WithContext(ctrl.ctx), nil
				case 2:
					ctrl.main.currentPage = pageSaveAction
				case 3:
					view := ctrl.mValue
					values := ctrl.mValuesList
					folders := ctrl.mFolderList
					folders.WithPrev(ctrl).WithNext(values)
					values.WithPrev(folders).WithNext(view)
					view.WithPrev(values).WithNext(nil)
					return folders.WithContext(ctrl.ctx), nil
					// return ctrl.mFolderList.WithContext(ctrl.ctx).WithNext(ctrl.mValuesList.WithNext(ctrl.mValue.WithPrev(ctrl))), nil
				case 4:
					return ctrl.mPasswordGenerate.WithPrev(ctrl).WithContext(ctrl.ctx), nil
				}

			case pageSaveAction:
				switch currentPage.cursor {
				case 0:
					fmodels := ctrl.mFolderList
					sdmodel := ctrl.mSaveData.WithPrev(fmodels)
					fmodels.WithNext(sdmodel).WithPrev(ctrl)
					return fmodels.WithContext(ctrl.ctx), nil
				case 1:
					return ctrl.mCreateFolder.WithPrev(ctrl).WithContext(ctrl.ctx), nil
				}
			}
		}
	}
	return ctrl, nil
}

// isShow следует ли показывать пользователю пункт меню
func (ctrl *Controller) isShow(pos int) bool {
	token, _ := ctrl.ctx.Value(ctxToken).(string)

	ch := ctrl.main.pages[ctrl.main.currentPage].choises[pos]
	if (ch.viewOnlyLoggin && token != "") ||
		(ch.viewOnlyUnLoggin && token == "") ||
		(!ch.viewOnlyLoggin && !ch.viewOnlyUnLoggin) {
		return true
	}
	return false
}

func (ctrl *Controller) WithNext(next model) model {
	return ctrl
}
func (ctrl *Controller) WithPrev(prev model) model {
	return ctrl
}
