# Group Cache

go语言实现的分布式缓存项目

分布式缓存是一种将缓存数据分布式存储在多个服务器节点上的技术，其目的是提高系统的性能、可扩展性和容错性。当客户端请求数据时，系统首先查询缓存，如果缓存中存在该数据，则直接返回给客户端，否则系统从后端存储系统中获取数据并将其写入缓存，以便下次查询时可以直接从缓存中获取。

以下是使用Go语言实现分布式缓存的优点：

1. 并发性能：Go语言天生具备高并发性能，能够轻松处理大量的并发请求。
2. 网络编程：Go语言的标准库中包含了完善的网络编程支持，可以轻松实现TCP、UDP、HTTP等协议的网络通信。
3. 轻量级协程：Go语言使用轻量级协程，可以轻松创建数千个协程，并发执行多个任务，提高系统的性能。
4. 内存管理：Go语言具有自动内存管理功能，能够有效地管理内存使用，避免内存泄漏等问题。
5. 丰富的第三方库：Go语言拥有丰富的第三方库，包括了诸如Redis、Memcached、etcd等常用的分布式缓存组件的客户端库，可以轻松与这些组件进行集成。

总的来说，使用Go语言实现分布式缓存具有高并发性、良好的网络编程支持、轻量级协程、自动内存管理等优点，可以有效提升系统的性能和可扩展性。

[TOC]



## 缓存淘汰

### 常见的缓存淘汰策略

常见的缓存淘汰策略有以下几种：

1. 先进先出 (FIFO)：缓存中最先进入的数据被最先淘汰。
2. 最近最少使用 (LRU)：缓存中最近最少使用的数据被最先淘汰。
3. 最少使用 (LFU)：缓存中最少使用的数据被最先淘汰。
4. 随机 (RAND)：随机选择一个缓存块进行淘汰。
5. 最近未使用 (NRU)：将缓存块标记为“已使用”或“未使用”，并在淘汰时优先淘汰未使用的缓存块。
6. 时间轮 (Clock)：类似于NRU策略，但是将缓存块放在一个时间轮上，并按照轮子旋转的顺序进行淘汰。
7. 热度 (Heat)：缓存中数据的热度被用于决定哪些缓存块被淘汰，热度高的数据被优先保留。

#### 先进先出 (FIFO)

优点：

- 实现简单，易于理解和部署。
- 对于缓存数据没有任何的优先级考虑，公平性较高。

缺点：

- 不考虑数据的访问频率，容易导致缓存中存储了很多很少使用的数据，而淘汰了一些常用的数据，缓存效率不高。
- FIFO无法应对数据访问模式变化的情况，不能适应高频数据和突发流量。

#### 最近最少使用 (LRU)

最近最少使用（LRU）算法是一种常见的缓存淘汰策略，它的基本思想是在缓存空间不足时，优先淘汰最近最少使用的缓存数据。

具体来说，LRU算法会在缓存中记录每个缓存数据最近一次被访问的时间戳。当缓存空间不足时，LRU算法会找到最近最少使用的数据，也就是最长时间没有被访问过的数据，并将其淘汰，以腾出空间存储新的数据。

LRU算法的实现可以采用链表和哈希表的结合来实现。具体地，可以使用双向链表记录缓存数据的访问顺序，每次访问缓存数据时，将其移动到链表头部，以表示该数据是最近访问过的。当缓存空间不足时，可以从链表尾部淘汰最近最少使用的数据。

LRU算法的优点包括：

1. 能够最大限度地利用缓存空间，减少缓存不命中率，提高缓存命中率。
2. 算法实现相对简单，容易理解和实现。
3. 在实际应用中广泛使用，许多缓存实现都支持LRU算法。

而LRU算法的缺点是：

1. 对于访问模式较为复杂的应用程序，LRU算法可能会出现“缓存失效预测不准”的情况，即某些数据虽然很少访问，但在未来可能会频繁访问，而被错误地淘汰。
2. 需要维护缓存数据的访问时间戳，会占用一定的内存空间和计算资源。

总的来说，LRU算法是一种较为简单和高效的缓存淘汰策略，在实际应用中具有广泛的应用和良好的效果。

#### 最少使用 (LFU)

优点：

- 比LRU更准确地考虑了缓存数据的访问频率，淘汰最少使用的数据，可以更好地利用缓存空间。
- 在访问频率较低的情况下可以起到一定的防止数据被淘汰的作用。

