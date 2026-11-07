//go:build sqlite

package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"

	_ "modernc.org/sqlite"
)

// init registers the SQLite provider constructor with promptsdb so that
// promptsdb.New(promptsdb.SQLITE_PROVIDER, ...) can build it without an
// import cycle (promptsdb cannot import this package back).
func init() {
	promptsdb.RegisterProvider(promptsdb.SQLITE_PROVIDER, func(config promptsdb.ProviderConfiguration, logfile string) (promptsdb.Provider, error) {
		return NewSQLiteProvider(SQLiteProviderConfiguration{
			Path:         config.SQLite.Path,
			Timeout:      config.SQLite.Timeout,
			ForeignKeys:  config.SQLite.ForeignKeys,
			MaxOpenConns: config.SQLite.MaxOpenConns,
			MaxIdleConns: config.SQLite.MaxIdleConns,
		}, logfile)
	})
}

// SQLiteProvider implements the promptsdb.Provider interface using SQLite as the backend.
type SQLiteProvider struct {
	mu     sync.RWMutex
	db     *sql.DB
	dbPath string
	log    *plog.Plogger
}

// SQLiteProviderConfiguration holds SQLite-specific configuration options.
type SQLiteProviderConfiguration struct {
	Path         string `yaml:"path"`
	Timeout      int    `yaml:"timeout"`        // connection timeout in seconds
	ForeignKeys  bool   `yaml:"foreign_keys"`   // enable foreign key constraints
	MaxOpenConns int    `yaml:"max_open_conns"` // max open connections
	MaxIdleConns int    `yaml:"max_idle_conns"` // max idle connections
}

// NewSQLiteProvider creates a new SQLiteProvider instance, initializes the database
// schema, and returns the provider ready for use.
func NewSQLiteProvider(config SQLiteProviderConfiguration, logfile string) (promptsdb.Provider, error) {
	p := plog.New(logfile)
	p.Write(plog.SERVER, "setting up new SQLite provider")

	if config.Path == "" {
		return nil, errors.New("sqlite provider: database path is required")
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	maxOpen := config.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 10
	}

	maxIdle := config.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 5
	}

	// Build DSN with pragmas
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(%v)&_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)",
		config.Path, config.ForeignKeys, timeout*1000)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite provider: failed to open database: %w", err)
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Duration(timeout) * time.Second)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite provider: failed to connect to database: %w", err)
	}

	provider := &SQLiteProvider{
		db:     db,
		dbPath: config.Path,
		log:    p,
	}

	// Initialize schema
	if err := provider.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite provider: schema migration failed: %w", err)
	}

	p.Write(plog.SERVER, "SQLite provider initialized successfully at "+config.Path)

	return provider, nil
}

