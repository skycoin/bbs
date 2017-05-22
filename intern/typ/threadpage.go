package typ

import (
	"encoding/hex"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// ThreadPage represents a ThreadPage as stored in cxo.
type ThreadPage struct {
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Posts  skyobject.References `skyobject:"schema=Post"`
}

func (tp ThreadPage) SerializeToHex() string {
	return hex.EncodeToString(encoder.Serialize(tp))
}

func (tp *ThreadPage) Deserialize(data []byte) error {
	return encoder.DeserializeRaw(data, tp)
}
