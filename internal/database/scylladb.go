package database

import (
	"log"
	"os"

	"github.com/gocql/gocql"
)

type ScyllaService interface {
	Health() map[string]string
}

type scyllaService struct {
	session *gocql.Session
}

func NewScylla() ScyllaService {
	cluster := gocql.NewCluster(os.Getenv("DB_HOST"))
	cluster.Port = 9042
	cluster.Keyspace = os.Getenv("DB_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Cannot connect to ScyllaDB:", err)
	}

	s := &scyllaService{
		session: session,
	}
	return s
}

func (s *scyllaService) Health() map[string]string {
	if err := s.session.Query(`SELECT now() FROM system.local`).Exec(); err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
