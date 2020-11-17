package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/lib/pq"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PGStore struct {
	db *sqlx.DB
}

type PGConnString string

func NewPGStore(cs PGConnString) (PGStore, error) {
	db, err := sqlx.Connect("postgres", string(cs))
	if err != nil {
		return PGStore{}, fmt.Errorf("NewPGStore#Connect: %w", err)
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return PGStore{db: db}, nil
}

func (pgs PGStore) GetBuff(ID uint64) (*Buff, error) {
	var answers pq.StringArray
	var q string
	rows := pgs.db.QueryRow(`select question, answers from buff where id=$1`, ID)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PGStore#GetBuff#sqlx.Get: %w", err)
	}
	if err := rows.Scan(&q, &answers); err != nil {
		return nil, fmt.Errorf("PGStore#Scan: %w", err)
	}
	return &Buff{
		ID:       ID,
		Question: q,
		Answers:  answers,
	}, nil
}

// I assume that SetBuff method will be used only for creating buffs, not updating it
func (pgs PGStore) SetBuff(b *Buff) (uint64, error) {
	var ID uint64
	if err := pgs.db.QueryRowx(`insert into buff (question, answers) values($1, $2) returning id`, b.Question, pq.Array(b.Answers)).Scan(&ID); err != nil {
		return 0, fmt.Errorf("PGStore#SetBuff#QueryRowx: %w", err)
	}
	return ID, nil
}

func (pgs PGStore) GetStream(ID uint64) (Stream, error) {
	var s Stream
	//Since all code uses uint64 i should follow same approach. But it is imposible to scan []uint64 from PG by pq lib
	//possible approach is to create UInt64Array and implement scanner for it. But for this particular task i'll just
	//use pq.Array of int64 and then copy to []uint64
	var buffIDs []int64
	if err := pgs.db.QueryRowx(`
		select s.id, s.name, array_agg(bts.buff_id)
		from stream s
		left join buff_to_stream bts on s.id = bts.stream_id
		where id = $1
		group by s.id`, ID).Scan(&s.ID, &s.Name, pq.Array(&buffIDs)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Stream{}, fmt.Errorf("GetStream#QueryRowx: %w", ErrNotFound)
		}
		return Stream{}, fmt.Errorf("PGStore#GetStream#Get: %w", err)
	}
	s.BuffIDs = make([]uint64, len(buffIDs))
	//it is also possible to do some fancy runtime efficient casting with unsafe and blackjack
	//but it is premature optimisation so skip it
	for i := range buffIDs {
		s.BuffIDs[i] = uint64(buffIDs[i])
	}
	return s, nil
}

func (pgs PGStore) SetStream(s Stream) (uint64, error) {
	var ID uint64
	t, err := pgs.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("PGStore#SetStream#Begin: %w", err)
	}
	if err := t.QueryRow(`insert into stream (name) values ($1) returning id`, s.Name).Scan(&ID); err != nil {
		return 0, fmt.Errorf("PGStore#SetStream#QueryRow: %w", err)
	}

	stmt, err := t.Preparex(pq.CopyIn("buff_to_stream", "stream_id", "buff_id"))
	if err != nil {
		return 0, fmt.Errorf("PGStore#SetStream#Prepare: %w", err)
	}
	for _, bID := range s.BuffIDs {
		_, err = stmt.Exec(ID, bID)
		if err != nil {
			return 0, fmt.Errorf("PGStore#SetStream#Exec: %w", err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return 0, fmt.Errorf("PGStore#SetStream#ExecFlush: %w", err)
	}
	defer stmt.Close()
	if err := t.Commit(); err != nil {
		return 0, fmt.Errorf("PGStore#AttachBuffs#Commit: %w", err)
	}

	return ID, nil
}
