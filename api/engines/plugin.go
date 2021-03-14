package engines

type PluginEngineRunner interface {
	Run(address string, specification *Specification) error
}
