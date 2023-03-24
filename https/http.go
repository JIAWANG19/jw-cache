package https

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"jw-cache/cache"
	pb "jw-cache/cachepb"
	"jw-cache/hashes"
	"jw-cache/nodes"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_jw_cache/" // 表示默认的基础路径，即缓存池中缓存项的URL前缀，默认为"/_jw_cache/"
	defaultReplicas = 50            // 表示默认的虚拟节点数，即每个节点在哈希环上的虚拟节点数，默认为50
)

// ConnectHTTPPool HTTP连接池
type ConnectHTTPPool struct {
	self       string                 // self 表示该池的连接的URL地址，即当前节点的地址
	basePath   string                 // basePath 表示该池的连接的基础路径，即缓存池中缓存项的URL前缀
	mu         sync.Mutex             // mu 互斥锁，用于保护节点列表的并发访问
	nodes      *hashes.Map            // nodes 哈希表，用于记录哈希值与节点的对应关系
	httpGetter map[string]*httpGetter // httpGetter 在当前节点获取不到缓存时，调用回调函数中其他节点获取
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
	group := cache.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)

	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()}) // 将消息对象序列化成二进制数据
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	//w.Write(view.ByteSlice())
	w.Write(body)
}

// Set 设置节点(初始化传入节点)，建立节点与哈希值的映射关系
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

// PickNode 当在当前节点获取不到值时，选择一个最可能获取到值的节点
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

// Get 发送http请求去其他节点获取值
func (p *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	// /baseURL?group=group&key=key
	u := fmt.Sprintf("%v%v/%v",
		p.baseURL,
		url.QueryEscape(in.Group),
		url.QueryEscape(in.Key))
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err = proto.Unmarshal(bytes, out); err != nil { // proto.Unmarshal() 将二进制数据反序列化为消息对象
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

// Get 发送http请求去其他节点获取值
//func (p *httpGetter) Get(group string, key string) ([]byte, error) {
//	// /baseURL?group=group&key=key
//	u := fmt.Sprintf("%v%v/%v", p.baseURL, url.QueryEscape(group), url.QueryEscape(key))
//	res, err := http.Get(u)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//	if res.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("server returned: %v", res.Status)
//	}
//
//	bytes, err := io.ReadAll(res.Body)
//	if err != nil {
//		return nil, fmt.Errorf("reading response body: %v", err)
//	}
//	return bytes, nil
//}
