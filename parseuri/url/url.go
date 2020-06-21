package url

import (
	"net/url"
	"path"
)

// JoinUrl(https://ip:port/*.*,path1,path2,path3...) => https://ip:port/*.*/path1/path2/path3
func JoinUrl(baseUrl string, paths ...string) string {
	u, _ := url.Parse(baseUrl)
	reqPath := path.Join(u.Path, path.Join(paths...))
	reqUrl, _ := url.ParseRequestURI(reqPath)
	u = u.ResolveReference(reqUrl)
	return u.String()
}
