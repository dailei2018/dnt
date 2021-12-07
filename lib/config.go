package lib

type Config struct {
	Names     []string
	Dnt2excel map[string]interface{}
	Dnt2Dnt   map[string]interface{}
	Dnt2Pak   map[string]interface{}
	Excel2dnt map[string]interface{}
}
