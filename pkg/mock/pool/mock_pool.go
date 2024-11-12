// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/pool/pool.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/pg-sharding/spqr/pkg/config"
	kr "github.com/pg-sharding/spqr/pkg/models/kr"
	pool "github.com/pg-sharding/spqr/pkg/pool"
	shard "github.com/pg-sharding/spqr/pkg/shard"
)

// MockConnectionKepper is a mock of ConnectionKepper interface.
type MockConnectionKepper struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionKepperMockRecorder
}

// MockConnectionKepperMockRecorder is the mock recorder for MockConnectionKepper.
type MockConnectionKepperMockRecorder struct {
	mock *MockConnectionKepper
}

// NewMockConnectionKepper creates a new mock instance.
func NewMockConnectionKepper(ctrl *gomock.Controller) *MockConnectionKepper {
	mock := &MockConnectionKepper{ctrl: ctrl}
	mock.recorder = &MockConnectionKepperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnectionKepper) EXPECT() *MockConnectionKepperMockRecorder {
	return m.recorder
}

// Discard mocks base method.
func (m *MockConnectionKepper) Discard(sh shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Discard", sh)
	ret0, _ := ret[0].(error)
	return ret0
}

// Discard indicates an expected call of Discard.
func (mr *MockConnectionKepperMockRecorder) Discard(sh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockConnectionKepper)(nil).Discard), sh)
}

// Put mocks base method.
func (m *MockConnectionKepper) Put(host shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", host)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockConnectionKepperMockRecorder) Put(host interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockConnectionKepper)(nil).Put), host)
}

// View mocks base method.
func (m *MockConnectionKepper) View() pool.Statistics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "View")
	ret0, _ := ret[0].(pool.Statistics)
	return ret0
}

// View indicates an expected call of View.
func (mr *MockConnectionKepperMockRecorder) View() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "View", reflect.TypeOf((*MockConnectionKepper)(nil).View))
}

// MockPool is a mock of Pool interface.
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool.
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance.
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// Connection mocks base method.
func (m *MockPool) Connection(clid uint, shardKey kr.ShardKey) (shard.Shard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connection", clid, shardKey)
	ret0, _ := ret[0].(shard.Shard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connection indicates an expected call of Connection.
func (mr *MockPoolMockRecorder) Connection(clid, shardKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connection", reflect.TypeOf((*MockPool)(nil).Connection), clid, shardKey)
}

// Discard mocks base method.
func (m *MockPool) Discard(sh shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Discard", sh)
	ret0, _ := ret[0].(error)
	return ret0
}

// Discard indicates an expected call of Discard.
func (mr *MockPoolMockRecorder) Discard(sh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockPool)(nil).Discard), sh)
}

// ForEach mocks base method.
func (m *MockPool) ForEach(cb func(shard.Shardinfo) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForEach", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForEach indicates an expected call of ForEach.
func (mr *MockPoolMockRecorder) ForEach(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEach", reflect.TypeOf((*MockPool)(nil).ForEach), cb)
}

// Put mocks base method.
func (m *MockPool) Put(host shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", host)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockPoolMockRecorder) Put(host interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockPool)(nil).Put), host)
}

// View mocks base method.
func (m *MockPool) View() pool.Statistics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "View")
	ret0, _ := ret[0].(pool.Statistics)
	return ret0
}

// View indicates an expected call of View.
func (mr *MockPoolMockRecorder) View() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "View", reflect.TypeOf((*MockPool)(nil).View))
}

// MockMultiShardPool is a mock of MultiShardPool interface.
type MockMultiShardPool struct {
	ctrl     *gomock.Controller
	recorder *MockMultiShardPoolMockRecorder
}

// MockMultiShardPoolMockRecorder is the mock recorder for MockMultiShardPool.
type MockMultiShardPoolMockRecorder struct {
	mock *MockMultiShardPool
}

// NewMockMultiShardPool creates a new mock instance.
func NewMockMultiShardPool(ctrl *gomock.Controller) *MockMultiShardPool {
	mock := &MockMultiShardPool{ctrl: ctrl}
	mock.recorder = &MockMultiShardPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMultiShardPool) EXPECT() *MockMultiShardPoolMockRecorder {
	return m.recorder
}

// Connection mocks base method.
func (m *MockMultiShardPool) Connection(clid uint, shardKey kr.ShardKey, host string) (shard.Shard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connection", clid, shardKey, host)
	ret0, _ := ret[0].(shard.Shard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connection indicates an expected call of Connection.
func (mr *MockMultiShardPoolMockRecorder) Connection(clid, shardKey, host interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connection", reflect.TypeOf((*MockMultiShardPool)(nil).Connection), clid, shardKey, host)
}

// Discard mocks base method.
func (m *MockMultiShardPool) Discard(sh shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Discard", sh)
	ret0, _ := ret[0].(error)
	return ret0
}

// Discard indicates an expected call of Discard.
func (mr *MockMultiShardPoolMockRecorder) Discard(sh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockMultiShardPool)(nil).Discard), sh)
}

// ForEach mocks base method.
func (m *MockMultiShardPool) ForEach(cb func(shard.Shardinfo) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForEach", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForEach indicates an expected call of ForEach.
func (mr *MockMultiShardPoolMockRecorder) ForEach(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEach", reflect.TypeOf((*MockMultiShardPool)(nil).ForEach), cb)
}

// ForEachPool mocks base method.
func (m *MockMultiShardPool) ForEachPool(cb func(pool.Pool) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForEachPool", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForEachPool indicates an expected call of ForEachPool.
func (mr *MockMultiShardPoolMockRecorder) ForEachPool(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEachPool", reflect.TypeOf((*MockMultiShardPool)(nil).ForEachPool), cb)
}

// Put mocks base method.
func (m *MockMultiShardPool) Put(host shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", host)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockMultiShardPoolMockRecorder) Put(host interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockMultiShardPool)(nil).Put), host)
}

// SetRule mocks base method.
func (m *MockMultiShardPool) SetRule(rule *config.BackendRule) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRule", rule)
}

// SetRule indicates an expected call of SetRule.
func (mr *MockMultiShardPoolMockRecorder) SetRule(rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRule", reflect.TypeOf((*MockMultiShardPool)(nil).SetRule), rule)
}

// View mocks base method.
func (m *MockMultiShardPool) View() pool.Statistics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "View")
	ret0, _ := ret[0].(pool.Statistics)
	return ret0
}