缺点：

- 实现较为复杂，需要维护每个数据的使用次数，对于缓存中的数据更新频率较高时，可能导致算法性能下降。
- 对于访问频率变化较快的数据，淘汰策略可能不太准确。

### 最近最少使用(LRU)的实现

group cache使用**最近最少使用(LRU)**作为缓存的淘汰策略

#### Cache

在 LRU 缓存算法中，Cache 是 LRU Cache 的基本数据结构，它用于存储和管理缓存中的数据。Cache 是一个有容量限制的缓存，缓存的数据以键值对的形式存储，可以快速地添加、查询、删除数据。同时，Cache 还需要支持淘汰算法，以保证缓存的容量不会超过规定的最大容量。

**基本数据结构**：

```go
type Cache struct {
    maxBytes  int64
    nowBytes  int64
    ll        *list.List
    cache     map[string]*list.Element
    OnEvicted func(key string, value Value)
}
```

其中，`Cache` 结构体包含以下成员：

- `maxBytes`：最大内存限制
- `nowBytes`：当前已使用的内存
- `ll`：双向链表，用于记录每个键值对的访问频率，频率越高的越靠近链表头部，越不容易被淘汰
- `cache`：哈希表，用于记录每个键对应的值在链表中的位置
- `OnEvicted`：记录被删除时的回调函数，可选参数

**方法**：

| 方法名 | 输入参数                                      | 输出参数                  | 功能描述                                                     |
| ------ | --------------------------------------------- | ------------------------- | ------------------------------------------------------------ |
| New    | maxBytes int64, onEvicted func(string, Value) | *Cache                    | 创建并初始化 Cache 对象                                      |
| Get    | key string                                    | value Value, success bool | 查找 Cache 中对应 key 的值，若查找成功，则将该节点移至队首   |
| Remove | 无                                            | 无                        | 删除 Cache 中最近最少使用的节点                              |
| Add    | key string, value Value                       | 无                        | 将键值对存入 Cache 中，并将被操作的节点移到队列的队首，若发现已使用内存超过了最大内存，则调用回收方法 |

#### Group

在分布式缓存中，Group是一个重要的概念，其作用主要是对一组缓存数据进行管理和封装，包括缓存命名空间的隔离、缓存数据的过期策略、缓存数据的加载、缓存数据的修改和删除等。通过对缓存数据进行分组管理，可以更方便地对缓存进行管理和维护。

**基本数据结构**：

```go
// Group lru算法中的分组，相当于命名空间
type Group struct {
	name      string              // 组名，用于区分不同的缓存组
	getter    Getter              // 回调函数，当在缓存中没有查询到数据时，去调用回调函数获取数据
	mainCache cache               // 缓存的具体实现，使用 lru.Cache 实现缓存淘汰策略
	nodes     nodes.NodePicker    // 节点选择器，用于选择要缓存到哪个节点，从哪个节点获取数据，实现分布式缓存
	loader    *singleflight.Group // 防止缓存击穿的实现，保证只有一个 goroutine 去加载缓存
}
```

**方法**：

| 方法                                           | 说明                                                         |
| ---------------------------------------------- | ------------------------------------------------------------ |
| load(key string)                               | 根据key加载缓存，会根据节点选择器选择节点，若选择到了节点，则会从该节点获取数据，否则会从回调函数中获取数据 |
| GetFromNode(node nodes.NodeGetter, key string) | 从指定节点中获取数据                                         |
| Get(key string)                                | 根据key获取组内的值，若没有获取到，会尝试从其他数据源获取    |
| getLocally(key string)                         | 调用Getter从其他数据源获取数据，若获取到数据，将该数据存入缓存中 |
| populateCache(key string, value ByteView)      | 将获取到的数据存入缓存中                                     |

## HTTP服务端

Go语言中的标准库中包含了一个HTTP包，也称为net/http包，提供了一个HTTP客户端和服务器的实现。这个包提供了一系列的函数和类型，可以用于创建HTTP服务器和客户端，并处理HTTP请求和响应。

### ConnectHTTPPool

**基本数据结构**：

```go
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
```

**方法**：

