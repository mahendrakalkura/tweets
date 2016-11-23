package main

import (
	"github.com/jmoiron/sqlx"
)

func get_programs(database *sqlx.DB) []Program {
	var programs []Program
	query := `
    SELECT
        programs.query,
        channels_programs.beginning_at,
        channels_programs.ending_at
    FROM programs
    INNER JOIN channels_programs ON channels_programs.program_id = programs.id
    WHERE
        channels_programs.beginning_at <= TIMEZONE('UTC', NOW())
        AND
        channels_programs.ending_at >= TIMEZONE('UTC', NOW())
    `
	rows, err := database.Query(query)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var program Program
		err = rows.Scan(&program)
		if err != nil {
			panic(err)
		}
		programs = append(programs, program)
	}
	return programs
}
