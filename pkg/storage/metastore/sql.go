// Copyright 2021 The Parca Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metastore

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/google/pprof/profile"
)

var _ ProfileMetaStore = &sqlMetaStore{}

type sqlMetaStore struct {
	db *sql.DB
}

func (s *sqlMetaStore) migrate() error {
	tables := []string{
		"PRAGMA foreign_keys = ON",
		`CREATE TABLE "mappings" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"mapping_id" 		INT64,
			"start"           	INT64,
			"limit"          	INT64,
			"offset"          	INT64,
			"file"           	TEXT,
			"build_id"         	TEXT,
			"has_functions"    	BOOLEAN,
			"has_filenames"    	BOOLEAN,
			"has_line_numbers"  BOOLEAN,
			"has_inline_frames" BOOLEAN,
			"size"				INT64,
			"build_id_or_file"	TEXT
		);`,
		`CREATE INDEX idx_mapping_id ON mappings (mapping_id);`,
		`CREATE INDEX idx_mapping_key ON mappings (size, offset, build_id_or_file);`,
		`CREATE TABLE "functions" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"function_id"	INT64,
			"name"       	TEXT,
			"system_name" 	TEXT,
			"filename"   	TEXT,
			"start_line"  	INT64
		);`,
		`CREATE INDEX idx_function_id ON functions (function_id);`,
		`CREATE INDEX idx_function_key ON functions (start_line, name, system_name, filename);`,
		`CREATE TABLE "lines" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"function_id"	INTEGER NOT NULL,
			"line" 		  	INT64,
			FOREIGN KEY (function_id) REFERENCES functions (id)
		);`,
		`CREATE TABLE "locations" (
			"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"location_id"			INT64,
			"mapping_id"  			INTEGER,
			"address"  				INT64,
			"is_folded" 			BOOLEAN,
			"normalized_address"	INT64,
			"lines"					TEXT,
			FOREIGN KEY (mapping_id) REFERENCES mappings (id)
		);`,
		`CREATE INDEX idx_location_id ON locations (location_id);`,
		`CREATE INDEX idx_location_key ON locations (normalized_address, mapping_id, is_folded, lines);`,
		`CREATE TABLE "location_lines" (
			"id" INTEGER PRIMARY KEY AUTOINCREMENT,
			"location_id" 	INTEGER NOT NULL,
			"line_id" 		INTEGER NOT NULL,
			FOREIGN KEY(location_id) REFERENCES locations (id),
			FOREIGN KEY(line_id) REFERENCES lines (id)
		);`,
	}
	// TODO(kakkoyun): Additional table between location and mapping? - mapInfo from pprof

	for _, t := range tables {
		statement, err := s.db.Prepare(t)
		if err != nil {
			return err
		}

		if _, err := statement.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (s *sqlMetaStore) GetLocationByKey(k LocationKey) (*profile.Location, error) {
	var (
		l           profile.Location
		mappingPKey *int
		id, address int64
		err         error
	)
	if k.MappingID > 0 {
		err = s.db.QueryRow(
			`SELECT "location_id", "address", "is_folded", "mapping_id"
					FROM "locations" l
					JOIN "mappings" m ON l.mapping_id = m.id
					WHERE l.normalized_address=? AND l.is_folded=? AND l.lines=? AND m.id=? `,
			int64(k.Addr), k.IsFolded, k.Lines, int64(k.MappingID),
		).Scan(&id, &address, &l.IsFolded, &mappingPKey)
	} else {
		err = s.db.QueryRow(
			`SELECT "location_id", "address", "is_folded"
					FROM "locations"
					WHERE normalized_address=? AND mapping_id IS NULL AND is_folded=? AND lines=?`,
			int64(k.Addr), k.IsFolded, k.Lines,
		).Scan(&id, &address, &l.IsFolded)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	l.ID = uint64(id)
	l.Address = uint64(address)

	if mappingPKey != nil {
		mapping, err := s.getMappingByPrimaryKey(*mappingPKey)
		if err != nil {
			return nil, err
		}
		l.Mapping = mapping
	}

	lines, err := s.getLocationLines(l.ID)
	if err != nil {
		return nil, err
	}
	l.Line = lines

	return &l, nil
}

func (s *sqlMetaStore) GetLocationByID(id uint64) (*profile.Location, error) {
	var (
		l              profile.Location
		mappingPKey    *int
		locID, address int64
	)
	err := s.db.QueryRow(
		`SELECT "location_id", "address", "is_folded", "mapping_id"
				FROM "locations"
				WHERE location_id=?`, id,
	).Scan(&locID, &address, &l.IsFolded, &mappingPKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrLocationNotFound
		}
		return nil, err
	}
	l.ID = uint64(locID)
	l.Address = uint64(address)

	if mappingPKey != nil {
		mapping, err := s.getMappingByPrimaryKey(*mappingPKey)
		if err != nil {
			return nil, err
		}
		l.Mapping = mapping
	}

	lines, err := s.getLocationLines(l.ID)
	if err != nil {
		return nil, err
	}
	l.Line = lines

	return &l, nil
}

func (s *sqlMetaStore) CreateLocation(l *profile.Location) error {
	k := MakeLocationKey(l)
	var res sql.Result
	if l.Mapping != nil {
		stmt, err := s.db.Prepare(
			`INSERT INTO "locations" (location_id, address, is_folded, mapping_id, normalized_address, lines)
					values(?,?,?,?,?,?)`,
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		var mappingID int
		err = s.db.QueryRow(`SELECT "id" FROM "mappings" WHERE mapping_id=?`, int64(l.Mapping.ID)).Scan(&mappingID)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrMappingNotFound
			}
			return err
		}

		res, err = stmt.Exec(int64(l.ID), int64(l.Address), l.IsFolded, mappingID, int64(k.Addr), k.Lines)
		if err != nil {
			return err
		}
	} else {
		stmt, err := s.db.Prepare(
			`INSERT INTO "locations" (
                         location_id, address, is_folded, normalized_address, lines
                         ) values(?,?,?,?,?)`,
		)
		if err != nil {
			return err
		}
		defer stmt.Close()

		res, err = stmt.Exec(int64(l.ID), int64(l.Address), l.IsFolded, int64(k.Addr), k.Lines)
		if err != nil {
			return err
		}
	}

	locID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if err := s.createLines(l.Line, locID); err != nil {
		return err
	}

	return nil
}

func (s *sqlMetaStore) UpdateLocation(l *profile.Location) error {
	k := MakeLocationKey(l)
	var res sql.Result
	if l.Mapping != nil {
		stmt, err := s.db.Prepare(
			`UPDATE "locations" SET address=?, is_folded=?, mapping_id=?, normalized_address=?, lines=? WHERE location_id=?`,
		)

		if err != nil {
			return err
		}
		defer stmt.Close()

		var mappingID int
		err = s.db.QueryRow(`SELECT "id" FROM "mappings" WHERE mapping_id=?`, int64(l.Mapping.ID)).Scan(&mappingID)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrMappingNotFound
			}
			return err
		}

		res, err = stmt.Exec(int64(l.Address), l.IsFolded, mappingID, int64(k.Addr), k.Lines, int64(l.ID))
		if err != nil {
			return err
		}
	} else {
		stmt, err := s.db.Prepare(
			`UPDATE "locations" SET address=?, is_folded=? WHERE location_id=?`,
		)

		if err != nil {
			return err
		}
		defer stmt.Close()

		res, err = stmt.Exec(int64(l.Address), l.IsFolded, int64(l.ID))
		if err != nil {
			return err
		}
	}

	locID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	if err := s.createLines(l.Line, locID); err != nil {
		return err
	}

	return nil
}

