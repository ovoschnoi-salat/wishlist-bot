package main

import (
	"sync"
)

const (
	DefaultState      uint8 = iota
	NewWishState      uint8 = iota
	ReadUrlState      uint8 = iota
	ReadNewTitleState uint8 = iota
	NewFriendState    uint8 = iota
)

var DefaultUserState = UserState{
	ChosenWish: 0,
	State:      DefaultState,
}

func init() {
	userStates = &statesStorage{states: make(map[int64]UserState)}
}

var userStates StateStorage

type StateStorage interface {
	SetUserWholeState(id int64, state UserState) error
	SetUserState(id int64, state uint8) error
	SetUsersChosenWish(id int64, wish int64) error
	GetUserState(id int64) (UserState, error)
	DeleteUserState(id int64) error
}

type statesStorage struct {
	m      sync.Mutex
	states map[int64]UserState
}

func (s *statesStorage) SetUserWholeState(id int64, state UserState) error {
	s.m.Lock()
	defer s.m.Unlock()
	s.states[id] = state
	return nil
}

type UserState struct {
	ChosenWish int64
	State      uint8
}

func (s *statesStorage) SetUserState(id int64, newState uint8) error {
	s.m.Lock()
	defer s.m.Unlock()
	state := s.states[id]
	state.State = newState
	s.states[id] = state
	return nil
}

func (s *statesStorage) SetUsersChosenWish(id int64, newWish int64) error {
	s.m.Lock()
	defer s.m.Unlock()
	state := s.states[id]
	state.ChosenWish = newWish
	s.states[id] = state
	return nil
}

func (s *statesStorage) GetUserState(id int64) (UserState, error) {
	s.m.Lock()
	defer s.m.Unlock()
	state := s.states[id]
	//if !ok {
	//	return UserState{}, errors.New("no user found")
	//}
	return state, nil
}

func (s *statesStorage) DeleteUserState(id int64) error {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.states, id)
	return nil
}
