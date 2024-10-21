package datapipe

// Conf represents the configuration settings required to connect to a Datapipe service.
// It includes the following fields:
// - DatapipeURL: The URL of the Datapipe service.
// - DatapipeCertPath: The file path to the Datapipe client certificate.
// - DatapipeKeyPath: The file path to the Datapipe client key.
// - DatapipeCAPath: The file path to the Datapipe Certificate Authority (CA) certificate.
// - DatapipeTokenEndPoint: The endpoint for obtaining a token for the Datapipe service.
type Conf struct {
	DatapipeURL           string
	DatapipeCertPath      string
	DatapipeKeyPath       string
	DatapipeTokenEndPoint string
	LogLevel              string
}
