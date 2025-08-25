package store

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"learnlang-backend/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB
var (
	hasVocabTranslation bool
	probedSchema        bool
)

// Init sets the global DB connection for the store package.
func Init(d *sql.DB) {
	db = d
}

// Close closes the global DB connection (optional helper).
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// InitFromEnv opens a DB connection from env vars and sets it for the store.
// It supports DATABASE_URL directly, or POSTGRES_* vars with sensible defaults.
func InitFromEnv() error {
	if db != nil {
		return nil
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		user := getenvDefault("POSTGRES_USER", "learnlang")
		pass := getenvDefault("POSTGRES_PASSWORD", "learnlang")
		host := getenvDefault("POSTGRES_HOST", "localhost")
		port := getenvDefault("POSTGRES_PORT", "5432")
		name := getenvDefault("POSTGRES_DB", "learnlang")
		// Auto-switch to test DB when under `go test`, unless explicitly overridden
		if flag.Lookup("test.v") != nil {
			name = getenvDefault("POSTGRES_TEST_DB", "learnlang_test")
		}
		// Explicit envs still take precedence
		if os.Getenv("TEST_MODE") == "1" {
			name = getenvDefault("POSTGRES_TEST_DB", "learnlang_test")
		} else if v := os.Getenv("POSTGRES_TEST_DB"); v != "" {
			name = v
		}
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
	}
	// Use pgx stdlib driver
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxIdleTime(5 * time.Minute)
	conn.SetConnMaxLifetime(60 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return err
	}
	Init(conn)
	// Probe schema once
	_ = probeSchema(context.Background())
	return nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// probeSchema detects optional columns to keep compatibility with older DBs.
func probeSchema(ctx context.Context) error {
	if db == nil || probedSchema {
		return nil
	}
	cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var one int
	// information_schema is portable enough for our use
	err := db.QueryRowContext(cctx, `SELECT 1 FROM information_schema.columns WHERE table_name='vocabs' AND column_name='translation'`).Scan(&one)
	hasVocabTranslation = err == nil
	probedSchema = true
	return nil
}

// LanguagesList returns all supported languages.
func LanguagesList() []models.Language {
	if db == nil {
		return []models.Language{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, `SELECT id, name, code FROM languages ORDER BY name`)
	if err != nil {
		return []models.Language{}
	}
	defer rows.Close()
	var out []models.Language
	for rows.Next() {
		var l models.Language
		if err := rows.Scan(&l.ID, &l.Name, &l.Code); err == nil {
			out = append(out, l)
		}
	}
	return out
}

// LanguageExists checks if a language code is supported.
func LanguageExists(code string) bool {
	if db == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var one int
	err := db.QueryRowContext(ctx, `SELECT 1 FROM languages WHERE code=$1`, code).Scan(&one)
	return err == nil
}

// GetLanguageByCode returns the language for the given code, if present.
func GetLanguageByCode(code string) (models.Language, bool) {
	if db == nil {
		return models.Language{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var l models.Language
	err := db.QueryRowContext(ctx, `SELECT id, name, code FROM languages WHERE code=$1`, code).Scan(&l.ID, &l.Name, &l.Code)
	if err != nil {
		return models.Language{}, false
	}
	return l, true
}

// GetAllPacks returns all packs.
func GetAllPacks() []models.Pack {
	if db == nil {
		return []models.Pack{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, `SELECT id, name, lang_id, user_id, public FROM packs ORDER BY name`)
	if err != nil {
		return []models.Pack{}
	}
	defer rows.Close()
	var out []models.Pack
	for rows.Next() {
		var p models.Pack
		if err := rows.Scan(&p.ID, &p.Name, &p.LangID, &p.UserID, &p.Public); err == nil {
			out = append(out, p)
		}
	}
	return out
}

// parsePackKey parses a key of the form userID:langID:packName (all lowercased by caller).
func parsePackKey(key string) (userID, langID, name string, err error) {
	parts := strings.Split(key, ":")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid pack key")
	}
	return parts[0], parts[1], parts[2], nil
}

// PackExistsByKey reports whether the composite pack key already exists.
func PackExistsByKey(key string) bool {
	if db == nil || key == "" {
		return false
	}
	userID, langID, name, err := parsePackKey(key)
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var one int
	err = db.QueryRowContext(ctx, `SELECT 1 FROM packs WHERE lower(user_id)=lower($1) AND lower(lang_id)=lower($2) AND lower(name)=lower($3)`, userID, langID, name).Scan(&one)
	return err == nil
}

// CreatePack stores the pack.
func CreatePack(p models.Pack, _ string) {
	if db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = db.ExecContext(ctx, `INSERT INTO packs (id, name, lang_id, user_id, public) VALUES ($1, $2, $3, $4, $5)`, p.ID, p.Name, p.LangID, p.UserID, p.Public)
}

// GetPackIDByKey returns the pack ID for the composite key if exists, else empty string.
func GetPackIDByKey(key string) string {
	if db == nil || key == "" {
		return ""
	}
	userID, langID, name, err := parsePackKey(key)
	if err != nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id string
	err = db.QueryRowContext(ctx, `SELECT id FROM packs WHERE lower(user_id)=lower($1) AND lower(lang_id)=lower($2) AND lower(name)=lower($3)`, userID, langID, name).Scan(&id)
	if err != nil {
		return ""
	}
	return id
}

// parseVocabKeyByPack parses a key packID:name (lowercased name by caller).
func parseVocabKeyByPack(key string) (packID, name string, err error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid vocab key")
	}
	return parts[0], parts[1], nil
}

// VocabExistsByKey reports whether the composite vocab key already exists.
func VocabExistsByKey(key string) bool {
	if db == nil || key == "" {
		return false
	}
	packID, name, err := parseVocabKeyByPack(key)
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var one int
	err = db.QueryRowContext(ctx, `SELECT 1 FROM vocabs WHERE pack_id=$1 AND lower(name)=lower($2)`, packID, name).Scan(&one)
	return err == nil
}

// CreateVocab stores a vocab.
func CreateVocab(v models.Vocab, _ string) {
	if db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = probeSchema(ctx)
	if hasVocabTranslation {
		_, _ = db.ExecContext(ctx, `INSERT INTO vocabs (id, image, name, translation, pack_id) VALUES ($1, $2, $3, $4, $5)`, v.ID, v.Image, v.Name, v.Translation, v.PackID)
		return
	}
	// fallback for older schema without translation column
	_, _ = db.ExecContext(ctx, `INSERT INTO vocabs (id, image, name, pack_id) VALUES ($1, $2, $3, $4)`, v.ID, v.Image, v.Name, v.PackID)
}

// ListVocabs filters by user, lang and optional pack IDs.
func ListVocabs(userID, langID string, packIDs []string) []models.Vocab {
	if db == nil {
		return []models.Vocab{}
	}
	_ = probeSchema(context.Background())
	base := `SELECT v.id, v.image, v.name, ` + func() string {
		if hasVocabTranslation {
			return "COALESCE(v.translation, '') as translation, "
		}
		return "'' as translation, "
	}() + `v.pack_id
             FROM vocabs v
             JOIN packs p ON p.id = v.pack_id
             WHERE p.user_id = $1 AND p.lang_id = $2`
	args := []any{userID, langID}
	if len(packIDs) > 0 {
		// Build IN clause safely
		placeholders := make([]string, len(packIDs))
		for i := range packIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+3)
			args = append(args, packIDs[i])
		}
		base += " AND v.pack_id IN (" + strings.Join(placeholders, ",") + ")"
	}
	base += " ORDER BY v.name"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, base, args...)
	if err != nil {
		return []models.Vocab{}
	}
	defer rows.Close()
	var out []models.Vocab
	for rows.Next() {
		var v models.Vocab
		if err := rows.Scan(&v.ID, &v.Image, &v.Name, &v.Translation, &v.PackID); err == nil {
			out = append(out, v)
		}
	}
	return out
}

