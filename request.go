package api

type Request struct {
	kind string
	path RequestPath
	prms RequestParams
}

type RequestPath string

const (
	PathEventParticipants RequestPath = "event/participants"
)
