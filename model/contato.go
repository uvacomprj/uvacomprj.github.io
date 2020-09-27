package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"uva-bot/utils"
)

// MeioContato é um tipo agregado contendo apenas informações de contato
// como email, site, etc. Sempre fará parte de um objeto Contato.
type MeioContato struct {
	Tipo     string   `json:"tipo"`
	Valor    string   `json:"valor"`
	FnidTags []string `json:"tags"`
	Tags     []*Tag   `json:"-"`
}

var meioContatoTypenane = "MeioContato"

// GetDomainTypename obtém o nome do tipo de MeioContato.
// Usado na renderização de respostas do bot.
func (mc *MeioContato) GetDomainTypename() *string { return &meioContatoTypenane }

// Contato representa informações de contato (pessoal, professor, aluno, etc).
//
// Contatos podem ser renderizados de forma diferente a depender
// do seu tipo. Os seguintes tipos são reconhecidos:
//
// - `docente`: Um professor, coordenador ou integrante do corpo docente.
// - `discente`: Um aluno (ex. representante, integrante de equipe na UVA, etc)
//
type Contato struct {
	Fnid         string         `json:"fnid"`
	FnidOriginal *string        `json:"fnidReal,omitempty"`
	Original     *Contato       `json:"-"`
	Tipo         string         `json:"tipo"`
	Nome         *string        `json:"nome"`
	FnidTags     []string       `json:"tags"`
	Tags         []*Tag         `json:"-"`
	MeiosContato []*MeioContato `json:"meiosContato"`
}

type ContatoPredicateFn func(c *Contato) (bool, error)

var contatoTypenane = "Contato"
var contatoRepository *ContatoRepository = nil

func (c *Contato) GetDomainTypename() string {
	return contatoTypenane
}

func (c *Contato) GetFnid() string {
	return c.Fnid
}

func (c *Contato) GetOriginal() *Contato {
	if c.Original != nil {
		return c.Original
	}
	return c
}

func (c *Contato) IsAlias() bool {
	return c.Original != nil
}

type ContatoRepository struct {
	Entities  []*Contato
	fnidIndex map[string]*Contato
}

func GetContatoRepository() *ContatoRepository {
	if contatoRepository == nil {
		tagRepository := GetTagRepository()

		filename := "data/contatos.json"

		log.Println("Lendo arquivo do repositório de " + contatoTypenane + ": " + filename)
		data, fileErr := ioutil.ReadFile(filename)
		utils.CheckError(fileErr)

		repository := ContatoRepository{
			Entities:  []*Contato{},
			fnidIndex: map[string]*Contato{},
		}

		log.Println("Carregando dados para repositório de " + contatoTypenane)
		jsonErr := json.Unmarshal(data, &repository.Entities)
		utils.CheckError(jsonErr)

		for entityIdx := 0; entityIdx < len(repository.Entities); entityIdx++ {
			entity := repository.Entities[entityIdx]

			normalizedFnid, normalizeErr := utils.NormalizeFnid(entity.Fnid)
			utils.CheckError(normalizeErr)

			if normalizedFnid != entity.Fnid {
				log.Panicf("ERRO: Em %s[%s]: FNID \"%s\" não está devidamente normalizado! (deveria ser: \"%s\")", contatoTypenane, entity.Fnid, entity.Fnid, normalizedFnid)
			}

			_, duplicated := repository.fnidIndex[entity.Fnid]

			if duplicated {
				log.Panicf("ERRO: Em %s: FNID \"%s\" está duplicado!", contatoTypenane, entity.Fnid)
			}

			repository.fnidIndex[entity.Fnid] = entity

			if entity.Nome == nil {
				entity.Nome = &entity.Fnid
			}

			if len(entity.Tipo) == 0 {
				entity.Tipo = "pessoa"
			}

			entity.Tags = []*Tag{}
			for tagIdx, tagFnid := range entity.FnidTags {
				tag, tagErr := tagRepository.FindByFnid(tagFnid)
				if tagErr != nil {
					log.Printf("AVISO: Em %s[%s].tags[%d]: %s\n", contatoTypenane, entity.Fnid, tagIdx, tagErr.Error())
					tag = &Tag{
						Fnid: tagFnid,
						Nome: &tagFnid,
					}
				}
				entity.Tags = append(entity.Tags, tag)
			}

			for meioContatoIdx, meioContato := range entity.MeiosContato {
				meioContato.Tags = []*Tag{}

				for tagIdx, tagFnid := range meioContato.FnidTags {
					tag, tagErr := tagRepository.FindByFnid(tagFnid)
					if tagErr != nil {
						log.Printf("AVISO: Em %s[%s].meiosContato[%d].tags[%d]: %s\n", contatoTypenane, entity.Fnid, meioContatoIdx, tagIdx, tagErr.Error())
						tag = &Tag{
							Fnid: tagFnid,
							Nome: &tagFnid,
						}
					}
					entity.Tags = append(entity.Tags, tag)
				}
			}
		}

		for entityIdx := 0; entityIdx < len(repository.Entities); entityIdx++ {
			aliasEntity := repository.Entities[entityIdx]
			if aliasEntity.FnidOriginal == nil {
				continue
			}
			realEntity, realEntityOk := repository.fnidIndex[*aliasEntity.FnidOriginal]
			if !realEntityOk {
				panic(contatoTypenane + ": Substituto \"" + (*aliasEntity.FnidOriginal) + "\" de entidade \"" + aliasEntity.Fnid + "\" não existe!")
			}
			realEntity.Original = realEntity
		}

		contatoRepository = &repository
	}

	return contatoRepository
}

func (r *ContatoRepository) FindByFnid(fnid string) (*Contato, error) {
	entity, ok := r.fnidIndex[fnid]
	if !ok {
		return nil, newFnidNotFoundError(contatoTypenane, fnid)
	}
	return entity, nil
}

func (r *ContatoRepository) FindByPredicateFn(predicate ContatoPredicateFn) ([]*Contato, error) {
	result := make([]*Contato, 0)
	for i := 0; i < len(r.Entities); i++ {
		entity := r.Entities[i]
		valid, err := predicate(entity)
		if err != nil {
			return nil, err
		}
		if valid {
			result = append(result, entity)
		}
	}
	return result, nil
}