| 方法名         | 描述                                                     |
| -------------- | -------------------------------------------------------- |
| Log            | 方便打印日志                                             |
| ServeHTTP      | 处理HTTP请求，用于获取缓存值                             |
| Set            | 设置节点，并建立节点与哈希值的映射关系                   |
| PickNode       | 当当前节点获取不到缓存值时，选择一个最可能获取到值的节点 |
| httpGetter.Get | 发送HTTP请求去其他节点获取缓存值                         |

其中，最核心的方法就是`ServeHTTP`方法

```go
// ServerHTTP HTTP 请求解析
func (p *ConnectHTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   // 请求需要请求前缀
   if !strings.HasPrefix(r.URL.Path, p.basePath) {
      panic("ConnectHTTPPool serving unexpected path: " + r.URL.Path)
   }
   p.Log("%s %s", r.Method, r.URL.Path)
   // 请求格式应当为：/basePath/groupName/key
   parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
   if len(parts) != 2 {
      http.Error(w, "bad request", http.StatusBadRequest)
      return
   }

   groupName, key := parts[0], parts[1]
   // 获取组
   group := cache.GetGroup(groupName)
   if group == nil {
      http.Error(w, "no such group: "+groupName, http.StatusNotFound)
      return
   }
   // 从组中获取数据
   view, err := group.Get(key)
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
   }

   w.Header().Set("Content-Type", "application/octet-stream")
   w.Write(view.ByteSlice())
}
```

## 一致性哈希

### 一致性哈希算法

一致性哈希算法（Consistent Hashing）是一种用于解决分布式系统中缓存、负载均衡等问题的算法。它的基本思想是将数据和服务器都映射到同一个哈希环上，并保证数据尽可能均匀地分布在哈希环上，同时将服务器映射到哈希环上的位置也尽可能均匀分布。这样，在需要查找缓存数据时，只需根据数据的哈希值在哈希环上查找到对应的位置，然后按顺时针方向找到最近的服务器即可。

具体来说，一致性哈希算法的实现步骤如下：

1. 将所有的缓存服务器和数据都映射到一个哈希环上，可以使用哈希函数对缓存数据和服务器进行哈希，得到一个哈希值，并将其映射到哈希环上的一个位置。
2. 在哈希环上查找缓存数据对应的位置，并按照顺时针方向查找到最近的服务器位置。
3. 当添加或删除一个服务器时，只需重新计算该服务器在哈希环上的位置，并将其负责的缓存数据迁移到该服务器上。

一致性哈希算法的优点包括：

1. 在添加或删除缓存服务器时，只需要重新计算和迁移部分缓存数据，而不需要全局重新分配，因此具有较好的可扩展性和灵活性。
2. 能够很好地解决节点故障导致的缓存失效问题，因为只需将该节点负责的缓存数据迁移到其他节点上即可。
3. 由于哈希环的使用，能够保证缓存数据的分布相对均匀，缓存命中率较高。

一致性哈希算法的缺点是：

1. 可能存在哈希值分布不均匀的问题，导致某些节点负载过重或过轻。
2. 需要维护一份节点列表，并保证节点列表的同步性，增加了实现的复杂度。

总的来说，一致性哈希算法是一种常用的分布式缓存和负载均衡的解决方案，它具有良好的可扩展性、灵活性和可靠性，是一种比较成熟和有效的算法。

**group cache使用虚拟节点的方式来解决节点负载不均衡的问题**

### Map

**基本数据结构**：

```go
// HashFunc 哈希函数，用于计算的哈希值
type HashFunc func(data []byte) uint32

type Map struct {
   hash     HashFunc       // hash 用于计算哈希值的哈希函数
   replicas int            // replicas 一个真实节点对应虚拟节点的数量
   keys     []int          // keys 该变量是一个有序列表，包含所有的虚拟节点
   hashMap  map[int]string // hashMap 该变量是一个哈希表，存储虚拟节点的哈希值和对应的真实节点名称
}
```

**方法**：

| 方法名 | 功能                        | 参数                              | 返回值 |
| ------ | --------------------------- | --------------------------------- | ------ |
| New    | 创建一个Map对象             | replicas: int, hashFunc: HashFunc | *Map   |
| Add    | 添加一个或多个节点到Map对象 | keys ...string                    | void   |
| Get    | 根据key获取节点名称         | key: string                       | string |

## 防止缓存击穿

### 缓存击穿

