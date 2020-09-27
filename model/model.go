package model

import "errors"

// DomainType descreve um tipo de domínio.
//
// Tipos de domínio podem ser entidades ou objetos de valor.
//
// Entidades sempre terão um FNID (_Formatted Name ID_), que é usado não
// apenas para exibição (quando não houver informação mais adequada para
// nomear o objeto), mas também para identificá-lo unicamente.
// FNIDs são únicos apenas para o tipo ao qual eles se encontram inseridos.
//
// Objetos de valor não possuem FNIDs, e quase sempre se encontrarão
// agregados a algum tipo. Por conta desta agregação, o tempo de vida
// de um objeto de valor é o mesmo do seu objeto-pai.
type DomainType interface {

	// GetDomainTypename obtém o nome do tipo do objeto de domínio.
	// É utilizado primariamente para determinar como renderizar a
	// entidade nas respostas do bot.
	GetDomainTypename() string
}

// Entity descreve um tipo entidade no domínio.
type Entity interface {

	// GetFnid retorna o nome identificador único da entidade.
	GetFnid() string

	GetOriginal() *interface{}
}

func newFnidNotFoundError(entityTypeName string, fnid string) error {
	return errors.New(entityTypeName + " \"" + fnid + "\" não foi encontrado(a).")
}

func Init() {
	GetTagRepository()
	GetContatoRepository()
	GetCursoRepository()
}
