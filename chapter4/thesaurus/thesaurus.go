package thesaurus

type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}
