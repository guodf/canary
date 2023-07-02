package systray

import (
	"runtime"

	"github.com/getlantern/systray"
)

type AppMenu struct {
	Title   string
	ToolTip string
	Icon    string
	OnClick func()
	OnExit  func()
	GetIcon func(string) []byte
	Menus   []*Menu
}

type Menu struct {
	Title   string
	ToolTip string
	Icon    string
	OnClick func()
	Menus   []*Menu
}

var appMenu = AppMenu{
	Title:   "托盘",
	ToolTip: "托盘",
	Icon:    "",
	OnClick: func() {

	},
	Menus: []*Menu{
		{
			Title:   "托盘",
			ToolTip: "托盘",
			Icon:    "",
			OnClick: func() {

			},
			Menus: []*Menu{
				{
					Title:   "托盘",
					ToolTip: "托盘",
					Icon:    "",
					OnClick: func() {

					},
					Menus: []*Menu{
						{
							Title:   "托盘",
							ToolTip: "托盘",
							Icon:    "",
							OnClick: func() {

							},
						},
						{
							Title:   "托盘",
							ToolTip: "托盘",
							Icon:    "",
							OnClick: func() {

							},
						},
					},
				},
				{
					Title:   "托盘",
					ToolTip: "托盘",
					Icon:    "",
					OnClick: func() {

					},
				},
			},
		},
		{
			Title:   "退出",
			ToolTip: "退出",
			Icon:    "",
			OnClick: func() {

			},
		},
	},
}

func InitSystrayMenu(menu AppMenu) {
	appMenu = menu
	sysType := runtime.GOOS

	// if sysType == "linux" {
	// }

	if sysType == "windows" {
		systray.Run(onReady, onExit)
	}
}

func listenClick(menu *systray.MenuItem, click func()) {
	if click != nil {
		<-menu.ClickedCh
		click()
	}
}

func addMenu(menuItem *systray.MenuItem, menu *Menu) {
	subMenu := menuItem.AddSubMenuItem(menu.Title, menu.ToolTip)
	go listenClick(subMenu, menu.OnClick)
	menuItem.SetIcon(appMenu.GetIcon(menu.Icon))
	for _, item := range menu.Menus {
		addMenu(subMenu, item)
	}
}

func onReady() {
	systray.SetIcon(appMenu.GetIcon(appMenu.Icon))
	systray.SetTitle(appMenu.Title)
	systray.SetTooltip(appMenu.ToolTip)

	for _, item := range appMenu.Menus {
		menuItem := systray.AddMenuItem(item.Title, item.ToolTip)
		go listenClick(menuItem, item.OnClick)
		menuItem.SetIcon(appMenu.GetIcon(item.Icon))

		for _, menu := range item.Menus {
			addMenu(menuItem, menu)
		}
	}
}
func onExit() {
	if appMenu.OnExit != nil {
		appMenu.OnExit()
	}
}
