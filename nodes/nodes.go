package nodes

import pb "jw-cache/cachepb"

type NodePicker interface {
	PickNode(key string) (node NodeGetter, ok bool)
}

//type NodeGetter interface {
//	Get(group string, key string) ([]byte, error)
//}

type NodeGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
