package raccoon

func IderFromNodeKeyer[Edg Edge[N], N any](nodeKeyer NodeKeyer[N]) Ider[Edg]{
    return &edgeIder[Edg, N]{
        nodeKeyer: nodeKeyer,
        edgeKeyer: NewEdgeKeyer(nodeKeyer, Key{}, nil, nil),
    }
}


type edgeIder[Edg Edge[N], N any] struct {
    nodeKeyer NodeKeyer[N]
    edgeKeyer EdgeKeyer[N]
}

func (i *edgeIder[Edg, N]) Id(edg Edg) []byte {
    return i.edgeKeyer.OutgoingKey(edg.GetSource(), edg.GetDest()).ToBytes()
}