缓存击穿是指当一个非常热门的数据缓存过期或者被清空时，同时有大量的请求涌入数据库或其他持久层，导致缓存失效，从而产生大量的请求直接落到了持久层，导致持久层负载过高，甚至宕机的情况。缓存击穿通常会导致系统的响应时间变慢，甚至崩溃。

缓存击穿一般是由以下原因导致的：

1. 缓存失效：当一个热门的数据缓存过期或者被清空时，同时有大量的请求涌入系统，导致缓存失效。
2. 热点数据：当某个数据被频繁访问时，它就成为了热点数据。当这个热点数据的缓存失效时，会导致大量的请求涌入系统，从而导致缓存击穿。
3. 数据库压力：当缓存失效后，大量请求直接访问数据库，导致数据库压力过大。

为了解决缓存击穿问题，通常采用以下方法：

1. 缓存预热：在系统启动时，提前将热点数据加载到缓存中，避免在缓存失效时，大量请求直接访问数据库。
2. 数据库查询结果为空时，也要将这个key对应的value设置到缓存中，这样下次相同的请求就能从缓存中获取到数据，避免访问数据库。
3. 采用分布式缓存：通过将缓存数据分布到多个节点上，可以降低单个节点的负载压力，避免出现缓存击穿的情况。
4. 采用布隆过滤器：通过布隆过滤器判断一个 key 是否在缓存中存在，如果不存在，直接返回，避免了大量请求直接访问数据库的情况。
5. 缓存数据设置不同的过期时间：通过给不同的 key 设置不同的过期时间，避免多个 key 同时过期导致缓存击穿的情况。

### 使用分布式节点防止缓存击穿

分布式节点防止缓存击穿的原理是在缓存集群的多个节点之间分配缓存数据，从而将缓存数据分散到多个节点上，同时在缓存失效时，避免大量的请求集中到一个节点上进行数据库查询操作。

具体来说，当一个缓存键失效时，多个并发请求会同时到达缓存集群中的不同节点，因为缓存数据被分配到了多个节点上。这时，每个节点会先检查自己的缓存中是否存在该键的缓存数据，如果存在，就返回缓存数据给请求方；如果不存在，该节点会向数据库发起查询请求，查询到数据后再更新自己的缓存，并返回结果给请求方。

同时，为了防止多个节点同时查询数据库，需要通过某种机制对请求进行调度和限制。一种常见的方式是使用分布式锁，对相同的缓存键进行加锁，避免多个节点同时访问数据库。另外一种方式是使用一致性哈希算法，将缓存键分配到不同的节点上，从而保证每个节点都有机会响应请求。

**定义接口**：

```go
type NodePicker interface { // 节点选择器接口
	PickNode(key string) (node NodeGetter, ok bool)
}

type NodeGetter interface { // 从远程节点获取值
	Get(in *pb.Request, out *pb.Response) error
}
```

**实现方法**：

```go
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
```

### 使用互斥锁防止缓存击穿

假设在一个分布式系统中，有一个名为 Tom 的缓存 key。现在有 N 个并发请求通过 HTTP 访问缓存系统，其中 8003 节点向 8001 节点发起了 N 次请求。如果对于数据库访问没有做任何限制，那么很可能会产生 N 次对数据库的访问，从而容易导致缓存击穿和穿透问题。即使已经针对数据库做了防护，但是由于 HTTP 请求是非常耗费资源的操作，因此在针对相同的 key 进行访问时，8003 节点向 8001 节点发起三次请求是没有必要的。因此，在分布式系统中，可以采用互斥锁的方式来保证只有一个请求去访问节点 。

```go
// Do 防止缓存击穿的实现，当相同的key并发的请求时，该方法可以保证fn函数只被调用一次
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
   g.mu.Lock()           // 先上锁
   if g.callMap == nil { // 延迟加载
      g.callMap = make(map[string]*call)
   }
   if c, ok := g.callMap[key]; ok { // 如果有相同的key正在请求，则等待
      g.mu.Unlock()       // 解锁
      c.wg.Wait()         // 等待key请求完成
      return c.val, c.err // 返回key请求的结果
   }
   aCall := new(call)
   aCall.wg.Add(1)        // 发起请求前加锁，使请求结束前的所有与该请求相同的key阻塞
   g.callMap[key] = aCall // 添加到 g.callMap
   g.mu.Unlock()

   aCall.val, aCall.err = fn() // 调用方法获取key的值
   aCall.wg.Done()             // 请求解锁

   g.mu.Lock()
   delete(g.callMap, key) // 更新 g.callMap
   g.mu.Unlock()

   return aCall.val, aCall.err
}
```

