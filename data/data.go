// Package data implements sql operations over the NoName e.V. website-db.
package data

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	driver  = flag.String("driver", "postgres", "Der benutzte sql-Treiber")

	// If the connection string is empty, the postgres driver will extract the
	// config out of environment variables (PGHOST, PGUSER, PGDATABASE, ...).
	connect = flag.String("connect", "", "Die Verbindungsspezifikation")
)

// OpenDB opens a connection to the database with parameters derived from flags.
func OpenDB() (*sql.DB, error) {
	return sql.Open(*driver, *connect)
}

// Querier is an interface used to query the database. Both *sql.DB and *sql.Tx
// implement it.
type Querier interface {
	// Query executes a query that returns rows, typically a SELECT. The args
	// are for any placeholder parameters in the query.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until Row's
	// Scan method is called.
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Execer is an interface used to write to the database. Both *sql.DB and
// *sql.Tx implement it.
type Execer interface {
	// Exec executes a query without returning any rows. The args are for any
	// placeholder parameters in the query.
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// NullTime represents a time.Time that may be null. NullTime implements the
// sql.scanner interface so it can be used as a scan destination, similar to
// sql.NullString.
type NullTime struct {
	Time  time.Time
	Valid bool // Vaid is true if Time is not NULL
}

// Scan implements the sql.scanner interface.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*n = NullTime{}
		return nil
	}
	n.Valid = true

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("can't save %T as time.Time", value)
	}
	n.Time = t
	return nil
}

// scanner is an interface to retrieve values from an sql row. Both *sql.Row
// and *sql.Rows implement it.
type scanner interface {
	// Scan copies the columns from the matched row into the values pointed at
	// by dest. See the documentation on Rows.Scan for details. If more than
	// one row matches the query, Scan uses the first row and discards the
	// rest. If no row matches the query, Scan returns ErrNoRows.
	Scan(dest ...interface{}) error
}

// row is an abstraction over the data elements to share code.
type item interface {
	// selectFragment returns an SQL query fragment of the form "SELECT field1,
	// field2,... " with all columns for this type.
	selectFragment() string

	// scanFrom calls the scanners Scan method to retrieve the fields in the
	// order output by selectFragment.
	scanFrom(scanner) error
}

// queryRow is a convenience function to get a single data row.
func queryRow(q Querier, r item, constrict string, args ...interface{}) error {
	row := q.QueryRow(r.selectFragment()+constrict, args...)
	return r.scanFrom(row)
}

// TerminIterator is an iterator over a subset of the termine table.
type TerminIterator struct {
	rows *sql.Rows
	err  error
	t    *Termin
}

// Next advances to the next element. It returns false if there is none.
func (it *TerminIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if !it.rows.Next() {
		it.err = it.rows.Err()
		return false
	}
	it.t = new(Termin)
	it.err = it.t.scanFrom(it.rows)
	return it.err == nil
}

// Val returns the current row. Requires Next() to be true.
func (it *TerminIterator) Val() *Termin {
	return it.t
}

// Close closes the iterator and returns the last error that occured when
// reading. It must be called after being done with the iterator. No other
// methods may be called after Close.
func (it *TerminIterator) Close() error {
	err := it.rows.Close()
	it.rows = nil
	if it.err == nil {
		it.err = err
	}
	return it.err
}

// One is a convenience method for queries that expect exactly one result. It
// returns an error if there is not exactly one result. Must be the only method
// being called.
func (it *TerminIterator) One() (*Termin, error) {
	defer it.rows.Close()
	if !it.Next() {
		if it.err == nil {
			it.err = sql.ErrNoRows
		}
		return nil, it.err
	}
	t := it.Val()
	if it.Next() {
		return nil, errors.New("more than one result")
	}
	return t, nil
}

// First is a convenience method for queries where only the first result is of
// interest. It returns an error if there is no result. Must be the only method
// being called.
func (it *TerminIterator) First() (*Termin, error) {
	defer it.rows.Close()
	if !it.Next() {
		if it.err == nil {
			it.err = sql.ErrNoRows
		}
		return nil, it.err
	}
	return it.Val(), nil
}

