package sysproxy

import (
	"C"
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/getlantern/byteexec"
)
import "github.com/guodf/canary/basic"

//go:embed binaries/windows/sysproxy.exe
var sysproxy []byte

func ensureElevatedOnDarwin(be *byteexec.Exec) (err error) {
	return nil
}

func detach(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}
}

const (
	INTERNET_PER_CONN_FLAGS                = 1  // 用于设置开关，PROXY_TYPE_DIRECT|PROXY_TYPE_PROXY|PROXY_TYPE_AUTO_PROXY_URL|PROXY_TYPE_AUTO_DETECT
	INTERNET_PER_CONN_PROXY_SERVER         = 2  // 用于设置手动代理服务地址
	INTERNET_PER_CONN_PROXY_BYPASS         = 3  // 用于设置手动代理跳过的本地代理地址
	INTERNET_PER_CONN_AUTOCONFIG_URL       = 4  // 用于设置脚本代理地址
	INTERNET_OPTION_REFRESH                = 37 // 用于刷新使之前设置的代理信息
	INTERNET_OPTION_PER_CONNECTION_OPTION  = 75 // 用于查询系统代理设置
	INTERNET_OPTION_PROXY_SETTINGS_CHANGED = 95 // 用于保存设置的代理信息

	PROXY_TYPE_DIRECT         = 0x00000001 // 禁用所有代理
	PROXY_TYPE_PROXY          = 0x00000002 // 启用手动代理
	PROXY_TYPE_AUTO_PROXY_URL = 0x00000004 // 启用脚本代理
	PROXY_TYPE_AUTO_DETECT    = 0x00000008 // 启用自动检测设置

)

var (
	winInet              = syscall.NewLazyDLL("wininet.dll")
	InternetSetOptionW   = winInet.NewProc("InternetSetOptionW")
	InternetOpenW        = winInet.NewProc("InternetOpenW")
	InternetQueryOptionW = winInet.NewProc("InternetQueryOptionW")
)

/*
	typedef struct {
		  DWORD dwOption;
		  union {
		    DWORD    dwValue;
		    LPSTR    pszValue;
		    FILETIME ftValue;
		  } Value;
	} INTERNET_PER_CONN_OPTIONA, *LPINTERNET_PER_CONN_OPTIONA;

	typedef struct _FILETIME {
		  DWORD dwLowDateTime;
		  DWORD dwHighDateTime;
		} FILETIME, *PFILETIME, *LPFILETIME;
*/
type INTERNET_PER_CONN_OPTION struct {
	dwOption uint32
	dwValue  uint64 // 注意 32位 和 64位 struct 和 union 内存对齐
}
type INTERNET_PER_CONN_OPTION_LIST struct {
	dwSize        uint32
	pszConnection *uint16
	dwOptionCount uint32
	dwOptionError uint32
	pOptions      uintptr
}

func EnableAutoProxy() error {
	options := [3]INTERNET_PER_CONN_OPTION{}
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	options[0].dwValue = PROXY_TYPE_AUTO_DETECT
	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = 1
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))

	return setProxy(list)
}

func GetProxyStatus() SysProxySetting {
	options := [4]INTERNET_PER_CONN_OPTION{
		{dwOption: INTERNET_PER_CONN_FLAGS},
		{dwOption: INTERNET_PER_CONN_PROXY_SERVER},
		{dwOption: INTERNET_PER_CONN_PROXY_BYPASS},
		{dwOption: INTERNET_PER_CONN_AUTOCONFIG_URL},
	}

	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.dwOptionCount = uint32(len(options))
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))
	var listSize uint32 = uint32(unsafe.Sizeof(list))
	sysProxy := SysProxySetting{}
	if _, _, err := InternetQueryOptionW.Call(0, INTERNET_OPTION_PER_CONNECTION_OPTION, uintptr(unsafe.Pointer(&list)), uintptr(unsafe.Pointer(&listSize))); err != nil {
		for _, v := range *(*[len(options)]INTERNET_PER_CONN_OPTION)(unsafe.Pointer(list.pOptions)) {
			if v.dwOption == INTERNET_PER_CONN_FLAGS {
				sysProxy.EnableAutoProxy = v.dwValue&PROXY_TYPE_AUTO_DETECT == PROXY_TYPE_AUTO_DETECT
				sysProxy.PacProxy.Enabled = v.dwValue&PROXY_TYPE_AUTO_PROXY_URL == PROXY_TYPE_AUTO_PROXY_URL
				sysProxy.LocalProxy.Enabled = v.dwValue&PROXY_TYPE_PROXY == PROXY_TYPE_PROXY
				continue
			}
			if v.dwOption == INTERNET_PER_CONN_AUTOCONFIG_URL {
				sysProxy.PacProxy.PacAddr = basic.C_Utf16ToUtf8(uintptr(v.dwValue))
				continue
			}
			if v.dwOption == INTERNET_PER_CONN_PROXY_BYPASS {
				sysProxy.LocalProxy.BYPASSAddr = basic.C_Utf16ToUtf8(uintptr(v.dwValue))
				continue
			}
			if v.dwOption == INTERNET_PER_CONN_PROXY_SERVER {
				sysProxy.LocalProxy.Server = basic.C_Utf16ToUtf8(uintptr(v.dwValue))
				continue
			}
			log.Println(v)
		}
	}
	return sysProxy
}

func disableProxy() error {
	options := [3]INTERNET_PER_CONN_OPTION{}
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	options[0].dwValue = PROXY_TYPE_DIRECT
	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = 1
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))

	return setProxy(list)
}

func setGlobalProxy(proxy string, bypass string) error {
	if len(bypass) == 0 {
		bypass = "<local>"
	}
	var count uint32 = 3
	options := [3]INTERNET_PER_CONN_OPTION{} //make([]INTERNET_PER_CONN_OPTION, count)
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	options[0].dwValue = PROXY_TYPE_PROXY
	options[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER
	options[1].dwValue = uint64(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(proxy))))
	options[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS
	options[2].dwValue = uint64(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(bypass))))
	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = count
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))
	return setProxy(list)
}

func setPacProxy(pac string) error {
	options := [3]INTERNET_PER_CONN_OPTION{}
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	options[0].dwValue = PROXY_TYPE_AUTO_PROXY_URL
	options[1].dwOption = INTERNET_PER_CONN_AUTOCONFIG_URL
	options[1].dwValue = uint64(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pac))))
	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = 2
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))

	return setProxy(list)
}

func setProxy(list INTERNET_PER_CONN_OPTION_LIST) error {
	_, _, err := InternetSetOptionW.Call(0, INTERNET_OPTION_PER_CONNECTION_OPTION, uintptr(unsafe.Pointer(&list)), uintptr(unsafe.Sizeof(list)))
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_PER_CONNECTION_OPTION Error: %s", err)
	}
	_, _, err = InternetSetOptionW.Call(0, INTERNET_OPTION_PROXY_SETTINGS_CHANGED, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_PROXY_SETTINGS_CHANGED Error: %s", err)
	}
	_, _, err = InternetSetOptionW.Call(0, INTERNET_OPTION_REFRESH, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_REFRESH Error: %s", err)
	}
	return nil
}
