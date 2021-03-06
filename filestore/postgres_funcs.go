/*
 * Copyright (c) 2013-2014, Jeremy Bingham (<jbingham@gmail.com>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/* Postgres specific functions for file store */

package filestore

import (
	"database/sql"
	"github.com/ctdk/goas/v2/logger"
	"github.com/ctdk/goiardi/datastore"
	"log"
	"strings"
)

func (f *FileStore) savePostgreSQL() error {
	tx, err := datastore.Dbh.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO goiardi.file_checksums (organization_id, checksum) VALUES (1, $1)", f.Chksum)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func deleteHashesPostgreSQL(fileHashes []string) {
	if len(fileHashes) == 0 {
		return // nothing to do
	}
	tx, err := datastore.Dbh.Begin()
	if err != nil {
		log.Fatal(err)
	}
	deleteQuery := "DELETE FROM goiardi.file_checksums WHERE checksum = ANY($1::varchar(32)[])"
	_, err = tx.Exec(deleteQuery, "{"+strings.Join(fileHashes, ",")+"}")
	if err != nil && err != sql.ErrNoRows {
		logger.Debugf("Error %s trying to delete hashes", err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
