package models

import (
	"context"
	"errors"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/kTowkA/GophKeeper/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// initCFM инициализация модели для создания папки (директории)
func (ctrl *Controller) initCFM() {
	inputs := make([]textinputWithReq, 2)

	t := textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 50
	t.Placeholder = "название папки *"
	t.Focus()
	t.PromptStyle = styleRed
	t.TextStyle = styleRed
	inputs[0] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	t = textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 150
	t.Placeholder = "описание"
	inputs[1] = textinputWithReq{
		Model: t,
		isReq: false,
	}

	mcf := &modelWithInputs{
		opt:           ctrl.opt,
		inputs:        inputs,
		button:        "-> Создать <-",
		focusButton:   styleBlue.Render("-> Создать <-"),
		header:        "Заполните обязательные поля",
		successHeader: "Вы успешно создали папку",
		footer:        "Для возврата в предыдущее меню нажмите esc. Для выхода нажмите ctrl+c",
	}

	mcf.execFunc = func() error {
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		token, ok := mcf.ctx.Value(ctxToken).(string)
		if !ok {
			return ErrTokenUndefined
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(tokenTitle, token))
		resp, err := mcf.opt.Service().CreateFolder(ctx, &grpc.CreateFolderRequest{
			Title:       inputs[0].Model.Value(),
			Description: inputs[1].Model.Value(),
		})
		if err != nil {
			return err
		}
		if !resp.CreateFolderStatus {
			return errors.New(resp.CreateFolderMessage)
		}
		mcf.status = true
		return nil
	}

	ctrl.mCreateFolder = mcf
}

// initLM инициализация модели для входа пользователя
func (ctrl *Controller) initLM() {
	inputs := make([]textinputWithReq, 2)

	t := textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 50
	t.Placeholder = "логин *"
	t.Focus()
	t.PromptStyle = styleRed
	t.TextStyle = styleRed
	inputs[0] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	t = textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 50
	t.Placeholder = "пароль *"
	t.EchoCharacter = '*'
	t.EchoMode = textinput.EchoPassword
	inputs[1] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	ml := &modelWithInputs{
		opt:           ctrl.opt,
		inputs:        inputs,
		button:        "-> Авторизация <-",
		focusButton:   styleBlue.Render("-> Авторизация <-"),
		header:        "Заполните обязательные поля",
		successHeader: "Вы успешно авторизировались",
		footer:        "Для возврата в предыдущее меню нажмите esc. Для выхода нажмите ctrl+c",
	}
	ml.execFunc = func() error {
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		resp, err := ml.opt.Service().Login(ctx, &grpc.LoginRequest{
			Login:    inputs[0].Model.Value(),
			Password: inputs[1].Model.Value(),
		})
		if err != nil {
			return err
		}
		if !resp.LoginStatus {
			return errors.New(resp.LoginMessage)
		}
		ml.ctx = context.WithValue(context.Background(), ctxToken, resp.Token)
		ml.ctx = context.WithValue(ml.ctx, ctxPassword, inputs[1].Model.Value())
		ml.ctx = context.WithValue(ml.ctx, ctxUsername, inputs[0].Model.Value())
		ml.status = true
		return nil
	}

	ctrl.mLogin = ml
}

// initPGM инициализация модели для генерации пароля
func (ctrl *Controller) initPGM() {
	t := textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 3
	t.Placeholder = "4-100"
	t.Focus()
	t.PromptStyle = styleRed
	t.TextStyle = styleRed
	inputs := make([]textinputWithReq, 1)
	inputs[0] = textinputWithReq{
		Model: t,
		isReq: true,
	}
	mpg := &modelWithInputs{
		opt:           ctrl.opt,
		inputs:        inputs,
		button:        "-> Сгенерировать пароль<-",
		focusButton:   styleBlue.Render("-> Сгенерировать пароль <-"),
		header:        "Введите желаемую длину пароля (пароль будет не менее 4 символов)",
		successHeader: "Ваш сгенерированный пароль",
		footer:        "Для возврата в предыдущее меню нажмите esc. Для выхода нажмите ctrl+c",
	}
	// mPasswordGenerate.returnModel = Init(mPasswordGenerate.opt)
	mpg.execFunc = func() error {
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		val, err := strconv.Atoi(mpg.inputs[0].Value())
		if err != nil {
			return err
		}
		resp, err := mpg.opt.Service().GeneratePassword(ctx, &grpc.GeneratePasswordRequest{
			Length: int64(val),
		})
		if err != nil {
			return err
		}
		mpg.result = resp.Password
		mpg.status = true
		return nil
	}
	ctrl.mPasswordGenerate = mpg
}

// initRM инициализация модели для регистрации пользователя
func (ctrl *Controller) initRM() {
	inputs := make([]textinputWithReq, 2)

	t := textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 50
	t.Placeholder = "логин *"
	t.Focus()
	t.PromptStyle = styleRed
	t.TextStyle = styleRed
	inputs[0] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	t = textinput.New()
	t.Cursor.Style = styleRed
	t.CharLimit = 50
	t.Placeholder = "пароль *"
	t.EchoCharacter = '*'
	t.EchoMode = textinput.EchoPassword
	inputs[1] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	mr := &modelWithInputs{
		opt:           ctrl.opt,
		inputs:        inputs,
		button:        "-> Регистрация <-",
		focusButton:   styleBlue.Render("-> Регистрация <-"),
		header:        "Заполните обязательные поля",
		successHeader: "Вы успешно зарегистрировались",
		footer:        "Для возврата в предыдущее меню нажмите esc. Для выхода нажмите ctrl+c",
	}

	mr.execFunc = func() error {
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		resp, err := mr.opt.Service().Register(ctx, &grpc.RegisterRequest{
			Login:    inputs[0].Model.Value(),
			Password: inputs[1].Model.Value(),
		})
		if err != nil {
			return err
		}
		if !resp.RegisterStatus {
			return errors.New(resp.RegisterMessage)
		}
		mr.status = true
		return nil
	}

	ctrl.mRegistration = mr
}

// initSDM итнициализация модели для сохранения данных
func (ctrl *Controller) initSDM() {
	inputs := make([]textinputWithReq, 3)

	t := textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 50
	t.Placeholder = "название данных *"
	t.Focus()
	t.PromptStyle = styleRed
	t.TextStyle = styleRed
	inputs[0] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	t = textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 150
	t.Placeholder = "описание"
	inputs[1] = textinputWithReq{
		Model: t,
		isReq: false,
	}

	t = textinput.New()
	t.Cursor.Style = cursorStyle
	t.CharLimit = 150
	t.Placeholder = "значение *"
	inputs[2] = textinputWithReq{
		Model: t,
		isReq: true,
	}

	msd := &modelWithInputs{
		opt:           ctrl.opt,
		inputs:        inputs,
		button:        "-> Сохранить <-",
		focusButton:   styleBlue.Render("-> Сохранить <-"),
		header:        "Заполните обязательные поля",
		successHeader: "Вы успешно сохранили данные",
		footer:        "Для возврата в предыдущее меню нажмите esc. Для выхода нажмите ctrl+c",
	}

	msd.execFunc = func() error {
		password, ok := msd.ctx.Value(ctxPassword).(string)
		if password == "" || !ok {
			err := errors.New("не установлен пароль для сохранения")
			msd.err = err
			return err
		}
		evalue, err := ctrl.opt.Crypter().Encrypt([]byte(inputs[2].Model.Value()), password)
		if err != nil {
			return err
		}

		folder, ok := msd.ctx.Value(ctxFolder).(string)
		if folder == "" || !ok {
			err := errors.New("не выбрана папка для сохранения")
			msd.err = err
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		token, ok := msd.ctx.Value(ctxToken).(string)
		if !ok {
			return ErrTokenUndefined
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(tokenTitle, token))
		resp, err := msd.opt.Service().Save(ctx, &grpc.SaveRequest{
			Value: &grpc.KeeperValue{
				Title:       inputs[0].Model.Value(),
				Description: inputs[1].Model.Value(),
				Value:       evalue,
			},
			Folder: folder,
		})
		if err != nil {
			return err
		}
		if !resp.SaveStatus {
			return errors.New(resp.SaveMessage)
		}
		msd.status = true
		return nil
	}

	ctrl.mSaveData = msd
}

// initLVM инициализация модели для отображения данных внутри папки
func (ctrl *Controller) initLVM() {

	items := []list.Item{}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Выберите данные для просмотра"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	mfl := &modelWithList{
		list: l,
		opt:  ctrl.opt,
	}
	mfl.getValues = func() ([]string, error) {
		if mfl.ctx == nil {
			return nil, errors.New("не установлен контекст")
		}
		token, ok := mfl.ctx.Value(ctxToken).(string)
		if !ok {
			return nil, ErrTokenUndefined
		}
		folder, ok := mfl.ctx.Value(ctxFolder).(string)
		if !ok || folder == "" {
			return nil, errors.New("не выбрана папка с данными")
		}
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(tokenTitle, token))
		resp, err := mfl.opt.Service().Values(ctx, &grpc.ValuesInFolderRequest{
			Folder: folder,
		})
		if err != nil {
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.NotFound {
				return nil, nil
			}
			return nil, err
		}
		return resp.Values, nil
	}
	ctrl.mValuesList = mfl
}

