package https

import (
	"JWCache/dao"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_jw_cache/"

type ConnectHTTPPool struct {
	self     string
	basePath string
}

// NewPool 新建连接池
func NewPool(self string) *ConnectHTTPPool {
	return &ConnectHTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 打印日志
func (p *ConnectHTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServerHTTP HTTP 请求解析
func (p *ConnectHTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 请求需要请求前缀
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("ConnectHTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName, key := parts[0], parts[1]
	group := dao.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
