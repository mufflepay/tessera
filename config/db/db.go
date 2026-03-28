package db

import (
	"backend/config"
	"backend/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	// postgres driver
	"gorm.io/driver/postgres"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

type IDatabase interface {
	GetDB() (*gorm.DB, error)
}

// Database represents a database connection.
type Database struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB, rdb *redis.Client) IDatabase {
	return &Database{db: db}
}

// ConnectDB creates a connection to a Postgres and Redis database.
func ConnectDB(config *config.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to Postgres database
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err.Error())
	}
	// db.Logger = logger.Default.LogMode(logger.Info)

	// GORM using database/sql to maintain connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	//db.Exec("CREATE OR REPLACE FUNCTION copy_ticket_status_to_history() RETURNS trigger LANGUAGE plpgsql AS $function$ BEGIN IF TG_OP = 'UPDATE' THEN INSERT INTO ticket_status_histories (ticket_id, status_id, status, changed_by, created_on, modified_on) VALUES (OLD.ticket_id, OLD.status_id, OLD.status, OLD.changed_by, OLD.created_on, NOW()); END IF; RETURN OLD; END; $function$;")
	//db.Exec("CREATE OR REPLACE FUNCTION copy_ticket_priority_to_history() RETURNS trigger LANGUAGE plpgsql AS $function$ BEGIN IF TG_OP = 'UPDATE' THEN INSERT INTO ticket_priority_histories (ticket_id, priority_id, priority, changed_by, created_on, modified_on) VALUES (OLD.ticket_id, OLD.priority_id, OLD.priority, OLD.changed_by, OLD.created_on, NOW()); END IF; RETURN OLD; END; $function$;")

	log.Println("Running Migrations")
	err = db.AutoMigrate(&models.Ticket{}, &models.TicketNote{}, &models.TicketAttachment{}, &models.TicketStatus{}, &models.TicketPriority{})
	if err != nil {
		log.Fatal("Migration Failed:  \n", err.Error())
	}

	//db.Exec("CREATE OR REPLACE TRIGGER ticket_status_update_trigger AFTER UPDATE ON ticket_statuses FOR EACH ROW EXECUTE FUNCTION public.copy_ticket_status_to_history();")
	//db.Exec("CREATE OR REPLACE TRIGGER ticket_priority_update_trigger AFTER UPDATE ON ticket_priorities FOR EACH ROW EXECUTE FUNCTION public.copy_ticket_priority_to_history();")

	log.Println("🚀 Connected Successfully to the Database")

	// // Connect to Redis database
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     config.RedisHost + ":" + config.RedisPort,
	// 	Password: "",
	// 	DB:       0, // use default DB(db0)
	// })

	// if err := rdb.Ping().Err(); err != nil {
	// 	return nil, err
	// }

	// err = rdb.Set("test", "connected", 60*time.Second).Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("🚀 Connected Successfully to Redis")

	return &Database{db: db}, nil
}

// GetDB returns the underlying sql.DB instance.
func (d *Database) GetDB() (*gorm.DB, error) {
	if d.db == nil {
		return nil, fmt.Errorf("could not get database")
	}
	return d.db, nil
}

// GetRedis returns the underlying Redis client instance.
// func (d *Database) GetRedis() (*redis.Client, error) {
// 	if d.rdb == nil {
// 		return nil, fmt.Errorf("could not get Redis client")
// 	}
// 	return d.rdb, nil
// }

func SetupDatabase() (*gorm.DB, error) {
	// Load environment variables
	loadConfig, err := config.LoadConfig("./")
	if err != nil {
		return nil, err
	}

	// Check if running in development environment
	if env := loadConfig.GoDev; env != "development" {
		return nil, errors.New("this server is only intended for development use")
	}

	// Connect to the database
	db, err := ConnectDB(&loadConfig)
	if err != nil {
		return nil, err
	}

	dbGet, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	// dbRedis, err := db.GetRedis()
	// if err != nil {
	// 	return nil, nil, err
	// }

	return dbGet, nil

	// return &Database{db: dbGet, rdb: dbRedis}, nil
}
