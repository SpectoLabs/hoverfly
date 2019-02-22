package models

import (
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/SpectoLabs/hoverfly/core/state"
)

type Simulation struct {
	matchingPairs           []RequestMatcherResponsePair
	ResponseDelays          ResponseDelays
	ResponseDelaysLogNormal ResponseDelaysLogNormal
	RWMutex                 sync.RWMutex
}

func NewSimulation() *Simulation {

	return &Simulation{
		matchingPairs:           []RequestMatcherResponsePair{},
		ResponseDelays:          &ResponseDelayList{},
		ResponseDelaysLogNormal: &ResponseDelayLogNormalList{},
	}
}

func (this *Simulation) AddPair(pair *RequestMatcherResponsePair) {
	var duplicate bool
	this.RWMutex.Lock()
	for _, savedPair := range this.matchingPairs {
		duplicate = reflect.DeepEqual(pair.RequestMatcher, savedPair.RequestMatcher)
		if duplicate {
			break
		}
	}
	if !duplicate {
		this.matchingPairs = append(this.matchingPairs, *pair)
	}
	this.RWMutex.Unlock()
}

func (this *Simulation) AddPairInSequence(pair *RequestMatcherResponsePair, state *state.State) {
	var duplicate bool

	updates := map[int]RequestMatcherResponsePair{}

	var counter int
	sequenceKey := "sequence:0"

	pairNoState := pair.RequestMatcher
	pairNoState.RequiresState = nil

	this.RWMutex.Lock()
	for i, savedPair := range this.matchingPairs {

		savedPairNoState := savedPair.RequestMatcher
		savedPairNoState.RequiresState = nil

		duplicate = reflect.DeepEqual(pairNoState, savedPairNoState)
		if duplicate {
			counter = counter + 1

			if savedPair.RequestMatcher.RequiresState == nil {
				savedPair.RequestMatcher.RequiresState = map[string]string{}
			}

			if savedPair.Response.TransitionsState == nil {
				savedPair.Response.TransitionsState = map[string]string{}
			}

			if pair.RequestMatcher.RequiresState == nil {
				pair.RequestMatcher.RequiresState = map[string]string{}
			}
			sequenceKey = state.GetNewSequenceKey()
			for key := range savedPair.RequestMatcher.RequiresState {
				if strings.Contains(key, "sequence:") {
					sequenceKey = key
					break
				}
			}

			sequenceState := savedPair.RequestMatcher.RequiresState[sequenceKey]
			nextSequenceState := ""
			if sequenceState == "" {
				sequenceState = "1"
				nextSequenceState = "2"
				state.PatchState(map[string]string{sequenceKey: "1"})

			} else {
				currentSequenceState, _ := strconv.Atoi(sequenceState)
				nextSequenceState = strconv.Itoa(currentSequenceState + 1)
			}
			savedPair.RequestMatcher.RequiresState[sequenceKey] = sequenceState
			savedPair.Response.TransitionsState[sequenceKey] = nextSequenceState
			updates[i] = savedPair
		}
	}

	for i, updatedPair := range updates {
		this.matchingPairs[i] = updatedPair
	}

	if counter != 0 {
		pair.RequestMatcher.RequiresState[sequenceKey] = strconv.Itoa(counter + 1)
	}

	this.matchingPairs = append(this.matchingPairs, *pair)
	this.RWMutex.Unlock()
}

func (this *Simulation) GetMatchingPairs() []RequestMatcherResponsePair {
	this.RWMutex.RLock()
	pairs := this.matchingPairs
	this.RWMutex.RUnlock()
	return pairs
}

func (this *Simulation) DeleteMatchingPairs() {
	var pairs []RequestMatcherResponsePair
	this.RWMutex.Lock()
	this.matchingPairs = pairs
	this.RWMutex.Unlock()
}
