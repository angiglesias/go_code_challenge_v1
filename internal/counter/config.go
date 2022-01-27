package counter

var (
	DefaultConfig = Config{Cors: false, Listen: "0.0.0.0:3000"}
)

// Config contains configuration parameters to instantiate counter service
type Config struct {
	Cors   bool   // Activate CORS features
	Listen string //TCP address where server will listen to incoming connections (format "IP:Port")
}
