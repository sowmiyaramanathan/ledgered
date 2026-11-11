package main

import (
	"fmt"
	"io"
	"ledgerd/internal/entities"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&entities.Account{},
		&entities.Entry{},
		&entities.Transaction{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to load gorm schema: %v", err))
	}
	io.WriteString(os.Stdout, stmts)
}