// QueryTermine queries the termine table for all rows where cond is true. cond
// must be a valid SQL fragment following a "SELECT field,..." statement on
// termine. If possible, one of the more specific functions should be used.
func QueryTermine(q Querier, conds string, args ...interface{}) *TerminIterator {
	sel := (*Termin)(nil).selectFragment()
	rows, err := q.Query(sel+conds, args...)
	return &TerminIterator{rows, err, nil}
}

// FutureTermine returns an iterator over all future meetings (starting with
// today's) in chronological order.
func FutureTermine(q Querier) *TerminIterator {
	return QueryTermine(q, "WHERE date >= $1 ORDER BY date ASC", time.Now())
}

// LastTermine returns an iterator over all past meetings (including today's),
// in reverse chronological order.
func LastTermine(q Querier) *TerminIterator {
	return QueryTermine(q, "WHERE date <= $1 ORDER BY date DESC", time.Now())
}

// Termin is the representation of a meeting.
type Termin struct {
	// Date (when set) is the date of the meeting.
	Date NullTime

	// Stammtisch (when set) is whether this meeting is a Stammtisch.
	Stammtisch sql.NullBool

	// Vortrag (when set) contains the id of this meetings talk.
	Vortrag sql.NullInt64

	// Location is the location of a potential Stammtisch.
	Location string

	// Override is a short string to display if a meeting isn't happening.
	Override string

	// OverrideLong is a long description of a meeting that isn't happening, to
	// be sent by E-Mail in the announcement.
	OverrideLong string
}

// GetTermin returns the meeting of the specified date.
func GetTermin(q Querier, date time.Time) (*Termin, error) {
	t := new(Termin)
	r := q.QueryRow(t.selectFragment()+"WHERE date = $1", date)
	err := t.scanFrom(r)
	return t, err
}

func (t *Termin) selectFragment() string {
	return "SELECT date, stammtisch, vortrag, location, override, override_long FROM termine "
}

func (t *Termin) scanFrom(s scanner) error {
	return s.Scan(&t.Date, &t.Stammtisch, &t.Vortrag, &t.Location, &t.Override, &t.OverrideLong)
}

// GetVortrag returns the talk of this meeting. It returns nil, nil, if there
// is no talk.
func (t *Termin) GetVortrag(q Querier) (*Vortrag, error) {
	if !t.Vortrag.Valid {
		return nil, nil
	}
	return GetVortrag(q, int(t.Vortrag.Int64))
}

// Update writes back the meeting data to the database.
func (t *Termin) Update(e Execer) error {
	var fields []string
	var values []interface{}

	add := func(k string, v interface{}) {
		values = append(values, v)
		fields = append(fields, fmt.Sprintf("%s = $%d", k, len(values)))
	}

	if t.Stammtisch.Valid {
		add("stammtisch", t.Stammtisch.Bool)
	} else {
		fields = append(fields, "stammtisch = NULL")
	}
	if t.Vortrag.Valid {
		add("vortrag", t.Vortrag.Int64)
	} else {
		fields = append(fields, "vortrag = NULL")
	}
	add("location", t.Location)
	add("override", t.Override)
	add("override_long", t.OverrideLong)

	values = append(values, t.Date.Time)
	query := fmt.Sprintf("UPDATE termine SET %s WHERE date = $%d", strings.Join(fields, ", "), len(values))

	_, err := e.Exec(query, values...)
	return err
}

// Insert adds this meeting to the database.
func (t *Termin) Insert(e Execer) error {
	if !t.Date.Valid {
		return errors.New("termin needs a date")
	}
	var fields []string
	var values []interface{}

	add := func(field string, value interface{}) {
		fields = append(fields, field)
		values = append(values, value)
	}

	add("date", t.Date.Time)
	if t.Stammtisch.Valid {
		add("stammtisch", t.Stammtisch.Bool)
	}
	if t.Vortrag.Valid {
		add("vortrag", t.Vortrag.Int64)
	}
	add("location", t.Location)
	add("override", t.Override)
	add("override_long", t.OverrideLong)

	placeholder := make([]string, 0, len(values))
	for i := 1; i <= len(values); i++ {
		placeholder = append(placeholder, "$"+strconv.Itoa(i))
	}

	query := fmt.Sprintf("INSERT INTO termine (%s) SELECT %s WHERE NOT EXISTS (SELECT 1 FROM termine WHERE date = $1)", strings.Join(fields, ", "), strings.Join(placeholder, ", "))
	_, err := e.Exec(query, values...)
	return err
}

