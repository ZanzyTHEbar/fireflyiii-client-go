package importers

import (
	"context"
	"time"

	"github.com/ZanzyTHEbar/errbuilder-go"
)

// ImporterConfig represents the configuration for an importer
type ImporterConfig struct {
	// Name of the importer
	Name string

	// Description of what this importer does
	Description string

	// Configuration map for importer-specific settings
	Settings map[string]interface{}

	// Schedule for automated imports (optional)
	Schedule *ImportSchedule

	// Mapping rules for data transformation
	Mappings map[string]string
}

// ImportSchedule represents the schedule for automated imports
type ImportSchedule struct {
	// Frequency of imports
	Interval time.Duration

	// Time of day to run import (optional)
	TimeOfDay *time.Time

	// Days of week to run import (optional)
	DaysOfWeek []time.Weekday

	// Whether the import is currently active
	Active bool
}

// ImportProgress represents the current progress of an import operation
type ImportProgress struct {
	// Total number of items to import
	Total int

	// Number of items processed so far
	Processed int

	// Number of items successfully imported
	Succeeded int

	// Number of items that failed to import
	Failed int

	// Current status message
	Status string

	// Any errors that occurred during import
	Errors errbuilder.ErrorMap

	// Start time of the import
	StartTime time.Time

	// End time of the import (if completed)
	EndTime *time.Time
}

// ImportResult represents the final result of an import operation
type ImportResult struct {
	// Whether the import was successful overall
	Success bool

	// Total number of items processed
	TotalProcessed int

	// Number of items successfully imported
	Succeeded int

	// Number of items that failed
	Failed int

	// Number of items skipped (e.g., duplicates)
	Skipped int

	// Any errors that occurred during import
	Errors errbuilder.ErrorMap

	// Start time of the import
	StartTime time.Time

	// End time of the import
	EndTime time.Time

	// Summary of the import operation
	Summary string
}

// ImportOptions represents options for the import operation
type ImportOptions struct {
	// Whether to detect and skip duplicates
	SkipDuplicates bool

	// Whether to apply Firefly III rules to imported transactions
	ApplyRules bool

	// Whether to perform a dry run (no actual import)
	DryRun bool

	// Default category for uncategorized items
	DefaultCategory string

	// Default account for unmatched accounts
	DefaultAccount string

	// Custom field mappings
	FieldMappings map[string]string
}

// Importer defines the interface for data importers
type Importer interface {
	// Initialize sets up the importer with the given configuration
	Initialize(ctx context.Context, config ImporterConfig) error

	// ValidateConfig checks if the configuration is valid
	ValidateConfig(config ImporterConfig) error

	// TestConnection tests the connection to the data source
	TestConnection(ctx context.Context) error

	// Import performs the import operation with the given options
	Import(ctx context.Context, options ImportOptions) (*ImportResult, error)

	// GetProgress returns the current progress of an ongoing import
	GetProgress(ctx context.Context) (*ImportProgress, error)

	// Cancel stops the current import operation
	Cancel(ctx context.Context) error

	// Cleanup performs any necessary cleanup after import
	Cleanup(ctx context.Context) error

	// GetCapabilities returns the supported features of this importer
	GetCapabilities() ImporterCapabilities
}

// ImporterCapabilities represents the features supported by an importer
type ImporterCapabilities struct {
	// Supported data types for import
	SupportedTypes []string

	// Whether the importer supports scheduled imports
	SupportsScheduling bool

	// Whether the importer supports progress tracking
	SupportsProgress bool

	// Whether the importer supports cancellation
	SupportsCancellation bool

	// Whether the importer supports dry run
	SupportsDryRun bool

	// Whether the importer supports duplicate detection
	SupportsDuplicateDetection bool

	// Whether the importer supports custom field mappings
	SupportsFieldMappings bool

	// Maximum number of items that can be imported at once
	MaxBatchSize *int

	// Supported authentication methods
	AuthMethods []string
}

// BaseImporter provides a basic implementation of the Importer interface
type BaseImporter struct {
	config     ImporterConfig
	progress   *ImportProgress
	cancelled  bool
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewBaseImporter creates a new BaseImporter instance
func NewBaseImporter() *BaseImporter {
	return &BaseImporter{
		progress: &ImportProgress{},
	}
}

// Initialize implements basic initialization for importers
func (b *BaseImporter) Initialize(ctx context.Context, config ImporterConfig) error {
	b.config = config
	b.ctx, b.cancelFunc = context.WithCancel(ctx)
	b.progress = &ImportProgress{
		StartTime: time.Now(),
	}
	return nil
}

// GetProgress returns the current progress
func (b *BaseImporter) GetProgress(ctx context.Context) (*ImportProgress, error) {
	return b.progress, nil
}

// Cancel stops the current import operation
func (b *BaseImporter) Cancel(ctx context.Context) error {
	if b.cancelFunc != nil {
		b.cancelled = true
		b.cancelFunc()
	}
	return nil
}

// Cleanup performs basic cleanup
func (b *BaseImporter) Cleanup(ctx context.Context) error {
	b.progress = nil
	b.cancelled = false
	return nil
}

// UpdateProgress updates the progress information
func (b *BaseImporter) UpdateProgress(processed, succeeded, failed int, status string) {
	if b.progress != nil {
		b.progress.Processed = processed
		b.progress.Succeeded = succeeded
		b.progress.Failed = failed
		b.progress.Status = status
	}
}

// IsCancelled returns whether the import has been cancelled
func (b *BaseImporter) IsCancelled() bool {
	return b.cancelled
}
