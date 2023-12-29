// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/binarymatt/kayak/internal/store/models"
	mock "github.com/stretchr/testify/mock"

	store "github.com/binarymatt/kayak/internal/store"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

type Store_Expecter struct {
	mock *mock.Mock
}

func (_m *Store) EXPECT() *Store_Expecter {
	return &Store_Expecter{mock: &_m.Mock}
}

// AddRecords provides a mock function with given fields: ctx, topic, records
func (_m *Store) AddRecords(ctx context.Context, topic string, records ...*models.Record) error {
	_va := make([]interface{}, len(records))
	for _i := range records {
		_va[_i] = records[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, topic)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...*models.Record) error); ok {
		r0 = rf(ctx, topic, records...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_AddRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddRecords'
type Store_AddRecords_Call struct {
	*mock.Call
}

// AddRecords is a helper method to define mock.On call
//   - ctx context.Context
//   - topic string
//   - records ...*models.Record
func (_e *Store_Expecter) AddRecords(ctx interface{}, topic interface{}, records ...interface{}) *Store_AddRecords_Call {
	return &Store_AddRecords_Call{Call: _e.mock.On("AddRecords",
		append([]interface{}{ctx, topic}, records...)...)}
}

func (_c *Store_AddRecords_Call) Run(run func(ctx context.Context, topic string, records ...*models.Record)) *Store_AddRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*models.Record, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*models.Record)
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *Store_AddRecords_Call) Return(_a0 error) *Store_AddRecords_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_AddRecords_Call) RunAndReturn(run func(context.Context, string, ...*models.Record) error) *Store_AddRecords_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *Store) Close() {
	_m.Called()
}

// Store_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Store_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *Store_Expecter) Close() *Store_Close_Call {
	return &Store_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *Store_Close_Call) Run(run func()) *Store_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Store_Close_Call) Return() *Store_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *Store_Close_Call) RunAndReturn(run func()) *Store_Close_Call {
	_c.Call.Return(run)
	return _c
}

// CommitConsumerPosition provides a mock function with given fields: ctx, consumer
func (_m *Store) CommitConsumerPosition(ctx context.Context, consumer *models.Consumer) error {
	ret := _m.Called(ctx, consumer)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) error); ok {
		r0 = rf(ctx, consumer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_CommitConsumerPosition_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CommitConsumerPosition'
type Store_CommitConsumerPosition_Call struct {
	*mock.Call
}

// CommitConsumerPosition is a helper method to define mock.On call
//   - ctx context.Context
//   - consumer *models.Consumer
func (_e *Store_Expecter) CommitConsumerPosition(ctx interface{}, consumer interface{}) *Store_CommitConsumerPosition_Call {
	return &Store_CommitConsumerPosition_Call{Call: _e.mock.On("CommitConsumerPosition", ctx, consumer)}
}

func (_c *Store_CommitConsumerPosition_Call) Run(run func(ctx context.Context, consumer *models.Consumer)) *Store_CommitConsumerPosition_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Consumer))
	})
	return _c
}

func (_c *Store_CommitConsumerPosition_Call) Return(_a0 error) *Store_CommitConsumerPosition_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_CommitConsumerPosition_Call) RunAndReturn(run func(context.Context, *models.Consumer) error) *Store_CommitConsumerPosition_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTopic provides a mock function with given fields: ctx, topic
func (_m *Store) CreateTopic(ctx context.Context, topic *models.Topic) error {
	ret := _m.Called(ctx, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_CreateTopic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTopic'
type Store_CreateTopic_Call struct {
	*mock.Call
}

// CreateTopic is a helper method to define mock.On call
//   - ctx context.Context
//   - topic *models.Topic
func (_e *Store_Expecter) CreateTopic(ctx interface{}, topic interface{}) *Store_CreateTopic_Call {
	return &Store_CreateTopic_Call{Call: _e.mock.On("CreateTopic", ctx, topic)}
}

func (_c *Store_CreateTopic_Call) Run(run func(ctx context.Context, topic *models.Topic)) *Store_CreateTopic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Topic))
	})
	return _c
}

