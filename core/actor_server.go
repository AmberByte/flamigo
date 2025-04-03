package flamigo

type serverActor struct {
	callingInterface string
}

func (s *serverActor) Type() string {
	return TypeActorServer
}

func (s *serverActor) Interface() string {
	return s.callingInterface
}

func NewServerActor(interfaceName string) Actor {
	return &serverActor{callingInterface: interfaceName}
}
