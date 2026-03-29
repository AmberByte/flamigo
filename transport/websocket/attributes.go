package websocket

import "github.com/amberbyte/flamigo/strategies"

const (
	AttrAckKey = "ws.ack"
)

type Metadata struct {
	AckKey string
}

func AttachMetadata(req *strategies.Request, metadata Metadata) {
	if metadata.AckKey != "" {
		req.SetAttribute(AttrAckKey, metadata.AckKey)
	}
}

func AckKey(ctx strategies.Context) (string, bool) {
	value, ok := ctx.Request().Attribute(AttrAckKey)
	if !ok {
		return "", false
	}
	ack, ok := value.(string)
	return ack, ok
}