func (_c *Store_CreateTopic_Call) Return(_a0 error) *Store_CreateTopic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_CreateTopic_Call) RunAndReturn(run func(context.Context, *models.Topic) error) *Store_CreateTopic_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTopic provides a mock function with given fields: ctx, topic
func (_m *Store) DeleteTopic(ctx context.Context, topic *models.Topic) error {
	ret := _m.Called(ctx, topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Topic) error); ok {
		r0 = rf(ctx, topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_DeleteTopic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTopic'
type Store_DeleteTopic_Call struct {
	*mock.Call
}

// DeleteTopic is a helper method to define mock.On call
//   - ctx context.Context
//   - topic *models.Topic
func (_e *Store_Expecter) DeleteTopic(ctx interface{}, topic interface{}) *Store_DeleteTopic_Call {
	return &Store_DeleteTopic_Call{Call: _e.mock.On("DeleteTopic", ctx, topic)}
}

func (_c *Store_DeleteTopic_Call) Run(run func(ctx context.Context, topic *models.Topic)) *Store_DeleteTopic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Topic))
	})
	return _c
}

func (_c *Store_DeleteTopic_Call) Return(_a0 error) *Store_DeleteTopic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_DeleteTopic_Call) RunAndReturn(run func(context.Context, *models.Topic) error) *Store_DeleteTopic_Call {
	_c.Call.Return(run)
	return _c
}

// FetchRecord provides a mock function with given fields: ctx, consumer
func (_m *Store) FetchRecord(ctx context.Context, consumer *models.Consumer) (*models.Record, error) {
	ret := _m.Called(ctx, consumer)

	var r0 *models.Record
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) (*models.Record, error)); ok {
		return rf(ctx, consumer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) *models.Record); ok {
		r0 = rf(ctx, consumer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Record)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Consumer) error); ok {
		r1 = rf(ctx, consumer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_FetchRecord_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchRecord'
type Store_FetchRecord_Call struct {
	*mock.Call
}

// FetchRecord is a helper method to define mock.On call
//   - ctx context.Context
//   - consumer *models.Consumer
func (_e *Store_Expecter) FetchRecord(ctx interface{}, consumer interface{}) *Store_FetchRecord_Call {
	return &Store_FetchRecord_Call{Call: _e.mock.On("FetchRecord", ctx, consumer)}
}

func (_c *Store_FetchRecord_Call) Run(run func(ctx context.Context, consumer *models.Consumer)) *Store_FetchRecord_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Consumer))
	})
	return _c
}

func (_c *Store_FetchRecord_Call) Return(_a0 *models.Record, _a1 error) *Store_FetchRecord_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_FetchRecord_Call) RunAndReturn(run func(context.Context, *models.Consumer) (*models.Record, error)) *Store_FetchRecord_Call {
	_c.Call.Return(run)
	return _c
}

// GetConsumerLag provides a mock function with given fields: ctx, consumer
func (_m *Store) GetConsumerLag(ctx context.Context, consumer *models.Consumer) (int64, error) {
	ret := _m.Called(ctx, consumer)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) (int64, error)); ok {
		return rf(ctx, consumer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) int64); ok {
		r0 = rf(ctx, consumer)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Consumer) error); ok {
		r1 = rf(ctx, consumer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_GetConsumerLag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConsumerLag'
type Store_GetConsumerLag_Call struct {
	*mock.Call
}

// GetConsumerLag is a helper method to define mock.On call
//   - ctx context.Context
//   - consumer *models.Consumer
func (_e *Store_Expecter) GetConsumerLag(ctx interface{}, consumer interface{}) *Store_GetConsumerLag_Call {
	return &Store_GetConsumerLag_Call{Call: _e.mock.On("GetConsumerLag", ctx, consumer)}
}

func (_c *Store_GetConsumerLag_Call) Run(run func(ctx context.Context, consumer *models.Consumer)) *Store_GetConsumerLag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Consumer))
	})
	return _c
}

func (_c *Store_GetConsumerLag_Call) Return(_a0 int64, _a1 error) *Store_GetConsumerLag_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_GetConsumerLag_Call) RunAndReturn(run func(context.Context, *models.Consumer) (int64, error)) *Store_GetConsumerLag_Call {
	_c.Call.Return(run)
	return _c
}

// GetConsumerPosition provides a mock function with given fields: ctx, consumer
func (_m *Store) GetConsumerPosition(ctx context.Context, consumer *models.Consumer) (string, error) {
	ret := _m.Called(ctx, consumer)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) (string, error)); ok {
		return rf(ctx, consumer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) string); ok {
		r0 = rf(ctx, consumer)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Consumer) error); ok {
		r1 = rf(ctx, consumer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_GetConsumerPosition_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConsumerPosition'
type Store_GetConsumerPosition_Call struct {
	*mock.Call
}

// GetConsumerPosition is a helper method to define mock.On call
//   - ctx context.Context
//   - consumer *models.Consumer
func (_e *Store_Expecter) GetConsumerPosition(ctx interface{}, consumer interface{}) *Store_GetConsumerPosition_Call {
	return &Store_GetConsumerPosition_Call{Call: _e.mock.On("GetConsumerPosition", ctx, consumer)}
}

func (_c *Store_GetConsumerPosition_Call) Run(run func(ctx context.Context, consumer *models.Consumer)) *Store_GetConsumerPosition_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Consumer))
	})
	return _c
}