// Vortrag is the representation of a Talk.
type Vortrag struct {
	// ID is this talks unique id.
	ID int

	// Date (if set) is the date of this talk.
	Date NullTime

	// Topic is the topic of this talk.
	Topic string

	// Abstract is a short summary of the talk.
	Abstract string

	// Speaker is the name of the speaker.
	Speaker string

	// InfoURL is a url where to find further information.
	InfoURL string

	// Password is the password to edit this talk.
	Password string
}

func (v *Vortrag) selectFragment() string {
	return "SELECT id, date, topic, abstract, speaker, infourl, password FROM vortraege "
}

func (v *Vortrag) scanFrom(s scanner) error {
	return s.Scan(&v.ID, &v.Date, &v.Topic, &v.Abstract, &v.Speaker, &v.InfoURL, &v.Password)
}

// GetVortrag returns the talk with the given id.
func GetVortrag(q Querier, id int) (*Vortrag, error) {
	v := new(Vortrag)
	r := q.QueryRow(v.selectFragment()+"WHERE id = $1", id)
	err := v.scanFrom(r)
	return v, err
}

// Link is the representation of an informational link.
type Link struct {
	// Kind (if given) is the type of this link.
	Kind sql.NullString

	// URL is the url of this link.
	URL *url.URL
}

// Links returns the informational links of this talk.
func (v *Vortrag) Links(q Querier) ([]Link, error) {
	return nil, nil
}

// Zusage is the representation of an RSVP.
type Zusage struct {
	// Nick is the nick this applies to.
	Nick sql.NullString

	// Kommt is whether this person intends to be there.
	Kommt bool

	// Kommentar is the optional comment given by this person.
	Kommentar string

	// Registered (if given) is the time this person RSVPed.
	Registered NullTime
}

func (z *Zusage) selectFragment() string {
	return "SELECT nick, kommt, kommentar, registered FROM zusagen "
}

func (z *Zusage) scanFrom(s scanner) error {
	return s.Scan(&z.Nick, &z.Kommt, &z.Kommentar, &z.Registered)
}

// ZusagenIterator is an iterator over a subset of the zusagen table.
type ZusagenIterator struct {
	rows *sql.Rows
	err  error
	z    *Zusage
}

// Next advances to the next element. It returns false if there is none.
func (it *ZusagenIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if !it.rows.Next() {
		it.err = it.rows.Err()
		return false
	}
	it.z = new(Zusage)
	it.err = it.z.scanFrom(it.rows)
	return it.err == nil
}

// Val returns the current row. Requires Next() to be true.
func (it *ZusagenIterator) Val() *Zusage {
	return it.z
}

// Close closes the iterator and returns the last error that occured when
// reading. It must be called after being done with the iterator. No other
// methods may be called after Close.
func (it *ZusagenIterator) Close() error {
	if it.rows != nil {
		err := it.rows.Close()
		it.rows = nil
		if it.err == nil {
			it.err = err
		}
	}
	return it.err
}

// Zusagen returns all rows of the zusagen table in unspecified order.
func Zusagen(q Querier) *ZusagenIterator {
	return QueryZusagen(q, "")
}

// QueryZusagen queries the zusagen table for all rows where cond is true. cond
// must be a valid SQL fragment following a "SELECT field,..." statement on
// zusagen. If possible, one of the more specific functions should be used.
func QueryZusagen(q Querier, conds string, args ...interface{}) *ZusagenIterator {
	sel := (*Zusage)(nil).selectFragment()
	rows, err := q.Query(sel+conds, args...)
	return &ZusagenIterator{rows, err, nil}
}