// Create stores a new prompt in the database.
func (s *SQLiteProvider) Create(prompt promptsdb.Prompt) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validatePrompt(prompt); err != nil {
		return err
	}

	if prompt.Id == "" {
		prompt.Id = prompt.Name
	}

	// Check for duplicate name
	var exists int
	err := s.db.QueryRow("SELECT COUNT(*) FROM prompts WHERE name = ?", prompt.Name).Scan(&exists)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to check existing prompt: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("sqlite provider: prompt with name %q already exists", prompt.Name)
	}

	tagsJSON, err := json.Marshal(prompt.Tags)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to marshal tags: %w", err)
	}

	argumentsJSON, err := json.Marshal(prompt.Arguments)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to marshal arguments: %w", err)
	}

	now := time.Now().UTC()

	_, err = s.db.Exec(
		`INSERT INTO prompts (id, name, title, description, arguments, content, tags, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		prompt.Id, prompt.Name, prompt.Title, prompt.Description,
		string(argumentsJSON), prompt.Content, string(tagsJSON),
		now, now,
	)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to create prompt: %w", err)
	}

	return nil
}

// Read retrieves a prompt by its ID (name) from the database.
func (s *SQLiteProvider) Read(promptId string) (promptsdb.Prompt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow(
		`SELECT id, name, title, description, arguments, content, tags, created_at, updated_at
		 FROM prompts WHERE name = ?`, promptId,
	)

	var prompt promptsdb.Prompt
	var argumentsJSON, tagsJSON string
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&prompt.Id, &prompt.Name, &prompt.Title, &prompt.Description,
		&argumentsJSON, &prompt.Content, &tagsJSON,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return promptsdb.Prompt{}, fmt.Errorf("sqlite provider: no prompt with id %q found", promptId)
		}
		return promptsdb.Prompt{}, fmt.Errorf("sqlite provider: failed to read prompt: %w", err)
	}

	if err := json.Unmarshal([]byte(argumentsJSON), &prompt.Arguments); err != nil {
		return promptsdb.Prompt{}, fmt.Errorf("sqlite provider: failed to unmarshal arguments: %w", err)
	}
	if err := json.Unmarshal([]byte(tagsJSON), &prompt.Tags); err != nil {
		return promptsdb.Prompt{}, fmt.Errorf("sqlite provider: failed to unmarshal tags: %w", err)
	}

	return prompt, nil
}

// Update modifies an existing prompt in the database.
func (s *SQLiteProvider) Update(prompt promptsdb.Prompt) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validatePrompt(prompt); err != nil {
		return err
	}

	if prompt.Id == "" {
		prompt.Id = prompt.Name
	}

	tagsJSON, err := json.Marshal(prompt.Tags)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to marshal tags: %w", err)
	}

	argumentsJSON, err := json.Marshal(prompt.Arguments)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to marshal arguments: %w", err)
	}

	now := time.Now().UTC()

	result, err := s.db.Exec(
		`UPDATE prompts SET title = ?, description = ?, arguments = ?, content = ?, tags = ?, updated_at = ?
		 WHERE name = ?`,
		prompt.Title, prompt.Description, string(argumentsJSON), prompt.Content,
		string(tagsJSON), now, prompt.Name,
	)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to update prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("sqlite provider: no prompt with name %q found to update", prompt.Name)
	}

	return nil
}

// Delete removes a prompt from the database by its ID (name).
func (s *SQLiteProvider) Delete(promptId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec("DELETE FROM prompts WHERE name = ?", promptId)
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to delete prompt: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite provider: failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("sqlite provider: no prompt with id %q found to delete", promptId)
	}

	return nil
}

// List returns all prompts, optionally filtered by the query parameters.
func (s *SQLiteProvider) List(query promptsdb.PromptQuery) ([]promptsdb.Prompt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rows *sql.Rows
	var err error

	switch {
	case query.NameStartsWith != "":
		rows, err = s.db.Query(
			`SELECT id, name, title, description, arguments, content, tags, created_at, updated_at
			 FROM prompts WHERE name LIKE ? ORDER BY name ASC`,
			query.NameStartsWith+"%",
		)
	case query.NameContains != "":
		rows, err = s.db.Query(
			`SELECT id, name, title, description, arguments, content, tags, created_at, updated_at
			 FROM prompts WHERE name LIKE ? ORDER BY name ASC`,
			"%"+query.NameContains+"%",
		)
	default:
		rows, err = s.db.Query(
			`SELECT id, name, title, description, arguments, content, tags, created_at, updated_at
			 FROM prompts ORDER BY name ASC`,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite provider: failed to list prompts: %w", err)
	}
	defer rows.Close()

	var prompts []promptsdb.Prompt
	for rows.Next() {
		var prompt promptsdb.Prompt
		var argumentsJSON, tagsJSON string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&prompt.Id, &prompt.Name, &prompt.Title, &prompt.Description,
			&argumentsJSON, &prompt.Content, &tagsJSON,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("sqlite provider: failed to scan prompt row: %w", err)
		}

		if err := json.Unmarshal([]byte(argumentsJSON), &prompt.Arguments); err != nil {
			return nil, fmt.Errorf("sqlite provider: failed to unmarshal arguments: %w", err)
		}
		if err := json.Unmarshal([]byte(tagsJSON), &prompt.Tags); err != nil {
			return nil, fmt.Errorf("sqlite provider: failed to unmarshal tags: %w", err)
		}

		prompts = append(prompts, prompt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite provider: error iterating rows: %w", err)
	}

	// Apply index range if specified
	if query.IndexFrom > 0 || query.IndexTo > 0 {
		from := query.IndexFrom
		to := len(prompts)
		if query.IndexTo > 0 && query.IndexTo < to {
			to = query.IndexTo
		}
		if from < len(prompts) {
			prompts = prompts[from:to]
		} else {
			prompts = []promptsdb.Prompt{}
		}
	}

	return prompts, nil
}

// Close closes the database connection.
func (s *SQLiteProvider) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Close()
}

// BeginTx starts a new database transaction for atomic multi-prompt operations.
func (s *SQLiteProvider) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

// CreateInTx stores a prompt within an existing transaction.
func CreateInTx(tx *sql.Tx, prompt promptsdb.Prompt) error {
	if err := validatePrompt(prompt); err != nil {
		return err
	}
	if prompt.Id == "" {
		prompt.Id = prompt.Name
	}

	tagsJSON, _ := json.Marshal(prompt.Tags)
	argumentsJSON, _ := json.Marshal(prompt.Arguments)
	now := time.Now().UTC()

	_, err := tx.Exec(
		`INSERT INTO prompts (id, name, title, description, arguments, content, tags, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		prompt.Id, prompt.Name, prompt.Title, prompt.Description,
		string(argumentsJSON), prompt.Content, string(tagsJSON),
		now, now,
	)
	return err
}

// validatePrompt checks that required fields are present.
func validatePrompt(prompt promptsdb.Prompt) error {
	if strings.TrimSpace(prompt.Name) == "" {
		return errors.New("sqlite provider: prompt name is required")
	}
	if strings.TrimSpace(prompt.Content) == "" {
		return errors.New("sqlite provider: prompt content is required")
	}
	if prompt.Id == "" {
		prompt.Id = uuid.New().String()
	}
	return nil
}
