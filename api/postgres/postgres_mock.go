package postgres

type MockPostgres struct {
	QueryRowFn      func(string, ...interface{}) Scanner
	QueryRowInvoked bool
}

func (m *MockPostgres) QueryRow(query string, args ...interface{}) Scanner {
	m.QueryRowInvoked = true
	return m.QueryRowFn(query, args...)
}

type MockScanner struct {
	ScanFn      func(dest ...interface{}) error
	ScanInvoked bool
}

func (s *MockScanner) Scan(dest ...interface{}) error {
	s.ScanInvoked = true
	return s.ScanFn(dest...)
}
