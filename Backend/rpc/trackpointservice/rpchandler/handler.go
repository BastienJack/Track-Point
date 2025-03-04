package rpchandler

import (
	trackpoint "commerce/idl/track_point/kitex_gen/track-point"
	"commerce/idl/track_point/kitex_gen/track-point/trackpointservice"
	"commerce/pkg/etcd"
	"commerce/pkg/middleware"
	"commerce/pkg/viper"
	"commerce/pkg/zap"
	"commerce/storage/db"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"

	"context"
	"encoding/json"
	"fmt"
	"net"
)

var (
	eventConfig = viper.Init("event")

	eventServiceName = eventConfig.Viper.GetString("server.name")
	eventServiceAddr = fmt.Sprintf("%s:%d",
		eventConfig.Viper.GetString("server.host"),
		eventConfig.Viper.GetInt("server.port"))

	eventEtcdAddr = fmt.Sprintf("%s:%d",
		eventConfig.Viper.GetString("etcd.host"),
		eventConfig.Viper.GetInt("etcd.port"))

	serviceLogger = zap.InitLogger()
)

func StartTrackPointService() error {
	// trackpointservice registry
	r, err := etcd.NewEtcdRegistry([]string{eventEtcdAddr})
	if err != nil {
		serviceLogger.Fatalf("Server register failed: %v", err)
	}

	// trackpointservice resolver for trackpointservice discovery
	addr, err := net.ResolveTCPAddr("tcp", eventServiceAddr)
	if err != nil {
		serviceLogger.Fatalf("Resolver tcp addr failed: %v", err)
	}

	// initialize rpc server
	s := trackpointservice.NewServer(new(TrackPointServiceImpl),
		server.WithServiceAddr(addr),
		server.WithMiddleware(middleware.CommonMiddleware),
		server.WithMiddleware(middleware.ServerMiddleware),
		server.WithRegistry(r),
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),
		server.WithMuxTransport(),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: eventServiceName}),
	)

	// run rpc server
	if err := s.Run(); err != nil {
		return err
	}

	return nil
}

// TrackPointServiceImpl implements the last userservice interface defined in the IDL.
type TrackPointServiceImpl struct{}

// SendEvent implements the TrackPointServiceImpl interface.
func (s *TrackPointServiceImpl) SendEvent(ctx context.Context, req *trackpoint.SendEventRequest) (resp *trackpoint.SendEventResponse, err error) {
	logger := zap.InitLogger()

	// unmarshall json data
	var event *struct {
		EventName   string `json:"event_name"`
		EventParams string `json:"event_params"`
	}

	err = json.Unmarshal([]byte(req.JsonEventParams), &event)
	if err != nil {
		logger.Errorf("Error: unmarshall json data error. Info: error=%+v.", err)

		res := &trackpoint.SendEventResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	}

	// convert event to db format
	dbEvent := &db.Event{
		EventName:   event.EventName,
		EventParams: []byte(event.EventParams),
	}

	// add event to db
	err = db.AddEvent(ctx, dbEvent)
	if err != nil {
		logger.Errorf("Error: add event error. Info: error=%+v.", err)

		res := &trackpoint.SendEventResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	}

	// send event response
	res := &trackpoint.SendEventResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
	}

	return res, nil
}

// QueryEvent implements the TrackPointServiceImpl interface.
func (s *TrackPointServiceImpl) QueryEvent(ctx context.Context, req *trackpoint.QueryEventRequest) (resp *trackpoint.QueryEventResponse, err error) {
	logger := zap.InitLogger()

	// get request params
	var offset = int(req.Offset)
	var lmt = int(req.Limit)

	// exclude error params
	if offset < 0 || lmt < 0 {
		logger.Errorf("Error: invalid query params.")

		res := &trackpoint.QueryEventResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
			Events:     make([]*trackpoint.Event, 0),
		}

		return res, nil
	}

	// define event data structure
	var dbEvents []*db.Event
	var events []*trackpoint.Event

	/* query all events data */
	if lmt == 0 {
		dbEvents, err = db.GetAllEvents(ctx)

		// error response
		if err != nil {
			logger.Errorf("Error: query all events error. Info: error=%+v.", err)

			res := &trackpoint.QueryEventResponse{
				StatusCode: -1,
				StatusMsg:  "Error: server internal error.",
				Events:     make([]*trackpoint.Event, 0),
			}

			return res, nil
		}

		// success response
		for _, dbEvent := range dbEvents {
			events = append(events, &trackpoint.Event{
				EventId:     uint64(dbEvent.ID),
				EventName:   dbEvent.EventName,
				EventParams: string(dbEvent.EventParams),
			})
		}

		res := &trackpoint.QueryEventResponse{
			StatusCode: 0,
			StatusMsg:  "Success",
			Events:     events,
		}

		return res, nil
	}

	/* query page events data */
	dbEvents, err = db.GetPageEvents(ctx, offset, lmt)

	// error response
	if err != nil {
		logger.Errorf("Error: query all events error. Info: error=%+v.", err)

		res := &trackpoint.QueryEventResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
			Events:     make([]*trackpoint.Event, 0),
		}

		return res, nil
	}

	// success response
	for _, dbEvent := range dbEvents {
		events = append(events, &trackpoint.Event{
			EventId:     uint64(dbEvent.ID),
			EventName:   dbEvent.EventName,
			EventParams: string(dbEvent.EventParams),
		})
	}

	res := &trackpoint.QueryEventResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		Events:     events,
	}

	return res, nil
}

// DeleteEvent implements the TrackPointServiceImpl interface.
func (s *TrackPointServiceImpl) DeleteEvent(ctx context.Context, req *trackpoint.DeleteEventRequest) (resp *trackpoint.DeleteEventResponse, err error) {
	logger := zap.InitLogger()

	// get request params
	var eventId = uint(req.EventId)

	err = db.DeleteEventById(ctx, eventId)
	if err != nil {
		logger.Errorf("Error: delete events error. Info: error=%+v.", err)

		res := &trackpoint.DeleteEventResponse{
			StatusCode: -1,
			StatusMsg:  "Error: server internal error.",
		}

		return res, nil
	}

	res := &trackpoint.DeleteEventResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
	}

	return res, nil
}
