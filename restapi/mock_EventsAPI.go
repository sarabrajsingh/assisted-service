// Code generated by mockery v1.0.0. DO NOT EDIT.

package restapi

import (
	context "context"

	events "github.com/filanov/bm-inventory/restapi/operations/events"
	middleware "github.com/go-openapi/runtime/middleware"

	mock "github.com/stretchr/testify/mock"
)

// MockEventsAPI is an autogenerated mock type for the EventsAPI type
type MockEventsAPI struct {
	mock.Mock
}

// ListEvents provides a mock function with given fields: ctx, params
func (_m *MockEventsAPI) ListEvents(ctx context.Context, params events.ListEventsParams) middleware.Responder {
	ret := _m.Called(ctx, params)

	var r0 middleware.Responder
	if rf, ok := ret.Get(0).(func(context.Context, events.ListEventsParams) middleware.Responder); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(middleware.Responder)
		}
	}

	return r0
}
