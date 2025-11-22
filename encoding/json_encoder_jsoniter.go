// +build jsoniter

package encoding

import jsoniter "github.com/json-iterator/go"

// JSON uses json‑iterator for zero‑allocation marshaling.
var JSON = jsoniter.ConfigFastest.Marshal
