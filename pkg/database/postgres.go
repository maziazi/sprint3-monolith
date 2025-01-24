package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"sprint3/pkg/config"
	"sync"
	"time"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func InitDB() {
	once.Do(func() {
		// Memuat konfigurasi dari .env
		cfg := config.LoadEnv()

		// Buat connection string PostgreSQL dengan konfigurasi optimal
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

		// Konfigurasi pool dengan opsi tambahan (timeout, max connections, dll.)
		poolConfig, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			log.Fatalf("Error parsing database config: %v", err)
		}

		// Atur parameter koneksi (sesuaikan dengan kebutuhan aplikasi)
		poolConfig.MaxConns = 10                       // Maksimal 10 koneksi
		poolConfig.MinConns = 2                        // Minimal 2 koneksi
		poolConfig.MaxConnLifetime = 30 * time.Minute  // Maksimal umur koneksi 30 menit
		poolConfig.MaxConnIdleTime = 5 * time.Minute   // Koneksi idle selama 5 menit akan ditutup
		poolConfig.HealthCheckPeriod = 1 * time.Minute // Cek kesehatan koneksi tiap 1 menit

		// Buat connection pool
		dbPool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		log.Println("✅ Connected to database successfully")
	})
}

func CloseDB() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("🛑 Database connection closed")
	}
}

func GetDBPool() *pgxpool.Pool {
	if dbPool == nil {
		log.Println("⚠️ Database connection is not initialized, calling InitDB()")
		InitDB()
	}
	return dbPool
}

func GetDB() *sql.DB {
	if dbPool == nil {
		log.Println("⚠️ Database connection is not initialized, calling InitDB()")
		InitDB()
	}

	// Convert pgxpool.Pool to *sql.DB
	db := stdlib.OpenDB(*dbPool.Config().ConnConfig)

	// Verify the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ SQL database instance is ready")
	return db
}
