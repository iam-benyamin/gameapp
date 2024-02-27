package protobufmapper

import (
	"gameapp/contract/goproto/presence"
	"gameapp/param"
)

func MapToGetPresenceResponseProtobuf(g param.GetPresenceResponse) *presence.GetPresenceResponse {
	r := &presence.GetPresenceResponse{}

	for _, item := range g.Items {
		r.Items = append(r.Items, &presence.GetPresenceItem{
			UserId:    uint64(item.UserID),
			Timestamp: item.Timestamp,
		})
	}

	return r
}

func MapGetPresenceResponseFromProtobuf(g *presence.GetPresenceResponse) param.GetPresenceResponse {
	r := param.GetPresenceResponse{}

	for _, item := range g.Items {
		r.Items = append(r.Items, param.GetPresenceItem{
			UserID:    uint(item.UserId),
			Timestamp: item.Timestamp,
		})
	}

	return r
}
