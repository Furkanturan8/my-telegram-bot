package services

// PingService, uygulamanın sağlık kontrolü için basit bir yanıt döndürür.
type PingService struct{}

func NewPingService() *PingService {
	return &PingService{}
}

func (s *PingService) Ping() string {
	return "OK"
}
