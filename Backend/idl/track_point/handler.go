package main

import (
	track_point "commerce/idl/track_point/kitex_gen/track-point"
	"context"
)

// TrackPointServiceImpl implements the last service interface defined in the IDL.
type TrackPointServiceImpl struct{}

// SendEvent implements the TrackPointServiceImpl interface.
func (s *TrackPointServiceImpl) SendEvent(ctx context.Context, req *track_point.SendEventRequest) (resp *track_point.SendEventResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryEvent implements the TrackPointServiceImpl interface.
func (s *TrackPointServiceImpl) QueryEvent(ctx context.Context, req *track_point.QueryEventRequest) (resp *track_point.QueryEventResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteEvent implements the TrackPointServiceImpl interface.
func (s *TrackPointServiceImpl) DeleteEvent(ctx context.Context, req *track_point.DeleteEventRequest) (resp *track_point.DeleteEventResponse, err error) {
	// TODO: Your code here...
	return
}
