// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/hospital/hospital.go

// Package mock_hospital_client is a generated GoMock package.
package mock_hospital_client

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	hospital "github.com/synthia-telemed/backend-api/pkg/hospital"
)

// MockSystemClient is a mock of SystemClient interface.
type MockSystemClient struct {
	ctrl     *gomock.Controller
	recorder *MockSystemClientMockRecorder
}

// MockSystemClientMockRecorder is the mock recorder for MockSystemClient.
type MockSystemClientMockRecorder struct {
	mock *MockSystemClient
}

// NewMockSystemClient creates a new mock instance.
func NewMockSystemClient(ctrl *gomock.Controller) *MockSystemClient {
	mock := &MockSystemClient{ctrl: ctrl}
	mock.recorder = &MockSystemClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSystemClient) EXPECT() *MockSystemClientMockRecorder {
	return m.recorder
}

// AssertDoctorCredential mocks base method.
func (m *MockSystemClient) AssertDoctorCredential(ctx context.Context, username, password string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssertDoctorCredential", ctx, username, password)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssertDoctorCredential indicates an expected call of AssertDoctorCredential.
func (mr *MockSystemClientMockRecorder) AssertDoctorCredential(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssertDoctorCredential", reflect.TypeOf((*MockSystemClient)(nil).AssertDoctorCredential), ctx, username, password)
}

// CategorizeAppointmentByStatus mocks base method.
func (m *MockSystemClient) CategorizeAppointmentByStatus(apps []*hospital.AppointmentOverview) *hospital.CategorizedAppointment {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CategorizeAppointmentByStatus", apps)
	ret0, _ := ret[0].(*hospital.CategorizedAppointment)
	return ret0
}

// CategorizeAppointmentByStatus indicates an expected call of CategorizeAppointmentByStatus.
func (mr *MockSystemClientMockRecorder) CategorizeAppointmentByStatus(apps interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CategorizeAppointmentByStatus", reflect.TypeOf((*MockSystemClient)(nil).CategorizeAppointmentByStatus), apps)
}

// CountAppointmentsWithFilters mocks base method.
func (m *MockSystemClient) CountAppointmentsWithFilters(ctx context.Context, filters *hospital.ListAppointmentsFilters) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAppointmentsWithFilters", ctx, filters)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAppointmentsWithFilters indicates an expected call of CountAppointmentsWithFilters.
func (mr *MockSystemClientMockRecorder) CountAppointmentsWithFilters(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAppointmentsWithFilters", reflect.TypeOf((*MockSystemClient)(nil).CountAppointmentsWithFilters), ctx, filters)
}

// FindAppointmentByID mocks base method.
func (m *MockSystemClient) FindAppointmentByID(ctx context.Context, appointmentID int) (*hospital.Appointment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAppointmentByID", ctx, appointmentID)
	ret0, _ := ret[0].(*hospital.Appointment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAppointmentByID indicates an expected call of FindAppointmentByID.
func (mr *MockSystemClientMockRecorder) FindAppointmentByID(ctx, appointmentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAppointmentByID", reflect.TypeOf((*MockSystemClient)(nil).FindAppointmentByID), ctx, appointmentID)
}

// FindDoctorAppointmentByID mocks base method.
func (m *MockSystemClient) FindDoctorAppointmentByID(ctx context.Context, appointmentID int) (*hospital.DoctorAppointment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDoctorAppointmentByID", ctx, appointmentID)
	ret0, _ := ret[0].(*hospital.DoctorAppointment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDoctorAppointmentByID indicates an expected call of FindDoctorAppointmentByID.
func (mr *MockSystemClientMockRecorder) FindDoctorAppointmentByID(ctx, appointmentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDoctorAppointmentByID", reflect.TypeOf((*MockSystemClient)(nil).FindDoctorAppointmentByID), ctx, appointmentID)
}

// FindDoctorByUsername mocks base method.
func (m *MockSystemClient) FindDoctorByUsername(ctx context.Context, username string) (*hospital.Doctor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindDoctorByUsername", ctx, username)
	ret0, _ := ret[0].(*hospital.Doctor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindDoctorByUsername indicates an expected call of FindDoctorByUsername.
func (mr *MockSystemClientMockRecorder) FindDoctorByUsername(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindDoctorByUsername", reflect.TypeOf((*MockSystemClient)(nil).FindDoctorByUsername), ctx, username)
}

// FindInvoiceByID mocks base method.
func (m *MockSystemClient) FindInvoiceByID(ctx context.Context, id int) (*hospital.InvoiceOverview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindInvoiceByID", ctx, id)
	ret0, _ := ret[0].(*hospital.InvoiceOverview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindInvoiceByID indicates an expected call of FindInvoiceByID.
func (mr *MockSystemClientMockRecorder) FindInvoiceByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindInvoiceByID", reflect.TypeOf((*MockSystemClient)(nil).FindInvoiceByID), ctx, id)
}

// FindPatientByGovCredential mocks base method.
func (m *MockSystemClient) FindPatientByGovCredential(ctx context.Context, cred string) (*hospital.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPatientByGovCredential", ctx, cred)
	ret0, _ := ret[0].(*hospital.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPatientByGovCredential indicates an expected call of FindPatientByGovCredential.
func (mr *MockSystemClientMockRecorder) FindPatientByGovCredential(ctx, cred interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPatientByGovCredential", reflect.TypeOf((*MockSystemClient)(nil).FindPatientByGovCredential), ctx, cred)
}

// FindPatientByID mocks base method.
func (m *MockSystemClient) FindPatientByID(ctx context.Context, id string) (*hospital.Patient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPatientByID", ctx, id)
	ret0, _ := ret[0].(*hospital.Patient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPatientByID indicates an expected call of FindPatientByID.
func (mr *MockSystemClientMockRecorder) FindPatientByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPatientByID", reflect.TypeOf((*MockSystemClient)(nil).FindPatientByID), ctx, id)
}

// ListAppointmentsByDoctorID mocks base method.
func (m *MockSystemClient) ListAppointmentsByDoctorID(ctx context.Context, doctorID string, date time.Time) ([]*hospital.AppointmentOverview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAppointmentsByDoctorID", ctx, doctorID, date)
	ret0, _ := ret[0].([]*hospital.AppointmentOverview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAppointmentsByDoctorID indicates an expected call of ListAppointmentsByDoctorID.
func (mr *MockSystemClientMockRecorder) ListAppointmentsByDoctorID(ctx, doctorID, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAppointmentsByDoctorID", reflect.TypeOf((*MockSystemClient)(nil).ListAppointmentsByDoctorID), ctx, doctorID, date)
}

// ListAppointmentsByPatientID mocks base method.
func (m *MockSystemClient) ListAppointmentsByPatientID(ctx context.Context, patientID string, since time.Time) ([]*hospital.AppointmentOverview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAppointmentsByPatientID", ctx, patientID, since)
	ret0, _ := ret[0].([]*hospital.AppointmentOverview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAppointmentsByPatientID indicates an expected call of ListAppointmentsByPatientID.
func (mr *MockSystemClientMockRecorder) ListAppointmentsByPatientID(ctx, patientID, since interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAppointmentsByPatientID", reflect.TypeOf((*MockSystemClient)(nil).ListAppointmentsByPatientID), ctx, patientID, since)
}

// ListAppointmentsWithFilters mocks base method.
func (m *MockSystemClient) ListAppointmentsWithFilters(ctx context.Context, filters *hospital.ListAppointmentsFilters, take, skip int) ([]*hospital.AppointmentOverview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAppointmentsWithFilters", ctx, filters, take, skip)
	ret0, _ := ret[0].([]*hospital.AppointmentOverview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAppointmentsWithFilters indicates an expected call of ListAppointmentsWithFilters.
func (mr *MockSystemClientMockRecorder) ListAppointmentsWithFilters(ctx, filters, take, skip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAppointmentsWithFilters", reflect.TypeOf((*MockSystemClient)(nil).ListAppointmentsWithFilters), ctx, filters, take, skip)
}

// PaidInvoice mocks base method.
func (m *MockSystemClient) PaidInvoice(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaidInvoice", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// PaidInvoice indicates an expected call of PaidInvoice.
func (mr *MockSystemClientMockRecorder) PaidInvoice(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaidInvoice", reflect.TypeOf((*MockSystemClient)(nil).PaidInvoice), ctx, id)
}

// SetAppointmentStatus mocks base method.
func (m *MockSystemClient) SetAppointmentStatus(ctx context.Context, appointmentID int, status hospital.SettableAppointmentStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAppointmentStatus", ctx, appointmentID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAppointmentStatus indicates an expected call of SetAppointmentStatus.
func (mr *MockSystemClientMockRecorder) SetAppointmentStatus(ctx, appointmentID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAppointmentStatus", reflect.TypeOf((*MockSystemClient)(nil).SetAppointmentStatus), ctx, appointmentID, status)
}
