package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"uva-bot/utils"
)

type Tema struct {
	Fnid      string   `json:"fnid"`
	Nome      *string  `json:"nome"`
	Pai       *Tema    `json:"-"`
	FnidPai   *string  `json:"pai"`
	Titulo    *string  `json:"titulo"`
	Descricao *string  `json:"descricao"`
	Url       *string  `json:"url"`
	FnidTags  []string `json:"tags"`
	Tags      []*Tag   `json:"-"`
}

type TemaRepository struct {
	Entities  []*Tema
	fnidIndex map[string]*Tema
}

type TemaPredicateFn func(c *Tema) (bool, error)

var temaTypenane = "Tema"
var temaRepository *TemaRepository = nil

func (c *Tema) GetDomainTypename() string {
	return temaTypenane
}

func (c *Tema) GetFnid() string {
	return c.Fnid
}

func GetTemaRepository() *TemaRepository {
	if temaRepository == nil {
		tagRepository := GetTagRepository()

		filename := "data/temas.json"

		log.Println("Lendo arquivo do repositório de " + temaTypenane + ": " + filename)
		data, fileErr := ioutil.ReadFile(filename)
		utils.CheckError(fileErr)

		repository := TemaRepository{
			Entities:  []*Tema{},
			fnidIndex: map[string]*Tema{},
		}

		log.Println("Carregando dados para repositório de " + temaTypenane)
		jsonErr := json.Unmarshal(data, &repository.Entities)
		utils.CheckError(jsonErr)

		for entityIdx := 0; entityIdx < len(repository.Entities); entityIdx++ {
			entity := repository.Entities[entityIdx]

			if entity.Nome == nil || len(*entity.Nome) == 0 {
				log.Panicf("ERRO: Em %s: Um dos temas não possui nome!", temaTypenane)
			}

			if entity.Titulo != nil {
				entity.Titulo = entity.Nome
			}

			if entity.FnidPai != nil {
				fnidPai, fnidPaiErr := utils.NormalizeFnid(*entity.FnidPai)
				utils.CheckError(fnidPaiErr)
				entity.FnidPai = &fnidPai
			}

			fnid := entity.Fnid

			if len(fnid) == 0 {
				fnid = (*entity.Nome)
				// if entity.FnidPai != nil {
				// 	fnid = (*entity.FnidPai) + "/" + fnid
				// }
			}

			normalizedFnid, normalizeErr := utils.NormalizeFnid(fnid)
			utils.CheckError(normalizeErr)
			entity.Fnid = normalizedFnid

			_, duplicated := repository.fnidIndex[entity.Fnid]

			if duplicated {
				log.Panicf("ERRO: Em %s: FNID \"%s\" está duplicado!", temaTypenane, entity.Fnid)
			}

			repository.fnidIndex[entity.Fnid] = entity

			entity.Tags = []*Tag{}
			for tagIdx, tagFnid := range entity.FnidTags {
				tag, tagErr := tagRepository.FindByFnid(tagFnid)
				if tagErr != nil {
					log.Printf("AVISO: Em %s[%s].tags[%d]: %s\n", temaTypenane, entity.Fnid, tagIdx, tagErr.Error())
					tag = &Tag{
						Fnid: tagFnid,
						Nome: &tagFnid,
					}
				}
				entity.Tags = append(entity.Tags, tag)
			}
		}

		for _, entity := range repository.Entities {
			if entity.FnidPai != nil {
				pai, paiOk := repository.fnidIndex[*entity.FnidPai]
				if !paiOk {
					log.Println("Tema-pai \"" + (*entity.FnidPai) + "\" (do tema \"" + entity.Fnid + "\") não existe!")
				}
				entity.Pai = pai
			}
		}

		temaRepository = &repository
	}

	return temaRepository
}

func (r *TemaRepository) FindAll() ([]*Tema, error) {
	return r.Entities, nil
}

func (r *TemaRepository) FindByFnid(fnid string) (*Tema, error) {
	fnid, err := utils.NormalizeFnid(fnid)
	if err != nil {
		return nil, err
	}
	entity, ok := r.fnidIndex[fnid]
	if !ok {
		return nil, newFnidNotFoundError(temaTypenane, fnid)
	}
	return entity, nil
}

func (r *TemaRepository) FindTemasRaiz() ([]*Tema, error) {
	return r.FindByPredicateFn(func(tema *Tema) (bool, error) {
		return tema.Pai == nil, nil
	})
}

func (r *TemaRepository) FindTemasFilhos(temaPai *Tema) ([]*Tema, error) {
	return r.FindByPredicateFn(func(tema *Tema) (bool, error) {
		return tema.Pai != nil && tema.Pai.Fnid == temaPai.Fnid, nil
	})
}

func (r *TemaRepository) CountTemasFilhos(temaPai *Tema) (int64, error) {
	return r.CountByPredicateFn(func(tema *Tema) (bool, error) {
		return tema.Pai != nil && tema.Pai.Fnid == temaPai.Fnid, nil
	})
}

func (r *TemaRepository) FindByPredicateFn(predicate TemaPredicateFn) ([]*Tema, error) {
	result := make([]*Tema, 0)
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

func (r *TemaRepository) CountByPredicateFn(predicate TemaPredicateFn) (int64, error) {
	result := int64(0)
	for i := 0; i < len(r.Entities); i++ {
		entity := r.Entities[i]
		valid, err := predicate(entity)
		if err != nil {
			return 0, err
		}
		if valid {
			result = result + 1
		}
	}
	return result, nil
}
