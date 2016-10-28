package resources

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Times/go-fthealth"
	"github.com/stretchr/testify/assert"
)

func TestHealthy(t *testing.T) {
	mockDb := new(MockDB)
	mockTx := new(MockTX)

	mockDb.On("Open").Return(mockTx, nil)
	mockTx.On("Ping").Return(nil)

	req := httptest.NewRequest("GET", "http://nothing/__health", nil)
	w := httptest.NewRecorder()

	Health(mockDb)(w, req)

	if w.Code != 200 {
		t.Fatal("Should return 200!")
	}

	health := fthealth.HealthResult{}
	err := json.NewDecoder(w.Body).Decode(&health)
	if err != nil {
		t.Fatal("Should return valid json!")
	}

	assert.NotEmpty(t, health.Name, "Should have a non-empty name")
	assert.NotEmpty(t, health.Description, "Should have a non-empty description")
	assert.NotEmpty(t, health.SchemaVersion, "Should have a non-empty schema version")
	assert.True(t, health.Ok, "Expect it's ok")

	assert.Len(t, health.Checks, 1, "Only one health check currently")
	check := health.Checks[0]

	assert.NotEmpty(t, check.Name, "Should have a non-empty name")
	assert.NotEmpty(t, check.PanicGuide, "Should have a non-empty panic guide")
	assert.Equal(t, uint8(1), check.Severity, "Severity 1")
	assert.NotEmpty(t, check.BusinessImpact, "Should have a non-empty business impact")
	assert.NotEmpty(t, check.TechnicalSummary, "Should have a non-empty technical summary")

	assert.True(t, check.Ok, "Expect it's ok")
}

func TestUnhealthy(t *testing.T) {
	mockDb := new(MockDB)
	mockTx := new(MockTX)

	mockDb.On("Open").Return(mockTx, nil)
	mockTx.On("Ping").Return(errors.New("we ain't looking too good"))

	req := httptest.NewRequest("GET", "http://nothing/__health", nil)
	w := httptest.NewRecorder()

	Health(mockDb)(w, req)

	if w.Code != 200 {
		t.Fatal("Should always return 200!")
	}

	health := fthealth.HealthResult{}
	err := json.NewDecoder(w.Body).Decode(&health)
	if err != nil {
		t.Fatal("Should return valid json!")
	}

	assert.NotEmpty(t, health.Name, "Should have a non-empty name")
	assert.NotEmpty(t, health.Description, "Should have a non-empty description")
	assert.NotEmpty(t, health.SchemaVersion, "Should have a non-empty schema version")
	assert.False(t, health.Ok, "Expect it's ok")

	assert.Len(t, health.Checks, 1, "Only one health check currently")
	check := health.Checks[0]

	assert.NotEmpty(t, check.Name, "Should have a non-empty name")
	assert.NotEmpty(t, check.PanicGuide, "Should have a non-empty panic guide")
	assert.Equal(t, uint8(1), check.Severity, "Severity 1")
	assert.NotEmpty(t, check.BusinessImpact, "Should have a non-empty business impact")
	assert.NotEmpty(t, check.TechnicalSummary, "Should have a non-empty technical summary")

	assert.False(t, check.Ok, "Expect it's not ok")
}

func TestWorkingGTG(t *testing.T) {
	mockDb := new(MockDB)
	mockTx := new(MockTX)

	mockDb.On("Open").Return(mockTx, nil)
	mockTx.On("Ping").Return(nil)

	req := httptest.NewRequest("GET", "http://nothing/at/__gtg", nil)
	w := httptest.NewRecorder()

	GTG(mockDb)(w, req)

	if w.Code != 200 {
		t.Fatal("Should return 200!")
	}

	t.Log("Return 200 as expected.")
}

func TestFailingGTG(t *testing.T) {
	mockDb := new(MockDB)
	mockTx := new(MockTX)

	mockDb.On("Open").Return(mockTx, nil)
	mockTx.On("Ping").Return(errors.New("omg we are not gtg"))

	req := httptest.NewRequest("GET", "http://nothing/at/__gtg", nil)
	w := httptest.NewRecorder()

	GTG(mockDb)(w, req)

	if w.Code != 500 {
		t.Fatal("Should return 500!")
	}

	t.Log("Return 500 as expected.")
}

func TestFailingDBGTG(t *testing.T) {
	mockDb := new(MockDB)

	mockDb.On("Open").Return(nil, errors.New("omg we are not gtg"))

	req := httptest.NewRequest("GET", "http://nothing/at/__gtg", nil)
	w := httptest.NewRecorder()

	GTG(mockDb)(w, req)

	if w.Code != 500 {
		t.Fatal("Should return 500!")
	}

	t.Log("Return 500 as expected.")
}
