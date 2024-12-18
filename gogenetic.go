package gogenetic

import (
	"errors"
	"math/rand"
)

type MutationFunction[Phenotype any] func(Phenotype) Phenotype
type CrossoverFunction[Phenotype any] func(Phenotype, Phenotype) Phenotype
type FitnessFunction[Phenotype any] func(Phenotype) int64
type DoesABeatBFunction[Phenotype any] func(pa Phenotype, pb Phenotype) bool

type CustomSelectionFunction[Phenotype any] func(*GoGenetic[Phenotype]) *[]Phenotype

type GoGenetic[Phenotype any] struct {
	// public
	GType        GeneticType
	Generation   int
	MutationRate float32
	// func
	Mutate              MutationFunction[Phenotype]
	Crossover           CrossoverFunction[Phenotype]
	Fitness             FitnessFunction[Phenotype]
	DoesABeatB          *DoesABeatBFunction[Phenotype]
	CustomSelectionFunc CustomSelectionFunction[Phenotype]
	// private
	popSize    int
	population []Phenotype
}

// return the population
func (g *GoGenetic[Phenotype]) Population() *[]Phenotype {
	return &g.population
}

// change the population size
func (g *GoGenetic[Phenotype]) ChangePopSize(newPopSize int) error {
	if newPopSize < 1 {
		return errors.New("invalid new pop size")
	}
	if newPopSize == g.popSize {
		return nil
	}

	newPopulation := make([]Phenotype, newPopSize)
	copy(newPopulation, g.population)
	g.population = newPopulation

	return nil
}

// randomizePopulation
func (g *GoGenetic[Phenotype]) randomizePopulation(times int) {
	if g.popSize == 1 {
		return
	}
	for i := 0; i < times; i++ {
		rpa := rand.Intn(g.popSize)
		rpb := rand.Intn(g.popSize)
		for rpa == rpb {
			rpb = rand.Intn(g.popSize)
		}

		aux := g.population[rpa]
		g.population[rpa] = g.population[rpb]
		g.population[rpb] = aux
	}
}

// methods
func (g *GoGenetic[Phenotype]) competition() *[]Phenotype {
	nextGeneration := make([]Phenotype, g.popSize)

	for i := 0; i < g.popSize-1; i += 2 {
		pa := g.population[i]
		pb := g.population[i+1]

		nextGeneration[i] = pa

		if g.DoesABeatB != nil {
			if (*g.DoesABeatB)(pa, pb) {
				pb = g.Crossover(pa, pb)
			}
		} else {
			if g.Fitness(pa) >= g.Fitness(pb) {
				pb = g.Crossover(pa, pb)
			}
		}

		if rand.Float32() < g.MutationRate {
			pb = g.Mutate(pb)
		}

		nextGeneration[i+1] = pb
	}

	return &nextGeneration
}

func (g *GoGenetic[Phenotype]) rank() *[]Phenotype {
	// TODO:
	return nil
}

func (g *GoGenetic[Phenotype]) roulette() *[]Phenotype {
	// TODO:
	return nil
}

// evolve function
func (g *GoGenetic[Phenotype]) Evolve() int {
	times := g.popSize
	g.randomizePopulation(times)

	var nextGeneration *[]Phenotype
	switch g.GType {
	case RANK:
		nextGeneration = g.rank()
	case ROULETTE:
		nextGeneration = g.roulette()
	case COMPETITION:
		nextGeneration = g.competition()
	case CUSTOM:
		fallthrough
	default:
		nextGeneration = g.CustomSelectionFunc(g)
	}

	if nextGeneration != nil {
		newPopulation := make([]Phenotype, g.popSize)
		copy(newPopulation, *nextGeneration)
		g.population = newPopulation
	}

	return g.Generation
}

// constructor 1
func Empty[Phenotype any]() *GoGenetic[Phenotype] {
	return &GoGenetic[Phenotype]{
		GType:      COMPETITION,
		population: make([]Phenotype, 1),
		Mutate:     nil,
		Crossover:  nil,
		Fitness:    nil,
		DoesABeatB: nil,
		popSize:    1,
	}
}

func New[Phenotype any](
	initialPopulation []Phenotype,
	mutationRate float32,
	gType GeneticType,
	mutate MutationFunction[Phenotype],
	crossover CrossoverFunction[Phenotype],
	fitness FitnessFunction[Phenotype],
	doesABeatB *DoesABeatBFunction[Phenotype]) (*GoGenetic[Phenotype], error) {

	popSize := len(initialPopulation)
	if popSize < 1 {
		return nil, errors.New("invalid initial population")
	}

	population := make([]Phenotype, popSize)
	copy(population, initialPopulation)
	return &GoGenetic[Phenotype]{
		Generation:   0,
		MutationRate: mutationRate,
		GType:        gType,
		Mutate:       mutate,
		Crossover:    crossover,
		Fitness:      fitness,
		DoesABeatB:   doesABeatB,
		population:   population,
		popSize:      popSize,
	}, nil
}