func (s *sqlMetaStore) GetLocations() ([]*profile.Location, error) {
	rows, err := s.db.Query(
		`SELECT l."location_id", l."address", l."is_folded", m."mapping_id",
       					m."start", m."limit", m."offset", m."file", m."build_id",
       					m."has_functions", m."has_filenames", m."has_line_numbers", m."has_inline_frames"
				FROM "locations" l
				LEFT JOIN "mappings" m ON l.mapping_id = m.id`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}
	defer rows.Close()

	locs := []*profile.Location{}
	for rows.Next() {
		l := &profile.Location{}
		var (
			mappingID       *int64
			start           *int64
			limit           *int64
			offset          *int64
			file            *string
			buildID         *string
			hasFunctions    *bool
			hasFilenames    *bool
			hasLineNumbers  *bool
			hasInlineFrames *bool
			locID           int64
			locAddress      int64
		)
		err := rows.Scan(
			&locID, &locAddress, &l.IsFolded,
			&mappingID, &start, &limit, &offset, &file, &buildID,
			&hasFunctions, &hasFilenames, &hasLineNumbers, &hasInlineFrames,
		)
		if err != nil {
			return nil, err
		}
		l.ID = uint64(locID)
		l.Address = uint64(locAddress)
		if mappingID != nil {
			l.Mapping = &profile.Mapping{
				ID:              uint64(*mappingID),
				Start:           uint64(*start),
				Limit:           uint64(*limit),
				Offset:          uint64(*offset),
				File:            *file,
				BuildID:         *buildID,
				HasFunctions:    *hasFunctions,
				HasFilenames:    *hasFilenames,
				HasLineNumbers:  *hasLineNumbers,
				HasInlineFrames: *hasInlineFrames,
			}
		}

		lines, err := s.getLocationLines(l.ID)
		if err != nil {
			return nil, err
		}
		l.Line = lines

		locs = append(locs, l)
	}
	return locs, nil
}

