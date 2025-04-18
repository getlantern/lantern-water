// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/getlantern/lantern-water/downloader (interfaces: WASMDownloader,torrentClient,torrentInfo)
//
// Generated by this command:
//
//	mockgen -package=downloader -destination=mocks_test.go . WASMDownloader,torrentClient,torrentInfo
//

// Package downloader is a generated GoMock package.
package downloader

import (
	context "context"
	io "io"
	reflect "reflect"

	events "github.com/anacrolix/chansync/events"
	torrent "github.com/anacrolix/torrent"
	gomock "go.uber.org/mock/gomock"
)

// MockWASMDownloader is a mock of WASMDownloader interface.
type MockWASMDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockWASMDownloaderMockRecorder
	isgomock struct{}
}

// MockWASMDownloaderMockRecorder is the mock recorder for MockWASMDownloader.
type MockWASMDownloaderMockRecorder struct {
	mock *MockWASMDownloader
}

// NewMockWASMDownloader creates a new mock instance.
func NewMockWASMDownloader(ctrl *gomock.Controller) *MockWASMDownloader {
	mock := &MockWASMDownloader{ctrl: ctrl}
	mock.recorder = &MockWASMDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWASMDownloader) EXPECT() *MockWASMDownloaderMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockWASMDownloader) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockWASMDownloaderMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWASMDownloader)(nil).Close))
}

// DownloadWASM mocks base method.
func (m *MockWASMDownloader) DownloadWASM(arg0 context.Context, arg1 io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadWASM", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DownloadWASM indicates an expected call of DownloadWASM.
func (mr *MockWASMDownloaderMockRecorder) DownloadWASM(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadWASM", reflect.TypeOf((*MockWASMDownloader)(nil).DownloadWASM), arg0, arg1)
}

// MocktorrentClient is a mock of torrentClient interface.
type MocktorrentClient struct {
	ctrl     *gomock.Controller
	recorder *MocktorrentClientMockRecorder
	isgomock struct{}
}

// MocktorrentClientMockRecorder is the mock recorder for MocktorrentClient.
type MocktorrentClientMockRecorder struct {
	mock *MocktorrentClient
}

// NewMocktorrentClient creates a new mock instance.
func NewMocktorrentClient(ctrl *gomock.Controller) *MocktorrentClient {
	mock := &MocktorrentClient{ctrl: ctrl}
	mock.recorder = &MocktorrentClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktorrentClient) EXPECT() *MocktorrentClientMockRecorder {
	return m.recorder
}

// AddMagnet mocks base method.
func (m *MocktorrentClient) AddMagnet(arg0 string) (torrentInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMagnet", arg0)
	ret0, _ := ret[0].(torrentInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMagnet indicates an expected call of AddMagnet.
func (mr *MocktorrentClientMockRecorder) AddMagnet(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMagnet", reflect.TypeOf((*MocktorrentClient)(nil).AddMagnet), arg0)
}

// Close mocks base method.
func (m *MocktorrentClient) Close() []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].([]error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MocktorrentClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MocktorrentClient)(nil).Close))
}

// MocktorrentInfo is a mock of torrentInfo interface.
type MocktorrentInfo struct {
	ctrl     *gomock.Controller
	recorder *MocktorrentInfoMockRecorder
	isgomock struct{}
}

// MocktorrentInfoMockRecorder is the mock recorder for MocktorrentInfo.
type MocktorrentInfoMockRecorder struct {
	mock *MocktorrentInfo
}

// NewMocktorrentInfo creates a new mock instance.
func NewMocktorrentInfo(ctrl *gomock.Controller) *MocktorrentInfo {
	mock := &MocktorrentInfo{ctrl: ctrl}
	mock.recorder = &MocktorrentInfoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktorrentInfo) EXPECT() *MocktorrentInfoMockRecorder {
	return m.recorder
}

// GotInfo mocks base method.
func (m *MocktorrentInfo) GotInfo() events.Done {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GotInfo")
	ret0, _ := ret[0].(events.Done)
	return ret0
}

// GotInfo indicates an expected call of GotInfo.
func (mr *MocktorrentInfoMockRecorder) GotInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GotInfo", reflect.TypeOf((*MocktorrentInfo)(nil).GotInfo))
}

// NewReader mocks base method.
func (m *MocktorrentInfo) NewReader() torrent.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewReader")
	ret0, _ := ret[0].(torrent.Reader)
	return ret0
}

// NewReader indicates an expected call of NewReader.
func (mr *MocktorrentInfoMockRecorder) NewReader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewReader", reflect.TypeOf((*MocktorrentInfo)(nil).NewReader))
}
