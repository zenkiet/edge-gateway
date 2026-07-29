package static

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"edge-gateway/internal/pkg/version"
)

type VersionEntry struct {
	ID      string
	DirPath string
	ModTime time.Time
}

type Manager struct {
	baseDir     string
	mappingPath string
	logger      *slog.Logger

	mu       sync.RWMutex
	versions map[string]*VersionEntry
	mapping  map[string]string
	latest   string

	stopWatch chan struct{}
}

func NewManager(baseDir string, logger *slog.Logger) *Manager {
	m := &Manager{
		baseDir:     baseDir,
		mappingPath: filepath.Join(baseDir, "mapping.json"),
		logger:      logger,
		versions:    make(map[string]*VersionEntry),
		mapping:     make(map[string]string),
		stopWatch:   make(chan struct{}),
	}

	m.ScanVersions()
	return m
}

func (m *Manager) ScanVersions() {
	versionsDir := filepath.Join(m.baseDir, "versions")
	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		m.logger.Error("Failed to scan versions", "error", err)
		return
	}

	newVersions := make(map[string]*VersionEntry)
	var latest string

	for _, de := range entries {
		if !de.IsDir() {
			continue
		}
		info, err := de.Info()
		if err != nil {
			continue
		}

		name := de.Name()
		newVersions[name] = &VersionEntry{
			ID:      name,
			DirPath: filepath.Join(versionsDir, name),
			ModTime: info.ModTime(),
		}

		if latest == "" || compareVersions(name, latest) > 0 {
			latest = name
		}
	}

	m.mu.Lock()
	m.versions = newVersions
	m.latest = latest
	m.mu.Unlock()

	m.logger.Info("Versions scanned", "count", len(newVersions), "latest", latest)
}

func (m *Manager) ReloadMapping() {
	data, err := os.ReadFile(m.mappingPath)
	if err != nil {
		m.logger.Warn("Mapping file not found or unreadable, using empty map")
		return
	}

	var newMapping map[string]string
	if err := json.Unmarshal(data, &newMapping); err != nil {
		m.logger.Error("Failed to parse mapping.json", "error", err)
		return
	}

	m.mu.Lock()
	m.mapping = newMapping
	m.mu.Unlock()
	m.logger.Info("Mapping reloaded", "entries", len(newMapping))
}

func (m *Manager) ResolveBundleID(posVersion string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	posVersion = strings.TrimSpace(posVersion)
	if bundleID, ok := m.mapping[posVersion]; ok {
		return strings.TrimPrefix(bundleID, "v")
	}
	return ""
}

func (m *Manager) GetVersion(bundleID string) (*VersionEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	id := strings.TrimPrefix(strings.TrimSpace(bundleID), "v")
	if id == "" || id == "latest" {
		id = m.latest
	}

	v, ok := m.versions[id]
	if !ok {
		v, ok = m.versions[m.latest]
	}
	return v, ok
}

func (m *Manager) IsLatest(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return strings.TrimPrefix(id, "v") == strings.TrimPrefix(m.latest, "v")
}

func (m *Manager) StartWatcher() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		var lastMod time.Time

		for {
			select {
			case <-m.stopWatch:
				return
			case <-ticker.C:
				info, err := os.Stat(m.mappingPath)
				if err == nil && info.ModTime().After(lastMod) {
					lastMod = info.ModTime()
					m.ReloadMapping()
				}
			}
		}
	}()
}

func (m *Manager) StopWatcher() {
	close(m.stopWatch)
}

func compareVersions(a, b string) int {
	vA, okA := version.ParseVersion(a)
	vB, okB := version.ParseVersion(b)

	if okA && okB {
		return vA.Compare(vB)
	}

	return strings.Compare(a, b)
}
