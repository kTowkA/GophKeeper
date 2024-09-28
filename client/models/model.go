package models

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// тип для значения контекста
type ctxTitle int

const (
	wait       = time.Second * 30
	tokenTitle = "token"
	ctxToken   = ctxTitle(iota)
	ctxUsername
	ctxPassword
	ctxFolder
	ctxValues
)

// model все наши кастомные модели должны реализовывать этот интерфейс для корректной передачи значений по стеку вызовов
type model interface {
	WithContext(ctx context.Context) model
	WithNext(next model) model
	WithPrev(prev model) model
	tea.Model
}

// Controller главная управляющая структура, которая содержит все используемые в настоящий момент модели
type Controller struct {
	ctx               context.Context
	main              *modelMain
	mCreateFolder     model
	mSaveData         model
	mLogin            model
	mRegistration     model
	mPasswordGenerate model
	mFolderList       model
	mValuesList       model
	mValue            model
	services          *Services
}

// NewController создание нашего контроллера и инициализация конкретных моделей
func NewController(ctx context.Context, services *Services) *Controller {
	ctrl := &Controller{
		services: services,
		ctx:      ctx,
	}
	ctrl.initMM()
	ctrl.initCFM()
	ctrl.initLM()
	ctrl.initPGM()
	ctrl.initRM()
	ctrl.initSDM()
	ctrl.initLFM()
	ctrl.initLVM()
	ctrl.initVM()
	return ctrl
}
