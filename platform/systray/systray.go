package systray

import (
	"runtime"

	"github.com/getlantern/systray"
)

type IMenu interface {
	GetMenuValue() *MenuValue // 绘制Menu
	Click()                   // Menu被点击
}

type DrawMenuEvent interface {
	DrawMenu() bool
}

type MenuValue struct {
	Title   string
	ToolTip string
	Icon    string
	OnClick func()
}

func (value *MenuValue) GetMenuValue() *MenuValue {
	return value
}

func (value *MenuValue) Click() {

}

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
	Title        string
	ToolTip      string
	Icon         string
	SelectIcon   string
	DisableIcon  string // 0 默认,1 选中 2 禁用
	Group        string // 分组，所有值相同的都在一个分组中
	SelectStatus int
	OnClick      func()
	Menus        []*Menu
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

func listenClick(menuItem *systray.MenuItem, menu *Menu) {
	if menu.OnClick != nil {
		<-menuItem.ClickedCh
		menu.OnClick()
		setIcon(menuItem, appMenu.GetIcon(menu.SelectIcon))
	}
}

func addMenu(menuItem *systray.MenuItem, menu *Menu) {
	subMenu := menuItem.AddSubMenuItem(menu.Title, menu.ToolTip)
	go listenClick(subMenu, menu)
	setIcon(subMenu, appMenu.GetIcon(menu.Icon))
	for _, item := range menu.Menus {
		addMenu(subMenu, item)
	}
}
func setIcon(menuItem *systray.MenuItem, icon []byte) {
	if len(icon) > 0 {
		menuItem.SetIcon(icon)
	}
}

func onReady() {
	systray.SetIcon(appMenu.GetIcon(appMenu.Icon))
	systray.SetTitle(appMenu.Title)
	systray.SetTooltip(appMenu.ToolTip)
	for _, item := range appMenu.Menus {
		menuItem := systray.AddMenuItem(item.Title, item.ToolTip)
		go listenClick(menuItem, item)
		setIcon(menuItem, appMenu.GetIcon(item.Icon))
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
