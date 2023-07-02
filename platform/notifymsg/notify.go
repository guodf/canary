package notifymsg

import "github.com/go-toast/toast"

func ShowMsg(appId, title, message string) {
	notify := toast.Notification{
		AppID:               appId,
		Title:               title,
		Message:             message,
		Icon:                "",
		ActivationType:      "",
		ActivationArguments: "",
		Actions:             []toast.Action{},
		Audio:               "",
		Loop:                false,
		Duration:            "",
	}
	notify.Push()
}
