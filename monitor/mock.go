package monitor

// MockMonitor is ...
type MockMonitor struct{}

// MonitorTx is ...
func (m *MockMonitor) MonitorTx(ReqID int) bool {
	return true
}
