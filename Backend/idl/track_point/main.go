package main

import (
	track_point "commerce/idl/track_point/kitex_gen/track-point/trackpointservice"
	"log"
)

func main() {
	svr := track_point.NewServer(new(TrackPointServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
