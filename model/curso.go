package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"uva-bot/utils"
)

var urlCursoTypenane = "UrlCurso"

type UrlCurso struct {
	Tipo  string `json:"fnid"`
	Valor string `json:"nome,omitempty"`
}

func (mc *UrlCurso) GetDomainTypename() *string { return &urlCursoTypenane }

type Curso struct {
	Fnid              string      `json:"fnid"`
	FnidOriginal      *string     `json:"fnidReal,omitempty"`
	Original          *Curso      `json:"-"`
	Nome              *string     `json:"nome"`
	FnidCoordenador   *string     `json:"coordenador"`
	Coordenador       *Contato    `json:"-"`
	FnidRepresentante *string     `json:"representante"`
	Representante     *Contato    `json:"-"`
	Urls              []*UrlCurso `json:"urls"`
	FnidTags          []string    `json:"tags"`
	Tags              []*Tag      `json:"-"`
}

type CursoRepository struct {
	Entities  []*Curso
	fnidIndex map[string]*Curso
}

type CursoPredicateFn func(c *Curso) (bool, error)

var cursoTypenane = "Curso"
var cursoRepository *CursoRepository = nil

func (c *Curso) GetDomainTypename() string {
	return cursoTypenane
}

func (c *Curso) GetFnid() string {
	return c.Fnid
}

func (c *Curso) GetOriginal() *Curso {
	if c.Original != nil {
		return c.Original
	}
	return c
}

func GetCursoRepository() *CursoRepository {
	if cursoRepository == nil {
		tagRepository := GetTagRepository()
		contatoRepository := GetContatoRepository()

		filename := "data/cursos.json"

		log.Println("Lendo arquivo do repositório de " + cursoTypenane + ": " + filename)
		data, fileErr := ioutil.ReadFile(filename)
		utils.CheckError(fileErr)

		repository := CursoRepository{
			Entities:  []*Curso{},
			fnidIndex: map[string]*Curso{},
		}

		log.Println("Carregando dados para repositório de " + cursoTypenane)
		jsonErr := json.Unmarshal(data, &repository.Entities)
		utils.CheckError(jsonErr)

		for entityIdx := 0; entityIdx < len(repository.Entities); entityIdx++ {
			entity := repository.Entities[entityIdx]

			normalizedFnid, normalizeErr := utils.NormalizeFnid(entity.Fnid)
			utils.CheckError(normalizeErr)

			if normalizedFnid != entity.Fnid {
				log.Panicf("ERRO: Em %s[%s]: FNID \"%s\" não está devidamente normalizado! (deveria ser: \"%s\")", cursoTypenane, entity.Fnid, entity.Fnid, normalizedFnid)
			}

			_, duplicated := repository.fnidIndex[entity.Fnid]

			if duplicated {
				log.Panicf("ERRO: Em %s: FNID \"%s\" está duplicado!", cursoTypenane, entity.Fnid)
			}

			repository.fnidIndex[entity.Fnid] = entity

			if entity.Nome == nil {
				entity.Nome = &entity.Fnid
			}

			entity.Tags = []*Tag{}
			for tagIdx, tagFnid := range entity.FnidTags {
				tag, tagErr := tagRepository.FindByFnid(tagFnid)
				if tagErr != nil {
					log.Printf("AVISO: Em %s[%s].tags[%d]: %s\n", cursoTypenane, entity.Fnid, tagIdx, tagErr.Error())
					tag = &Tag{
						Fnid: tagFnid,
						Nome: &tagFnid,
					}
				}
				entity.Tags = append(entity.Tags, tag)
			}

			if entity.FnidCoordenador != nil {
				coordenador, coordenadorErr := contatoRepository.FindByFnid(*entity.FnidCoordenador)
				if coordenadorErr != nil {
					log.Println("Erro ao buscar coordenador \"" + (*entity.FnidCoordenador) + "\" do curso \"" + entity.Fnid + "\"!")
					panic(coordenadorErr)
				}
				entity.Coordenador = coordenador
			}

			if entity.FnidRepresentante != nil {
				representante, representanteErr := contatoRepository.FindByFnid(*entity.FnidRepresentante)
				if representanteErr != nil {
					log.Println("Erro ao buscar representante \"" + (*entity.FnidRepresentante) + "\" do curso \"" + entity.Fnid + "\"!")
					panic(representanteErr)
				}
				entity.Representante = representante
			}
		}

		for entityIdx := 0; entityIdx < len(repository.Entities); entityIdx++ {
			aliasEntity := repository.Entities[entityIdx]
			if aliasEntity.FnidOriginal == nil {
				continue
			}
			realEntity, realEntityOk := repository.fnidIndex[*aliasEntity.FnidOriginal]
			if !realEntityOk {
				panic(cursoTypenane + ": Substituto \"" + (*aliasEntity.FnidOriginal) + "\" de entidade \"" + aliasEntity.Fnid + "\" não existe!")
			}
			realEntity.Original = realEntity
		}

		cursoRepository = &repository
	}

	return cursoRepository
}

func (r *CursoRepository) FindByFnid(fnid string) (*Curso, error) {
	entity, ok := r.fnidIndex[fnid]
	if !ok {
		return nil, newFnidNotFoundError(cursoTypenane, fnid)
	}
	return entity, nil
}

func (r *CursoRepository) FindByPredicateFn(predicate CursoPredicateFn) ([]*Curso, error) {
	result := make([]*Curso, 0)
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