// initLFM инициализация модели для отображения папок
func (ctrl *Controller) initLFM() {

	items := []list.Item{}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Выберите папку данных"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	mfl := &modelWithList{
		list: l,
		opt:  ctrl.opt,
	}
	mfl.getValues = func() ([]string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		token, ok := mfl.ctx.Value(ctxToken).(string)
		if !ok {
			return nil, ErrTokenUndefined
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(tokenTitle, token))
		resp, err := mfl.opt.Service().Folders(ctx, &grpc.FoldersRequest{})
		if err != nil {
			return nil, err
		}
		return resp.Folders, nil
	}
	ctrl.mFolderList = mfl
}

// initVM инициализация модели отображения данных
func (ctrl *Controller) initVM() {

	mv := &viewMoidel{
		opt:    ctrl.opt,
		header: "Сохраненные значения",
		footer: "Для возврата в предыдущее меню нажмите esc. Для выхода нажмите ctrl+c",
	}

	mv.getValues = func() ([]viewV, error) {
		if mv.ctx == nil {
			return nil, errors.New("не задан контекст")
		}
		token, ok := mv.ctx.Value(ctxToken).(string)
		if !ok {
			return nil, ErrTokenUndefined
		}
		folder, ok := mv.ctx.Value(ctxFolder).(string)
		if folder == "" || !ok {
			return nil, errors.New("не задана директория")
		}
		valueTitle, ok := mv.ctx.Value(ctxValues).(string)
		if valueTitle == "" || !ok {
			return nil, errors.New("не задано значение")
		}
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(tokenTitle, token))
		resp, err := mv.opt.Service().Load(ctx, &grpc.LoadRequest{
			Folder: folder,
			Title:  valueTitle,
		})
		if err != nil {
			return nil, err
		}
		values := make([]viewV, 3)
		values[0].title = "НАЗВАНИЕ"
		values[0].value = resp.Value.Title
		values[1].title = "ОПИСАНИЕ"
		values[1].value = resp.Value.Description
		password, ok := mv.ctx.Value(ctxPassword).(string)
		if password == "" || !ok {
			return nil, errors.New("не задан пароль")
		}

		val, err := mv.opt.Crypter().Decrypt(resp.Value.Value, password)
		if err != nil {
			return nil, err
		}
		values[2].title = "ЗНАЧЕНИЕ"
		values[2].value = string(val)
		return values, nil
	}
	ctrl.mValue = mv
}
