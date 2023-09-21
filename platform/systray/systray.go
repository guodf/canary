package systray

import (
	"runtime"

	"github.com/getlantern/systray"
)

type AppMenu struct {
	Title     string
	ToolTip   string
	Icon      string
	OnClick   func()
	OnExit    func()
	GetIcon   func(string) []byte
	Menus     []*Menu
	menuGroup map[string][]*Item
}

type Item struct {
	Menu     *Menu
	MenuItem *systray.MenuItem
}

type Menu struct {
	Id           string
	Title        string
	ToolTip      string
	Icon         string
	SelectIcon   string
	DefaultIcon  bool   // 使用默认图标
	Selected     bool   // 默认选中状态
	DisableIcon  string // 0 默认,1 选中 2 禁用
	Group        string // 分组，所有值相同的都在一个分组中
	SelectStatus int
	OnClick      func(*Menu)
	Menus        []*Menu
}

func (menu *Menu) Click(m *Menu) {
	menu.OnClick(menu)
}

var appMenu = &AppMenu{}
var isReady = false
var menusClick = map[string]*systray.MenuItem{}

// 触发某个菜单的点击事件
func ClickMenu(id string) {
	if click, ok := menusClick[id]; ok {
		click.ClickedCh <- struct{}{}
	}
}

func InitSystrayMenu(menu *AppMenu) {
	if isReady {
		return
	}
	appMenu = menu
	appMenu.menuGroup = make(map[string][]*Item)
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
		for _, item := range appMenu.menuGroup[menu.Group] {
			if item.Menu.Selected {
				item.Menu.Icon, item.Menu.SelectIcon = item.Menu.SelectIcon, item.Menu.Icon
				item.Menu.Selected = !item.Menu.Selected
			}
			setMenuItem(item.MenuItem, item.Menu)
		}
		menu.Icon, menu.SelectIcon = menu.SelectIcon, menu.Icon
		menu.Selected = !menu.Selected
		menu.OnClick(menu)
		setMenuItem(menuItem, menu)
	}
}

func setMenuItem(menuItem *systray.MenuItem, menu *Menu) {
	if len(menu.Group) > 0 {
		if _, ok := appMenu.menuGroup[menu.Group]; !ok {
			appMenu.menuGroup[menu.Group] = []*Item{}
		}
		has := false
		for _, item := range appMenu.menuGroup[menu.Group] {
			if item.Menu == menu {
				has = true
				break
			}
		}
		if !has {
			appMenu.menuGroup[menu.Group] = append(appMenu.menuGroup[menu.Group], &Item{
				Menu:     menu,
				MenuItem: menuItem,
			})
		}
	}
	menuItem.SetTitle(menu.Title)
	menuItem.SetTooltip(menu.ToolTip)
	if menu.DefaultIcon {
		if menu.Selected {
			menuItem.Check()
		} else {
			menuItem.Uncheck()
		}
	} else {
		menuItem.SetTemplateIcon(nil, appMenu.GetIcon(menu.Icon))
	}
	go listenClick(menuItem, menu)
}

func addMenuItem(menuItem *systray.MenuItem, menu *Menu) {
	if menuItem == nil {
		menuItem = systray.AddMenuItem(menu.Title, menu.ToolTip)
	} else {
		menuItem = menuItem.AddSubMenuItem(menu.Title, menu.ToolTip)
	}
	if len(menu.Id) > 0 {
		menusClick[menu.Id] = menuItem
	}
	setMenuItem(menuItem, menu)
	for _, item := range menu.Menus {
		addMenuItem(menuItem, item)
	}
}

func onReady() {
	systray.SetIcon(appMenu.GetIcon(appMenu.Icon))
	systray.SetTitle(appMenu.Title)
	systray.SetTooltip(appMenu.ToolTip)
	for _, item := range appMenu.Menus {
		addMenuItem(nil, item)
	}
}
func onExit() {
	if appMenu.OnExit != nil {
		appMenu.OnExit()
	}
}
