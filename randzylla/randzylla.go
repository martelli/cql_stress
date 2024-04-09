package randzylla

import (
	"fmt"
	"math/rand"

	"github.com/gocql/gocql"
)

const (
	NAME_LENGTH    = 20
	TEST_KEYSPACE  = "test_keyspace"
	TEST_TABLENAME = "test_table"
)

type randzylla struct {
	session *gocql.Session
}

func NewRandzylla(addr string) (*randzylla, error) {
	cluster := gocql.NewCluster(addr)
	var err error

	r := &randzylla{}

	r.session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	if err = r.TearDown(); err != nil {
		return nil, err
	}

	q_keyspace := fmt.Sprintf("CREATE KEYSPACE %s WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1}", TEST_KEYSPACE)
	if err = r.session.Query(q_keyspace).Exec(); err != nil {
		return nil, err
	}

	cluster.Keyspace = TEST_KEYSPACE
	r.session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	q_drop_table := fmt.Sprintf("DROP TABLE IF EXISTS %s", TEST_TABLENAME)
	if err = r.session.Query(q_drop_table).Exec(); err != nil {
		return nil, err
	}

	create_query := fmt.Sprintf("CREATE TABLE %s (id int PRIMARY KEY)", TEST_TABLENAME)
	if err = r.session.Query(create_query).Exec(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *randzylla) TearDown() error {
	q_drop := fmt.Sprintf("DROP KEYSPACE IF EXISTS %s", TEST_KEYSPACE)
	if err := r.session.Query(q_drop).Exec(); err != nil {
		return err
	}
	return nil
}

func (r *randzylla) GetInsertFunction() func() error {
	query := fmt.Sprintf("INSERT INTO %s (id) VALUES (?)", TEST_TABLENAME)

	f := func() error {
		return r.session.Query(query, rand.Int31()).Exec()
	}

	return f
}
