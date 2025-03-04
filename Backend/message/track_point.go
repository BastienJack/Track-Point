package message

import trackpoint "commerce/idl/track_point/kitex_gen/track-point"

type AddCommonParamRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddCommonParamResponse struct {
	Base
}

type SendEventRequest struct {
	Event string `json:"event"`
}

type SendEventResponse struct {
	Base
}

type QueryEventRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type QueryEventResponse struct {
	Base
	Events []*trackpoint.Event `json:"events"`
}

type DeleteEventRequest struct {
	EventId uint64 `json:"event_id"`
}

type DeleteEventResponse struct {
	Base
}
