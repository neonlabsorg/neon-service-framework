package service

type PostServiceStartedEvent struct {
	serviceName string
}

func (e PostServiceStartedEvent) Name() string {
	return "post.service.created"
}

func (e PostServiceStartedEvent) IsAsynchronous() bool {
	return false
}

func (e PostServiceStartedEvent) ServiceName() string {
	return e.serviceName
}

type PostServiceOnlineEvent struct {
	serviceName string
}

func (e PostServiceOnlineEvent) Name() string {
	return "post.service.online"
}

func (e PostServiceOnlineEvent) IsAsynchronous() bool {
	return false
}

func (e PostServiceOnlineEvent) ServiceName() string {
	return e.serviceName
}
