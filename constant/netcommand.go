package constant

var NetCommand = struct {
	State          string
	Link           string
	StartSession   string
	EndSession     string
	RestartService string
}{
	State:          "state",
	Link:           "link",
	StartSession:   "start_session",
	EndSession:     "end_session",
	RestartService: "restart_service",
}