// View indicates an expected call of View.
func (mr *MockMultiShardPoolMockRecorder) View() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "View", reflect.TypeOf((*MockMultiShardPool)(nil).View))
}

// MockPoolIterator is a mock of PoolIterator interface.
type MockPoolIterator struct {
	ctrl     *gomock.Controller
	recorder *MockPoolIteratorMockRecorder
}

// MockPoolIteratorMockRecorder is the mock recorder for MockPoolIterator.
type MockPoolIteratorMockRecorder struct {
	mock *MockPoolIterator
}

// NewMockPoolIterator creates a new mock instance.
func NewMockPoolIterator(ctrl *gomock.Controller) *MockPoolIterator {
	mock := &MockPoolIterator{ctrl: ctrl}
	mock.recorder = &MockPoolIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPoolIterator) EXPECT() *MockPoolIteratorMockRecorder {
	return m.recorder
}

// ForEachPool mocks base method.
func (m *MockPoolIterator) ForEachPool(cb func(pool.Pool) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForEachPool", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForEachPool indicates an expected call of ForEachPool.
func (mr *MockPoolIteratorMockRecorder) ForEachPool(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEachPool", reflect.TypeOf((*MockPoolIterator)(nil).ForEachPool), cb)
}

// MockDBPool is a mock of DBPool interface.
type MockDBPool struct {
	ctrl     *gomock.Controller
	recorder *MockDBPoolMockRecorder
}

// MockDBPoolMockRecorder is the mock recorder for MockDBPool.
type MockDBPoolMockRecorder struct {
	mock *MockDBPool
}

// NewMockDBPool creates a new mock instance.
func NewMockDBPool(ctrl *gomock.Controller) *MockDBPool {
	mock := &MockDBPool{ctrl: ctrl}
	mock.recorder = &MockDBPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBPool) EXPECT() *MockDBPoolMockRecorder {
	return m.recorder
}

// Connection mocks base method.
func (m *MockDBPool) Connection(clid uint, shardKey kr.ShardKey, host string) (shard.Shard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connection", clid, shardKey, host)
	ret0, _ := ret[0].(shard.Shard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connection indicates an expected call of Connection.
func (mr *MockDBPoolMockRecorder) Connection(clid, shardKey, host interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connection", reflect.TypeOf((*MockDBPool)(nil).Connection), clid, shardKey, host)
}

// Discard mocks base method.
func (m *MockDBPool) Discard(sh shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Discard", sh)
	ret0, _ := ret[0].(error)
	return ret0
}

// Discard indicates an expected call of Discard.
func (mr *MockDBPoolMockRecorder) Discard(sh interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockDBPool)(nil).Discard), sh)
}

// ForEach mocks base method.
func (m *MockDBPool) ForEach(cb func(shard.Shardinfo) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForEach", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForEach indicates an expected call of ForEach.
func (mr *MockDBPoolMockRecorder) ForEach(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEach", reflect.TypeOf((*MockDBPool)(nil).ForEach), cb)
}

// ForEachPool mocks base method.
func (m *MockDBPool) ForEachPool(cb func(pool.Pool) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForEachPool", cb)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForEachPool indicates an expected call of ForEachPool.
func (mr *MockDBPoolMockRecorder) ForEachPool(cb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForEachPool", reflect.TypeOf((*MockDBPool)(nil).ForEachPool), cb)
}

// Put mocks base method.
func (m *MockDBPool) Put(host shard.Shard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", host)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockDBPoolMockRecorder) Put(host interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockDBPool)(nil).Put), host)
}

// SetRule mocks base method.
func (m *MockDBPool) SetRule(rule *config.BackendRule) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRule", rule)
}

// SetRule indicates an expected call of SetRule.
func (mr *MockDBPoolMockRecorder) SetRule(rule interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRule", reflect.TypeOf((*MockDBPool)(nil).SetRule), rule)
}

// SetShuffleHosts mocks base method.
func (m *MockDBPool) SetShuffleHosts(arg0 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetShuffleHosts", arg0)
}

// SetShuffleHosts indicates an expected call of SetShuffleHosts.
func (mr *MockDBPoolMockRecorder) SetShuffleHosts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetShuffleHosts", reflect.TypeOf((*MockDBPool)(nil).SetShuffleHosts), arg0)
}

// ShardMapping mocks base method.
func (m *MockDBPool) ShardMapping() map[string]*config.Shard {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShardMapping")
	ret0, _ := ret[0].(map[string]*config.Shard)
	return ret0
}

// ShardMapping indicates an expected call of ShardMapping.
func (mr *MockDBPoolMockRecorder) ShardMapping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShardMapping", reflect.TypeOf((*MockDBPool)(nil).ShardMapping))
}

// View mocks base method.
func (m *MockDBPool) View() pool.Statistics {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "View")
	ret0, _ := ret[0].(pool.Statistics)
	return ret0
}

// View indicates an expected call of View.
func (mr *MockDBPoolMockRecorder) View() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "View", reflect.TypeOf((*MockDBPool)(nil).View))
}