func (s *sqlMetaStore) GetUnsymbolizedLocations() ([]*profile.Location, error) {
	rows, err := s.db.Query(
		`SELECT l."location_id", l."address", l."is_folded", m."mapping_id",
       					m."start", m."limit", m."offset", m."file", m."build_id",
       					m."has_functions", m."has_filenames", m."has_line_numbers", m."has_inline_frames"
				FROM "locations" l
				LEFT JOIN "mappings" m ON l.mapping_id = m.id
				LEFT JOIN "location_lines" ll ON l."id" = ll."location_id"
				WHERE ll."location_id" IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locs := []*profile.Location{}
	for rows.Next() {
		l := &profile.Location{}
		var (
			mappingID       *int64
			start           *int64
			limit           *int64
			offset          *int64
			file            *string
			buildID         *string
			hasFunctions    *bool
			hasFilenames    *bool
			hasLineNumbers  *bool
			hasInlineFrames *bool
			locID           int64
			locAddress      int64
		)
		err := rows.Scan(
			&locID, &locAddress, &l.IsFolded,
			&mappingID, &start, &limit, &offset, &file, &buildID,
			&hasFunctions, &hasFilenames, &hasLineNumbers, &hasInlineFrames,
		)
		if err != nil {
			return nil, err
		}
		l.ID = uint64(locID)
		l.Address = uint64(locAddress)
		if mappingID != nil {
			l.Mapping = &profile.Mapping{
				ID:              uint64(*mappingID),
				Start:           uint64(*start),
				Limit:           uint64(*limit),
				Offset:          uint64(*offset),
				File:            *file,
				BuildID:         *buildID,
				HasFunctions:    *hasFunctions,
				HasFilenames:    *hasFilenames,
				HasLineNumbers:  *hasLineNumbers,
				HasInlineFrames: *hasInlineFrames,
			}
		}

		locs = append(locs, l)
	}
	return locs, nil
}

func (s *sqlMetaStore) GetFunctionByKey(k FunctionKey) (*profile.Function, error) {
	var (
		fn profile.Function
		id int64
	)
	err := s.db.QueryRow(
		`SELECT "function_id", "name", "system_name", "filename", "start_line"
				FROM "functions"
				WHERE start_line=? AND name=? AND system_name=? AND filename=?`,
		k.StartLine, k.Name, k.SystemName, k.FileName,
	).Scan(&id, &fn.Name, &fn.SystemName, &fn.Filename, &fn.StartLine)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFunctionNotFound
		}
		return nil, err
	}
	fn.ID = uint64(id)
	return &fn, nil
}

func (s *sqlMetaStore) CreateFunction(fn *profile.Function) error {
	_, err := s.createFunction(fn)
	return err
}

func (s *sqlMetaStore) GetFunctions() ([]*profile.Function, error) {
	rows, err := s.db.Query(`SELECT "function_id", "name", "system_name", "filename", "start_line" FROM "functions"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	funcs := []*profile.Function{}
	for rows.Next() {
		f := profile.Function{}
		var id int64
		err := rows.Scan(&id, &f.Name, &f.SystemName, &f.Filename, &f.StartLine)
		if err != nil {
			return nil, err
		}
		f.ID = uint64(id)
		funcs = append(funcs, &f)
	}

	return funcs, nil
}

func (s *sqlMetaStore) GetMappingByKey(k MappingKey) (*profile.Mapping, error) {
	var (
		m                        profile.Mapping
		id, start, limit, offset int64
	)
	err := s.db.QueryRow(
		`SELECT "mapping_id", "start", "limit", "offset", "file", "build_id",
				"has_functions", "has_filenames", "has_line_numbers", "has_inline_frames"
				FROM "mappings"
				WHERE size=? AND offset=? AND build_id_or_file=?`,
		k.Size, k.Offset, k.BuildIDOrFile,
	).Scan(
		&id, &start, &limit, &offset, &m.File, &m.BuildID,
		&m.HasFunctions, &m.HasFilenames, &m.HasLineNumbers, &m.HasInlineFrames,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMappingNotFound
		}
		return nil, err
	}
	m.ID = uint64(id)
	m.Start = uint64(start)
	m.Limit = uint64(limit)
	m.Offset = uint64(offset)
	return &m, nil
}