func (_c *Store_GetConsumerPosition_Call) Return(_a0 string, _a1 error) *Store_GetConsumerPosition_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_GetConsumerPosition_Call) RunAndReturn(run func(context.Context, *models.Consumer) (string, error)) *Store_GetConsumerPosition_Call {
	_c.Call.Return(run)
	return _c
}

// GetRecords provides a mock function with given fields: ctx, topic, start, limit
func (_m *Store) GetRecords(ctx context.Context, topic string, start string, limit int) ([]*models.Record, error) {
	ret := _m.Called(ctx, topic, start, limit)

	var r0 []*models.Record
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) ([]*models.Record, error)); ok {
		return rf(ctx, topic, start, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) []*models.Record); ok {
		r0 = rf(ctx, topic, start, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Record)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) error); ok {
		r1 = rf(ctx, topic, start, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_GetRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRecords'
type Store_GetRecords_Call struct {
	*mock.Call
}

// GetRecords is a helper method to define mock.On call
//   - ctx context.Context
//   - topic string
//   - start string
//   - limit int
func (_e *Store_Expecter) GetRecords(ctx interface{}, topic interface{}, start interface{}, limit interface{}) *Store_GetRecords_Call {
	return &Store_GetRecords_Call{Call: _e.mock.On("GetRecords", ctx, topic, start, limit)}
}

func (_c *Store_GetRecords_Call) Run(run func(ctx context.Context, topic string, start string, limit int)) *Store_GetRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(int))
	})
	return _c
}

func (_c *Store_GetRecords_Call) Return(_a0 []*models.Record, _a1 error) *Store_GetRecords_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_GetRecords_Call) RunAndReturn(run func(context.Context, string, string, int) ([]*models.Record, error)) *Store_GetRecords_Call {
	_c.Call.Return(run)
	return _c
}

// Impl provides a mock function with given fields:
func (_m *Store) Impl() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Store_Impl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Impl'
type Store_Impl_Call struct {
	*mock.Call
}

// Impl is a helper method to define mock.On call
func (_e *Store_Expecter) Impl() *Store_Impl_Call {
	return &Store_Impl_Call{Call: _e.mock.On("Impl")}
}

func (_c *Store_Impl_Call) Run(run func()) *Store_Impl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Store_Impl_Call) Return(_a0 interface{}) *Store_Impl_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_Impl_Call) RunAndReturn(run func() interface{}) *Store_Impl_Call {
	_c.Call.Return(run)
	return _c
}

// ListTopics provides a mock function with given fields: ctx
func (_m *Store) ListTopics(ctx context.Context) ([]string, error) {
	ret := _m.Called(ctx)

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []string); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_ListTopics_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListTopics'
type Store_ListTopics_Call struct {
	*mock.Call
}

// ListTopics is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Store_Expecter) ListTopics(ctx interface{}) *Store_ListTopics_Call {
	return &Store_ListTopics_Call{Call: _e.mock.On("ListTopics", ctx)}
}

func (_c *Store_ListTopics_Call) Run(run func(ctx context.Context)) *Store_ListTopics_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Store_ListTopics_Call) Return(_a0 []string, _a1 error) *Store_ListTopics_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_ListTopics_Call) RunAndReturn(run func(context.Context) ([]string, error)) *Store_ListTopics_Call {
	_c.Call.Return(run)
	return _c
}

// LoadMeta provides a mock function with given fields: ctx, topic
func (_m *Store) LoadMeta(ctx context.Context, topic string) (*models.Topic, error) {
	ret := _m.Called(ctx, topic)

	var r0 *models.Topic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Topic, error)); ok {
		return rf(ctx, topic)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Topic); ok {
		r0 = rf(ctx, topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Topic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, topic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_LoadMeta_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LoadMeta'
type Store_LoadMeta_Call struct {
	*mock.Call
}

// LoadMeta is a helper method to define mock.On call
//   - ctx context.Context
//   - topic string
func (_e *Store_Expecter) LoadMeta(ctx interface{}, topic interface{}) *Store_LoadMeta_Call {
	return &Store_LoadMeta_Call{Call: _e.mock.On("LoadMeta", ctx, topic)}
}

func (_c *Store_LoadMeta_Call) Run(run func(ctx context.Context, topic string)) *Store_LoadMeta_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Store_LoadMeta_Call) Return(_a0 *models.Topic, _a1 error) *Store_LoadMeta_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_LoadMeta_Call) RunAndReturn(run func(context.Context, string) (*models.Topic, error)) *Store_LoadMeta_Call {
	_c.Call.Return(run)
	return _c
}

