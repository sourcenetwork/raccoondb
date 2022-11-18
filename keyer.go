package raccoon


type EdgeKeyer[Node any] struct {
    keyer NodeKeyer[Node]
    inPrefix []byte
    outPrefix []byte
    base Key
}

func NewEdgeKeyer[Node any](keyer NodeKeyer[Node], base Key, inPrefix []byte, outPrefix []byte) EdgeKeyer[Node] {
    return EdgeKeyer[Node] {
        keyer: keyer,
        base: base,
        inPrefix: inPrefix,
        outPrefix: outPrefix,
    }
}

// OutgoingKey represents a key from source to target
func (k *EdgeKeyer[Node]) OutgoingKey(src, dst Node) Key {
    srcKey := k.keyer.Key(src)
    dstKey := k.keyer.Key(dst)
    return k.base.Append(k.outPrefix, srcKey, dstKey)
}

// OutgoingKey represents a key from target to source
func (k *EdgeKeyer[Node]) IncomingKey(src, dst Node) Key { 
    srcKey := k.keyer.Key(src)
    dstKey := k.keyer.Key(dst)
    return k.base.Append(k.inPrefix, dstKey, srcKey)
}


// Build a pair of keys for iteration.
// IterKeys are used to return all edges starting from Node
func (k *EdgeKeyer[Node]) SucessorsIterKeys(src Node) (Key, Key) {
    nodeKey := k.keyer.Key(src)
    base := k.base.Append(k.outPrefix, nodeKey)

    minKey := k.keyer.MinKey()
    maxKey := k.keyer.MaxKey()
    return base.Append(minKey), base.Append(maxKey)
}

func (k *EdgeKeyer[Node]) AncestorsIterKey(src Node) (Key, Key) {
    nodeKey := k.keyer.Key(src)
    base := k.base.Append(k.inPrefix, nodeKey)

    minKey := k.keyer.MinKey()
    maxKey := k.keyer.MaxKey()
    return base.Append(minKey), base.Append(maxKey)
}

type SecondaryIdxKeyer[Edg Edge[N], N any] struct {
    keyer NodeKeyer[N]
    mapper Mapper[Edg]
    idxName string
    baseKey Key
}

func NewSecIdxKeyer[Edg Edge[N], N any](keyer NodeKeyer[N], mapper Mapper[Edg], idxName string, prefix []byte) SecondaryIdxKeyer[Edg, N] {
    baseKey := Key{}.Append(prefix, []byte(idxName))
    return SecondaryIdxKeyer[Edg, N] {
        keyer: keyer,
        mapper: mapper,
        idxName: idxName,
        baseKey: baseKey,
    }
}

// Return key for a record
// Format: /{prefix}/{idxName}/{val}/{recordKey}
func (k *SecondaryIdxKeyer[Edg, N]) Key(edg Edg) Key {
    idxVal := k.mapper(edg)
    key := k.baseKey.Append(idxVal)

    edgKey := k.edgeKey(edg)
    return key.Join(edgKey)
}

func (k *SecondaryIdxKeyer[Edg, N]) edgeKey(edg Edg) Key {
    src := edg.GetSource()
    tgt := edg.GetDest()

    srcKey := k.keyer.Key(src)
    tgtKey := k.keyer.Key(tgt)

    key := Key{}.Append(srcKey, tgtKey)
    return key
}


func (k *SecondaryIdxKeyer[Edg, N]) IterKeys(idx []byte) (Key, Key) {
    base := k.baseKey.Append(idx)

    min, max := k.edgeKeyDomain()
    min = base.Join(min)
    max = base.Join(max)

    return min, max
}

func (k *SecondaryIdxKeyer[Edg, N]) edgeKeyDomain() (Key, Key) {
    min := k.keyer.MinKey()
    max := k.keyer.MaxKey()

    minKey := Key{}.Append(min, min)
    maxKey := Key{}.Append(max, max)
    return minKey, maxKey
}