func (s *sqlMetaStore) CreateMapping(m *profile.Mapping) error {
	stmt, err := s.db.Prepare(
		`INSERT INTO "mappings" (
                        "mapping_id", "start", "limit", "offset", "file", "build_id",
                        "has_functions", "has_filenames", "has_line_numbers", "has_inline_frames",
                        "size", "build_id_or_file"
                        ) values(?,?,?,?,?,?,?,?,?,?,?,?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	k := MakeMappingKey(m)
	_, err = stmt.Exec(
		int64(m.ID), int64(m.Start), int64(m.Limit), int64(m.Offset), m.File, m.BuildID,
		m.HasFunctions, m.HasFilenames, m.HasLineNumbers, m.HasInlineFrames,
		int64(k.Size), k.BuildIDOrFile,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlMetaStore) Close() error {
	return s.db.Close()
}

func (s *sqlMetaStore) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (s *sqlMetaStore) getMappingByPrimaryKey(pkey int) (*profile.Mapping, error) {
	var (
		m                        profile.Mapping
		id, start, limit, offset int64
	)
	err := s.db.QueryRow(
		`SELECT "mapping_id", "start", "limit", "offset", "file", "build_id",
				"has_functions", "has_filenames", "has_line_numbers", "has_inline_frames"
				FROM "mappings" WHERE id=?`, pkey,
	).Scan(
		&id, &start, &limit, &offset, &m.File, &m.BuildID,
		&m.HasFunctions, &m.HasFilenames, &m.HasLineNumbers, &m.HasInlineFrames,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrMappingNotFound
		}
		return nil, err
	}
	m.ID = uint64(id)
	m.Start = uint64(start)
	m.Limit = uint64(limit)
	m.Offset = uint64(offset)
	return &m, nil
}

func (s *sqlMetaStore) getLocationLines(locationID uint64) ([]profile.Line, error) {
	var lines []profile.Line
	rows, err := s.db.Query(
		`SELECT ln."line", fn."function_id", fn."name", fn."system_name", fn."filename", fn."start_line"
				FROM "location_lines" ll
				LEFT JOIN "locations" loc ON ll."location_id" = loc."id"
				LEFT JOIN "lines" ln ON ll."line_id" = ln."id"
				LEFT JOIN "functions" fn ON ln."function_id" = fn."id"
				WHERE loc."location_id"=?`, int64(locationID),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ln := profile.Line{}
		fn := profile.Function{}
		var fnID int64
		err := rows.Scan(&ln.Line, &fnID, &fn.Name, &fn.SystemName, &fn.Filename, &fn.StartLine)
		if err != nil {
			return nil, err
		}
		fn.ID = uint64(fnID)
		ln.Function = &fn
		lines = append(lines, ln)
	}

	// To make tests stable.
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i].Line < lines[j].Line
	})

	return lines, nil
}

func (s *sqlMetaStore) getOrCreateFunction(f *profile.Function) (int64, error) {
	var functionID int64
	err := s.db.QueryRow(`SELECT "id" FROM "functions" WHERE function_id=?`, int64(f.ID)).Scan(&functionID)
	if err != nil {
		if err == sql.ErrNoRows {
			functionID, err = s.createFunction(f)
			if err != nil {
				return 0, err
			}
			return functionID, nil
		}
		return 0, err
	}
	return functionID, nil
}

func (s *sqlMetaStore) createFunction(fn *profile.Function) (int64, error) {
	stmt, err := s.db.Prepare(
		`INSERT INTO "functions" (
                         function_id, name, system_name, filename, start_line
                         ) values(?,?,?,?,?)`,
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(int64(fn.ID), fn.Name, fn.SystemName, fn.Filename, fn.StartLine)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *sqlMetaStore) createLines(lines []profile.Line, locID int64) error {
	if len(lines) > 0 {
		q := `INSERT INTO "lines" (line, function_id) VALUES `
		for i, ln := range lines {
			functionID, err := s.getOrCreateFunction(ln.Function)
			if err != nil {
				return err
			}
			q += fmt.Sprintf(`(%s, %s)`,
				strconv.FormatInt(ln.Line, 10),
				strconv.FormatInt(functionID, 10))
			if i != len(lines)-1 {
				q += ", "
			}
		}
		q += ";"
		stmt, err := s.db.Prepare(q)
		if err != nil {
			return err
		}
		defer stmt.Close()

		res, err := stmt.Exec()
		if err != nil {
			return err
		}

		// Assuming ids are auto-incremented, we populate locations_lines going backwards.
		rf, err := res.RowsAffected()
		if err != nil {
			return err
		}
		lastLineID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		q = `INSERT INTO "location_lines" (line_id, location_id) VALUES `
		for i := int64(0); i < rf; i++ {
			q += fmt.Sprintf(`(%s, %s)`,
				strconv.FormatInt(lastLineID-i, 10),
				strconv.FormatInt(locID, 10))
			if i != rf-1 {
				q += ", "
			}
		}
		q += ";"
		stmt, err = s.db.Prepare(q)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
