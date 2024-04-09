package randzylla

import (
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/gocql/gocql"
)

const (
	NAME_LENGTH    = 20
	TEST_KEYSPACE  = "test_keyspace"
	TEST_TABLENAME = "test_table"
)

type randzylla struct {
	session *gocql.Session
	serial atomic.Uint32
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

	create_query := fmt.Sprintf("CREATE TABLE %s (id INT PRIMARY KEY, va INT, vb INT, vc INT)", TEST_TABLENAME)
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

func (r *randzylla) getInsertFunction(valuesFunc func() []any) func() error {
	query := fmt.Sprintf("INSERT INTO %s (id, va, vb, vc) VALUES (?,?,?,?)", TEST_TABLENAME)

	f := func() error {
		return r.session.Query(query, valuesFunc()...).Exec()
	}

	return f
}

func (r *randzylla) GetRandomInsertFunction() func() error {
	randomValues := func() []any {
		return []any{
			rand.Int31(),
			rand.Int31(),
			rand.Int31(),
			rand.Int31(),
		}
	}
	return r.getInsertFunction(randomValues)
}

func (r *randzylla) GetSerialInsertFunction() func() error {
	serialValues := func() []any {
		return []any{
			r.serial.Add(1),
			r.serial.Add(1),
			r.serial.Add(1),
			r.serial.Add(1),
		}
	}
	return r.getInsertFunction(serialValues)
}
