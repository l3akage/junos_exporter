package subscribers

type SubscribersRPC struct {
	SubscribersInformation struct {
		Subscribers []Subscriber `xml:"subscriber"`
	} `xml:"subscribers-information"`
}

type Subscriber struct {
	Interface string `xml:"interface"`
	Username  string `xml:"user-name"`
}
