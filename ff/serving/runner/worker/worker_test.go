package worker

import (
	"errors"
	runner "github.com/featureform/serving/runner"
	"testing"
)

type MockRunner struct{}

type MockCompletionWatcher struct{}

func (m *MockRunner) Run() (runner.CompletionWatcher, error) {
	return &MockCompletionWatcher{}, nil
}

func (m *MockCompletionWatcher) Complete() bool {
	return false
}

func (m *MockCompletionWatcher) String() string {
	return ""
}

func (m *MockCompletionWatcher) Wait() error {
	return nil
}

func (m *MockCompletionWatcher) Err() error {
	return nil
}

type RunnerWithFailingWatcher struct{}

func (r *RunnerWithFailingWatcher) Run() (runner.CompletionWatcher, error) {
	return &FailingWatcher{}, nil
}

type FailingWatcher struct{}

func (f *FailingWatcher) Complete() bool {
	return false
}
func (f *FailingWatcher) String() string {
	return ""
}
func (f *FailingWatcher) Wait() error {
	return errors.New("Run failed")
}
func (f *FailingWatcher) Err() error {
	return errors.New("Run failed")
}

type FailingRunner struct{}

func (f *FailingRunner) Run() (runner.CompletionWatcher, error) {
	return nil, errors.New("Failed to run runner")
}

func registerMockRunnerFactoryFailingWatcher() error {
	mockRunnerFailingWatcher := &RunnerWithFailingWatcher{}
	mockFactory := func(config runner.Config) (runner.Runner, error) {
		return mockRunnerFailingWatcher, nil
	}
	if err := runner.RegisterFactory("test", mockFactory); err != nil {
		return err
	}
	return nil
}

func registerMockFailRunnerFactory() error {
	failRunner := &FailingRunner{}
	failRunnerFactory := func(config runner.Config) (runner.Runner, error) {
		return failRunner, nil
	}
	if err := runner.RegisterFactory("test", failRunnerFactory); err != nil {
		return err
	}
	return nil
}

func registerMockRunnerFactory() error {
	mockRunner := &MockRunner{}
	mockFactory := func(config runner.Config) (runner.Runner, error) {
		return mockRunner, nil
	}
	if err := runner.RegisterFactory("test", mockFactory); err != nil {
		return err
	}
	return nil
}

func TestBasicRunner(t *testing.T) {
	runner.ResetFactoryMap()
	if err := registerMockRunnerFactory(); err != nil {
		t.Fatalf("Error registering mock runner factory: %v", err)
	}
	config := runner.Config{}
	t.Setenv("CONFIG", string(config))
	t.Setenv("NAME", "test")
	if err := CreateAndRun(); err != nil {
		t.Fatalf("Error running mock runner: %v", err)
	}
}

func TestRunnerNoConfig(t *testing.T) {
	runner.ResetFactoryMap()
	if err := registerMockRunnerFactory(); err != nil {
		t.Fatalf("Error registering mock runner factory: %v", err)
	}
	t.Setenv("NAME", "test")
	if err := CreateAndRun(); err == nil {
		t.Fatalf("Failed to call error on missing config envvar")
	}
}

func TestRunnerNoName(t *testing.T) {
	runner.ResetFactoryMap()
	if err := registerMockRunnerFactory(); err != nil {
		t.Fatalf("Error registering mock runner factory: %v", err)
	}
	config := runner.Config{}
	t.Setenv("CONFIG", string(config))
	if err := CreateAndRun(); err == nil {
		t.Fatalf("Failed to call error on missing name envvar")
	}
}

func TestRunnerNoFactory(t *testing.T) {
	runner.ResetFactoryMap()
	config := runner.Config{}
	t.Setenv("CONFIG", string(config))
	t.Setenv("NAME", "test")
	if err := CreateAndRun(); err == nil {
		t.Fatalf("Failed to call error on missing runner factory")
	}
}

func TestRunnerCreateFail(t *testing.T) {
	runner.ResetFactoryMap()
	if err := registerMockFailRunnerFactory(); err != nil {
		t.Fatalf("Error registering mock runner factory: %v", err)
	}
	config := runner.Config{}
	t.Setenv("CONFIG", string(config))
	t.Setenv("NAME", "test")
	if err := CreateAndRun(); err == nil {
		t.Fatalf("Broken runner doesn't fail when run")
	}
}

func TestRunnerRunFail(t *testing.T) {
	runner.ResetFactoryMap()
	if err := registerMockRunnerFactoryFailingWatcher(); err != nil {
		t.Fatalf("Error registering mock runner factory: %v", err)
	}
	config := runner.Config{}
	t.Setenv("CONFIG", string(config))
	t.Setenv("NAME", "test")
	if err := CreateAndRun(); err == nil {
		t.Fatalf("Broken watcher does not return error")
	}
}