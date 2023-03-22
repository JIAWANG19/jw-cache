package nodes

type NodePicker interface {
	PickNode(key string) (node NodeGetter, ok bool)
}

type NodeGetter interface {
	Get(group string, key string) ([]byte, error)
}
