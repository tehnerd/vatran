package lb

import (
	"sync"

	"github.com/tehnerd/vatran/go/katran"
)

// Manager is a singleton manager for the Katran LoadBalancer instance.
// It provides thread-safe access to create, get, and close the load balancer.
type Manager struct {
	mu          sync.RWMutex
	lb          *katran.LoadBalancer
	config      *katran.Config
	initialized bool
	ready       bool
}

var (
	instance *Manager
	once     sync.Once
)

// GetManager returns the singleton Manager instance.
//
// Returns the Manager singleton.
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{}
	})
	return instance
}

// Create creates a new LoadBalancer instance with the provided configuration.
// If a LoadBalancer already exists, it returns an error.
//
// Parameters:
//   - cfg: Configuration for the load balancer.
//
// Returns an error if creation fails or if LB is already initialized.
func (m *Manager) Create(cfg *katran.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.initialized {
		return &katran.KatranError{
			Code:    katran.ErrAlreadyExists,
			Message: "load balancer is already initialized",
		}
	}

	lb, err := katran.New(cfg)
	if err != nil {
		return err
	}

	m.lb = lb
	m.config = cfg
	m.initialized = true
	m.ready = false
	return nil
}

// Get returns the LoadBalancer instance if it exists.
//
// Returns the LoadBalancer and a boolean indicating if it was found.
func (m *Manager) Get() (*katran.LoadBalancer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lb, m.initialized
}

// GetConfig returns the configuration used to create the LoadBalancer.
//
// Returns the Config and a boolean indicating if it was found.
func (m *Manager) GetConfig() (*katran.Config, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config, m.initialized
}

// Close closes the LoadBalancer instance and releases all resources.
//
// Returns an error if the LB is not initialized or if closing fails.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.initialized {
		return &katran.KatranError{
			Code:    katran.ErrNotFound,
			Message: "load balancer is not initialized",
		}
	}

	err := m.lb.Close()
	m.lb = nil
	m.config = nil
	m.initialized = false
	m.ready = false
	return err
}

// Status returns the current status of the LoadBalancer.
//
// Returns:
//   - initialized: Whether the LB instance has been created.
//   - ready: Whether BPF programs are loaded and attached.
func (m *Manager) Status() (initialized, ready bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized, m.ready
}

// SetReady marks the LoadBalancer as ready (BPF programs loaded and attached).
//
// Parameters:
//   - ready: Whether the LB is ready.
func (m *Manager) SetReady(ready bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ready = ready
}

// LoadBPFProgs loads BPF programs into the kernel.
//
// Returns an error if the LB is not initialized or if loading fails.
func (m *Manager) LoadBPFProgs() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.initialized {
		return &katran.KatranError{
			Code:    katran.ErrNotFound,
			Message: "load balancer is not initialized",
		}
	}

	return m.lb.LoadBPFProgs()
}

// AttachBPFProgs attaches loaded BPF programs to network interfaces.
//
// Returns an error if the LB is not initialized or if attaching fails.
func (m *Manager) AttachBPFProgs() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.initialized {
		return &katran.KatranError{
			Code:    katran.ErrNotFound,
			Message: "load balancer is not initialized",
		}
	}

	err := m.lb.AttachBPFProgs()
	if err == nil {
		m.ready = true
	}
	return err
}

// ReloadBalancerProg reloads the balancer BPF program at runtime.
//
// Parameters:
//   - path: Path to the new BPF program file.
//   - cfg: Optional new configuration. Pass nil to keep current config.
//
// Returns an error if the LB is not initialized or if reloading fails.
func (m *Manager) ReloadBalancerProg(path string, cfg *katran.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.initialized {
		return &katran.KatranError{
			Code:    katran.ErrNotFound,
			Message: "load balancer is not initialized",
		}
	}

	return m.lb.ReloadBalancerProg(path, cfg)
}
