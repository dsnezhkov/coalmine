package modules

type Processor interface {
	DemineFile(string, bool) (bool, error)
	GetHoneys() map[string][]string
}
