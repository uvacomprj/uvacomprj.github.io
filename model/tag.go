package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"uva-bot/utils"
)

type Tag struct {
	Fnid         string  `json:"fnid"`
	FnidOriginal *string `json:"fnidReal,omitempty"`
	Original     *Tag    `json:"-"`
	Nome         *string `json:"nome,omitempty"`
}

type TagPredicateFn func(c *Tag) (bool, error)

type TagRepository struct {
	Entities  []*Tag
	fnidIndex map[string]*Tag
}

var tagTypenane = "Tag"
var tagRepository *TagRepository = nil

func (c *Tag) GetDomainTypename() string {
	return contatoTypenane
}

func (c *Tag) GetFnid() string {
	return c.Fnid
}

func (c *Tag) GetOriginal() *Tag {
	if c.Original != nil {
		return c.Original
	}
	return c
}

func GetTagRepository() *TagRepository {
	if tagRepository == nil {
		filename := "data/tags.json"
		log.Println("Lendo arquivo do repositório de " + tagTypenane + ": " + filename)
		data, fileErr := ioutil.ReadFile(filename)
		utils.CheckError(fileErr)

		repository := TagRepository{
			Entities:  []*Tag{},
			fnidIndex: map[string]*Tag{},
		}

		log.Println("Carregando dados para repositório de " + contatoTypenane)
		jsonErr := json.Unmarshal(data, &repository.Entities)
		utils.CheckError(jsonErr)

		for i := 0; i < len(repository.Entities); i++ {
			entity := repository.Entities[i]

			normalizedFnid, normalizeErr := utils.NormalizeFnid(entity.Fnid)
			utils.CheckError(normalizeErr)

			if normalizedFnid != entity.Fnid {
				log.Panicf("ERRO: Em %s[%s]: FNID \"%s\" não está devidamente normalizado! (deveria ser: \"%s\")", tagTypenane, entity.Fnid, entity.Fnid, normalizedFnid)
			}

			_, duplicated := repository.fnidIndex[entity.Fnid]

			if duplicated {
				log.Panicf("ERRO: Em %s: FNID \"%s\" está duplicado!", tagTypenane, entity.Fnid)
			}

			repository.fnidIndex[entity.Fnid] = entity

			if entity.Nome == nil {
				entity.Nome = &entity.Fnid
			}
		}

		for entityIdx := 0; entityIdx < len(repository.Entities); entityIdx++ {
			aliasEntity := repository.Entities[entityIdx]
			if aliasEntity.FnidOriginal == nil {
				continue
			}
			realEntity, realEntityOk := repository.fnidIndex[*aliasEntity.FnidOriginal]
			if !realEntityOk {
				panic(tagTypenane + ": Substituto \"" + (*aliasEntity.FnidOriginal) + "\" de entidade \"" + aliasEntity.Fnid + "\" não existe!")
			}
			realEntity.Original = realEntity
		}

		tagRepository = &repository
	}

	return tagRepository
}

func (r *TagRepository) FindByFnid(fnid string) (*Tag, error) {
	tag, ok := r.fnidIndex[fnid]
	if !ok {
		return nil, newFnidNotFoundError(tagTypenane, fnid)
	}
	return tag, nil
}

func (r *TagRepository) FindByPredicateFn(predicate TagPredicateFn) ([]*Tag, error) {
	result := make([]*Tag, 0)
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
