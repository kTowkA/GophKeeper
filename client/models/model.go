package models

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kTowkA/GophKeeper/client/models/options"
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

type model interface {
	WithContext(ctx context.Context) model
	WithNext(next model) model
	WithPrev(prev model) model
	tea.Model
}

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
	opt               *options.Options
}

func NewController(ctx context.Context, opt *options.Options) *Controller {
	ctrl := &Controller{
		opt: opt,
		ctx: ctx,
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
