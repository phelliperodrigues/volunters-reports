package domain

// LocalidadeRepository define as operações de persistência para Localidade
type LocalidadeRepository interface {
	GetAll() (map[string]map[string]*Summary, error)
	Save(localidade *Localidade) error
}

// SetorRepository define as operações de persistência para Setor
type SetorRepository interface {
	GetAll() (map[string]*Setor, error)
	GetByLocalidade(localidade string) (*Setor, error)
}

// LivroRepository define as operações de persistência para Livros
type LivroRepository interface {
	GetByLocalidade(localidade string) (map[string]bool, error)
	GetAll() (map[string]map[string]bool, error)
}
