// Copyright 2022 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logictest

import (
	"fmt"
	"github.com/dolthub/sqllogictest/go/logictest/parser"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"strings"
)

var statementPrefixes = []string{
	"ALTER ",
	"ALTER TABLE ",
	"ALTER DATABASE ",
	"ALTER VIEW ",
	"ALTER EVENT ",
	"ALTER TABLESPACE ",
	"ALTER UNDO TABLESPACE ",
	"ALTER LOGFILE GROUP ",
	"ALTER SERVER ",
	"CREATE ",
	"CREATE TABLE ",
	"CREATE DATABASE ",
	"CREATE FUNCTION ",
	"CREATE PROCEDURE ",
	"CREATE VIEW ",
	"CREATE TRIGGER ",
	"CREATE INDEX ",
	"CREATE UNIQUE INDEX ",
	"CREATE ROLE ",
	"DROP ",
	"DROP INDEX ",
	"DROP PROCEDURE ",
	"DROP TABLE ",
	"DROP TRIGGER ",
	"DROP VIEW ",
	"RENAME TABLE ",
	"RENAME TABLES ",
	"TRUNCATE TABLE ",
	"CALL ",
	"DELETE ",
	"DO ",
	"LOAD ",
	"SELECT ",
	"SHOW ",
	"PREPARE ",
	"INSERT ",
	"UPDATE ",
}

// AnalyzeStatements looks at all the tests in the specified test filepaths and prints out basic stats on statement
// type usage.
func AnalyzeStatements(harness Harness, paths ...string) {
	testFilepaths := collectTestFiles(paths)

	statementCounts := map[string]int{}
	for _, filepath := range testFilepaths {
		analyzeStatementsFromTestFile(harness, filepath, statementCounts)
	}

	p := message.NewPrinter(language.English)
	fmt.Println("Statement Counts:")
	for _, prefix := range statementPrefixes {
		value, ok := statementCounts[strings.TrimSpace(prefix)]
		if !ok {
			value = 0
		}

		p.Printf(" -  %-25s: %12d\n", prefix, value)
	}
}

func analyzeStatementsFromTestFile(harness Harness, testFilepath string, statementCounts map[string]int) {
	err := harness.Init()
	if err != nil {
		panic(err)
	}

	testFile, err := os.Open(testFilepath)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = testFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	testRecords, err := parser.ParseTestFile(testFilepath)
	if err != nil {
		panic(err)
	}

	for _, record := range testRecords {
		if record.ShouldExecuteForEngine("mysql") == false {
			continue
		}

		for _, prefix := range statementPrefixes {
			if strings.HasPrefix(record.Query(), prefix) {
				statementCounts[strings.TrimSpace(prefix)]++
			}
		}
	}
}