// PruneOldRecords provides a mock function with given fields: ctx
func (_m *Store) PruneOldRecords(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_PruneOldRecords_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PruneOldRecords'
type Store_PruneOldRecords_Call struct {
	*mock.Call
}

// PruneOldRecords is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Store_Expecter) PruneOldRecords(ctx interface{}) *Store_PruneOldRecords_Call {
	return &Store_PruneOldRecords_Call{Call: _e.mock.On("PruneOldRecords", ctx)}
}

func (_c *Store_PruneOldRecords_Call) Run(run func(ctx context.Context)) *Store_PruneOldRecords_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Store_PruneOldRecords_Call) Return(_a0 error) *Store_PruneOldRecords_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_PruneOldRecords_Call) RunAndReturn(run func(context.Context) error) *Store_PruneOldRecords_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterConsumer provides a mock function with given fields: ctx, consumer
func (_m *Store) RegisterConsumer(ctx context.Context, consumer *models.Consumer) (*models.Consumer, error) {
	ret := _m.Called(ctx, consumer)

	var r0 *models.Consumer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) (*models.Consumer, error)); ok {
		return rf(ctx, consumer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Consumer) *models.Consumer); ok {
		r0 = rf(ctx, consumer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Consumer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Consumer) error); ok {
		r1 = rf(ctx, consumer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store_RegisterConsumer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterConsumer'
type Store_RegisterConsumer_Call struct {
	*mock.Call
}

// RegisterConsumer is a helper method to define mock.On call
//   - ctx context.Context
//   - consumer *models.Consumer
func (_e *Store_Expecter) RegisterConsumer(ctx interface{}, consumer interface{}) *Store_RegisterConsumer_Call {
	return &Store_RegisterConsumer_Call{Call: _e.mock.On("RegisterConsumer", ctx, consumer)}
}

func (_c *Store_RegisterConsumer_Call) Run(run func(ctx context.Context, consumer *models.Consumer)) *Store_RegisterConsumer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Consumer))
	})
	return _c
}

func (_c *Store_RegisterConsumer_Call) Return(_a0 *models.Consumer, _a1 error) *Store_RegisterConsumer_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Store_RegisterConsumer_Call) RunAndReturn(run func(context.Context, *models.Consumer) (*models.Consumer, error)) *Store_RegisterConsumer_Call {
	_c.Call.Return(run)
	return _c
}

// SnapshotItems provides a mock function with given fields:
func (_m *Store) SnapshotItems() <-chan store.DataItem {
	ret := _m.Called()

	var r0 <-chan store.DataItem
	if rf, ok := ret.Get(0).(func() <-chan store.DataItem); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan store.DataItem)
		}
	}

	return r0
}

// Store_SnapshotItems_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SnapshotItems'
type Store_SnapshotItems_Call struct {
	*mock.Call
}

// SnapshotItems is a helper method to define mock.On call
func (_e *Store_Expecter) SnapshotItems() *Store_SnapshotItems_Call {
	return &Store_SnapshotItems_Call{Call: _e.mock.On("SnapshotItems")}
}

func (_c *Store_SnapshotItems_Call) Run(run func()) *Store_SnapshotItems_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Store_SnapshotItems_Call) Return(_a0 <-chan store.DataItem) *Store_SnapshotItems_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_SnapshotItems_Call) RunAndReturn(run func() <-chan store.DataItem) *Store_SnapshotItems_Call {
	_c.Call.Return(run)
	return _c
}

// Stats provides a mock function with given fields:
func (_m *Store) Stats() map[string]*models.Topic {
	ret := _m.Called()

	var r0 map[string]*models.Topic
	if rf, ok := ret.Get(0).(func() map[string]*models.Topic); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*models.Topic)
		}
	}

	return r0
}

// Store_Stats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stats'
type Store_Stats_Call struct {
	*mock.Call
}

// Stats is a helper method to define mock.On call
func (_e *Store_Expecter) Stats() *Store_Stats_Call {
	return &Store_Stats_Call{Call: _e.mock.On("Stats")}
}

func (_c *Store_Stats_Call) Run(run func()) *Store_Stats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Store_Stats_Call) Return(_a0 map[string]*models.Topic) *Store_Stats_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_Stats_Call) RunAndReturn(run func() map[string]*models.Topic) *Store_Stats_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStore(t mockConstructorTestingTNewStore) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
