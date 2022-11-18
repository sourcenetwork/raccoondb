package raccoon

import (
	"testing"
        "github.com/stretchr/testify/assert"
        "google.golang.org/protobuf/reflect/protoreflect"
)

type testNode string

type testEdge struct {
    source testNode
    target testNode
}

func newTestEdg(source, target testNode) testEdge {
    return testEdge {
        source: source,
        target: target,
    }
}

func (edg *testEdge) GetSource() testNode {
    return edg.source
}

func (edg *testEdge) GetDest() testNode {
    return edg.target
}

func (edg *testEdge) ProtoReflect() protoreflect.Message{
    return nil
}

type testNodeKeyer struct { }
func (k *testNodeKeyer) Key(n testNode) []byte {
    return []byte(n)
}

func (k *testNodeKeyer) MaxKey() []byte {
    return []byte("zzz")
}

func (k *testNodeKeyer) MinKey() []byte {
    return []byte("aaa")
}


var _ NodeKeyer[testNode] = (*testNodeKeyer)(nil)
var _ Edge[testNode] = (*testEdge)(nil)

func TestEdgeKeyerOutgoingKey(t *testing.T) {
    inPrefix := []byte("in")
    outPrefix := []byte("out")
    keyer := NewEdgeKeyer[testNode](&testNodeKeyer{}, Key{}, inPrefix, outPrefix)

    key := keyer.OutgoingKey("abc", "def")
    
    assert.Equal(t, key.ToBytes(), []byte("/out/abc/def"))
}

func TestEdgeKeyerIncomingKey(t *testing.T) {
    inPrefix := []byte("in")
    outPrefix := []byte("out")
    keyer := NewEdgeKeyer[testNode](&testNodeKeyer{}, Key{}, inPrefix, outPrefix)

    key := keyer.IncomingKey("abc", "def")
    
    assert.Equal(t, key.ToBytes(), []byte("/in/def/abc"))
}

func TestEdgeKeyerSucessorIterKeys(t *testing.T) {
    inPrefix := []byte("in")
    outPrefix := []byte("out")
    keyer := NewEdgeKeyer[testNode](&testNodeKeyer{}, Key{}, inPrefix, outPrefix)

    min, max := keyer.SucessorsIterKeys("abc")
    
    assert.Equal(t, min.ToBytes(), []byte("/out/abc/aaa"))
    assert.Equal(t, max.ToBytes(), []byte("/out/abc/zzz"))
}

func TestEdgeKeyerAncestorIterKeys(t *testing.T) {
    inPrefix := []byte("in")
    outPrefix := []byte("out")
    keyer := NewEdgeKeyer[testNode](&testNodeKeyer{}, Key{}, inPrefix, outPrefix)

    min, max := keyer.AncestorsIterKey("abc")
    
    assert.Equal(t, min.ToBytes(), []byte("/in/abc/aaa"))
    assert.Equal(t, max.ToBytes(), []byte("/in/abc/zzz"))
}