该方法使用了 sync 包中的锁和 WaitGroup 实现了对相同 key 的并发请求进行阻塞，确保 fn 函数只被调用一次，从而避免了缓存穿透的问题。主要流程如下：

1. 首先获取锁，然后判断 g.callMap 是否为空，如果为空，则新建一个 g.callMap。
2. 接着判断 key 是否在 g.callMap 中已经有值，如果有则等待。否则，新建一个 call 对象，并将其添加到 g.callMap 中。
3. 在发起请求前，需要在 aCall 对象上加锁，这样就能保证在请求结束之前，所有与该请求相同的 key 都将被阻塞。
4. 然后调用 fn 方法获取 key 的值，获取完成后解锁。
5. 最后，需要再次获取锁来更新 g.callMap，并返回结果。

## 使用 Protobuf 进行通信

Protocol Buffers（简称 Protobuf）是一种轻量级的数据交换格式，它是由 Google 设计的，能够用于跨语言和平台的数据交换。Protobuf 通过使用二进制编码和压缩技术，使得数据传输和存储更加高效。它的基本原理如下：

1. 定义 .proto 文件：使用 Protobuf 时需要先定义一个 .proto 文件，该文件使用 Protobuf 官方提供的语法定义消息类型、字段名、类型和顺序等信息。通过编写 .proto 文件来描述需要传输的数据结构。
2. 编译 .proto 文件：在编写好 .proto 文件后，需要使用 Protobuf 的编译器将其编译成目标语言的代码，例如 C++、Java、Go 等。编译后生成的代码中包含了定义的消息类型、字段、序列化和反序列化等操作。
3. 序列化：将结构化数据序列化为二进制格式。序列化时，Protobuf 将消息按照 .proto 文件中的定义进行编码，并压缩数据。由于使用了二进制编码和压缩技术，因此序列化后的数据通常比 JSON、XML 等其他格式的数据更小，更加高效。
4. 反序列化：将二进制数据反序列化为结构化数据。反序列化时，Protobuf 读取二进制数据并按照 .proto 文件中的定义进行解码，生成结构化的数据。反序列化后的数据可以直接用于程序中的数据操作。

Protobuf 的优点包括：高效、跨语言、可扩展、易于维护和版本控制等。由于采用二进制编码和压缩技术，因此 Protobuf 在数据传输和存储方面的性能表现更加出色。同时，由于使用了标准化的 .proto 文件来描述数据结构，因此 Protobuf 可以支持版本控制和升级。此外，Protobuf 还支持向后和向前兼容，即旧版本的消息可以被新版本的代码读取，而新版本的消息也可以被旧版本的代码读取。

### 使用方法

1. 定义 .proto 文件：首先需要定义一个 .proto 文件来描述数据的结构和格式。这个文件中定义了需要传输的消息的字段名称、类型、顺序等信息。这个文件需要使用 Protobuf 官方提供的语法进行编写。

```protobuf
syntax = "proto3";

package cachepb;

option go_package = "jw-cache/cachepb";

message Request {
  string group = 1;
  string key = 2;
}

message Response {
  bytes value = 1;
}

service GroupCache {
  rpc Get(Request) returns (Response);
}
```

2. 编译 .proto 文件：编写好 .proto 文件后，需要使用 Protobuf 的编译器将其编译成 Go 语言代码。

3. 在 Go 项目中使用编译后的代码：编译 .proto 文件后会生成一些 .pb.go 文件，这些文件包含了 Protobuf 定义的数据结构和操作方法。在 Go 项目中需要使用这些代码来序列化和反序列化数据，或者使用这些代码生成消息对象进行通信。

```go
// ServerHTTP HTTP 请求解析
func (p *ConnectHTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   // ...
   view, err := group.Get(key)

   body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})  // // 将消息对象序列化成二进制数据
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
   }

   w.Header().Set("Content-Type", "application/octet-stream")
   //w.Write(view.ByteSlice())
   w.Write(body)
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
	if err = proto.Unmarshal(bytes, out); err != nil {  // proto.Unmarshal() 将二进制数据反序列化为消息对象
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
```
