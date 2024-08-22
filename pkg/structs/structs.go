package structs

type Data struct {
	Message       string
	Channel       string
	Guild         string
	HasMessageRef bool
	MessageRef    string
	SessionId     string
	Seq           int
}

func DiscordData() Data {
	return Data{}
}
