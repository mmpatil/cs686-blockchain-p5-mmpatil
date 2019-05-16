package p3

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	// Route{
	// 	"MPTReceive",
	// 	"POST",
	// 	"/mptReceive/{key}/{message}",
	// 	MPTReceive,
	// },
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	Route{
		"StartClient",
		"GET",
		"/startClient",
		StartClient,
	},
	Route{
		"StartAuthServer",
		"GET",
		"/startAuth",
		StartAuthServer,
	},
	Route{
		"SignUp",
		"GET",
		"/signup",
		SignUp,
	},
	Route{
		"SignIn",
		"GET",
		"/signin",
		SignIn,
	},
	Route{
		"RegisterClient",
		"POST",
		"/registerClient",
		RegisterClient,
	},
	Route{
		"UserRegister",
		"GET",
		"/register/{nationalId}",
		UserRegister,
	},
	Route{
		"StartRegistrationServer",
		"GET",
		"/startReg",
		StartRegistrationServer,
	},
	Route{
		"DisplayUsers",
		"GET",
		"/displayUsers",
		DisplayUsers,
	},
	Route{
		"CheckUser",
		"POST",
		"/checkUser",
		CheckUser,
	},
	Route{
		"Check",
		"POST",
		"/check",
		Check,
	},
	Route{
		"ClientVote",
		"POST",
		"/clientVote",
		ClientVote,
	},
	Route{
		"VoteDetails",
		"POST",
		"/voteDetails",
		VoteDetails,
	},
	Route{
		"Vote",
		"POST",
		"/vote",
		Vote,
	},
	Route{
		"GetPeerList",
		"GET",
		"/getPeerList",
		GetPeerList,
	},
	Route{
		"ShowMPT",
		"GET",
		"/showMPT",
		ShowMPT,
	},
}
