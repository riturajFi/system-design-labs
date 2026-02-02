package model

type Channel string

const (
	ChannelPushIOS Channel = "push_ios"
	ChannelPushAndroid Channel = "push_android"
	ChannelSMS Channel = "sms"
	ChannelEmail Channel = "email"
)