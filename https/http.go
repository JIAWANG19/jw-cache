package https

import (
	"JWCache/dao"
	"JWCache/hashes"
	"JWCache/nodes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_jw_cache/"
	defaultReplicas = 50
)

// ConnectHTTPPool HTTP连接池
type ConnectHTTPPool struct {
	self       string                 // self 该池的连接的url地址
	basePath   string                 // basePath 基本url前缀
	mu         sync.Mutex             // mu 锁
	nodes      *hashes.Map            // nodes 节点的哈希表
	httpGetter map[string]*httpGetter // httpGetter todo
}

// NewHTTPPool 新建连接池
func NewHTTPPool(self string) *ConnectHTTPPool {
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

// Set 设置节点(初始化传入节点)，建立节点与key的映射关系
func (p *ConnectHTTPPool) Set(nodes ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nodes = hashes.New(defaultReplicas, nil)
	p.nodes.Add(nodes...)
	p.httpGetter = make(map[string]*httpGetter, len(nodes))
	for _, node := range nodes {
		p.httpGetter[node] = &httpGetter{baseURL: node + p.basePath}
	}
}

// PickNode 选择真实节点
func (p *ConnectHTTPPool) PickNode(key string) (nodes.NodeGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if node := p.nodes.Get(key); node != "" && node != p.self {
		p.Log("Pick Node %s", node)
		return p.httpGetter[node], true
	}
	return nil, false
}

// httpGetter 主要实现实现实际的发送请求到真实节点去获取值的操作
type httpGetter struct {
	baseURL string
}

// Get 发送http请求去获取值
func (p *httpGetter) Get(group string, key string) ([]byte, error) {
	// /baseURL?group=group&key=key
	u := fmt.Sprintf("%v%v/%v", p.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	return bytes, nil
}