// ListVocabsByPackID returns all vocabs for a single pack.
func ListVocabsByPackID(packID string) []models.Vocab {
	if db == nil || packID == "" {
		return []models.Vocab{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = probeSchema(ctx)
	query := `SELECT id, image, name, `
	if hasVocabTranslation {
		query += `COALESCE(translation, '') as translation`
	} else {
		query += `'' as translation`
	}
	query += `, pack_id FROM vocabs WHERE pack_id=$1 ORDER BY name`
	rows, err := db.QueryContext(ctx, query, packID)
	if err != nil {
		return []models.Vocab{}
	}
	defer rows.Close()
	var out []models.Vocab
	for rows.Next() {
		var v models.Vocab
		if err := rows.Scan(&v.ID, &v.Image, &v.Name, &v.Translation, &v.PackID); err == nil {
			out = append(out, v)
		}
	}
	return out
}

// Reset clears data tables (packs, vocabs). Useful for tests.
func Reset() {
	if db == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Keep languages as-is (seeded by migrations)
	_, _ = db.ExecContext(ctx, `TRUNCATE TABLE vocabs, packs RESTART IDENTITY CASCADE`)
}

// GetPackByID returns a pack by ID if present.
func GetPackByID(id string) (models.Pack, bool) {
	if db == nil {
		return models.Pack{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var p models.Pack
	err := db.QueryRowContext(ctx, `SELECT id, name, lang_id, user_id, public FROM packs WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.LangID, &p.UserID, &p.Public)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Pack{}, false
		}
		return models.Pack{}, false
	}
	return p, true
}
