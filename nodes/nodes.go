package nodes

import pb "jw-cache/cachepb"

type NodePicker interface { // 节点选择器接口
	PickNode(key string) (node NodeGetter, ok bool)
}

//type NodeGetter interface {
//	Get(group string, key string) ([]byte, error)
//}

type NodeGetter interface { // 从远程节点获取值
	Get(in *pb.Request, out *pb.Response) error
}
