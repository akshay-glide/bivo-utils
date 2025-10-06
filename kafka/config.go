package kafka

type KafkaConfig struct {
	BootstrapServers string `json:"bootstrap_servers" validate:"required"`
	GroupID          string `json:"group_id,omitempty"`
	Topic            string `json:"topic,omitempty"`

	// Producer specific
	EnableIdempotence bool `json:"enable_idempotence"`

	// Common optional settings
	Acks           string               `json:"acks,omitempty"` // "all", "1", "0"
	AutoOffset     string               `json:"auto_offset"`    // "earliest", "latest"
	SecurityConfig *KafkaSecurityConfig // Optional
}

type KafkaSecurityConfig struct {
	EnableTLS     bool   `json:"enable_tls"`
	EnableSASL    bool   `json:"enable_sasl"`
	SaslMechanism string `json:"sasl_mechanism"` // "PLAIN", "SCRAM-SHA-256", etc.
	SaslUsername  string `json:"sasl_username"`
	SaslPassword  string `json:"sasl_password"`
}
