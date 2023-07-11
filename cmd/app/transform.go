package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"
)

type transformCmd struct {
	// cli options
	InputFilepath  string `help:"path to the file to transform" type:"existingfile" default:"./internal/adapters/namesdb/names.csv"`
	OutputFilepath string `help:"path to the file to write the transformed data to" type:"string" default:"./internal/adapters/namesdb/names_transformed.csv"`
	// Dependencies
	logger *slog.Logger
}

func (c *transformCmd) Run(cmdCtx *cmdContext) error {
	c.logger = cmdCtx.Logger.With("component", "transformCmd")

	// open source file
	fd, err := os.Open(c.InputFilepath)
	if err != nil {
		return fmt.Errorf("failed to open a source file descriptor: %w", err)
	}
	defer fd.Close()
	bfd := bufio.NewReader(fd)

	// open destination file
	ofd, err := os.OpenFile(c.OutputFilepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open a destination file descriptor: %w", err)
	}
	defer ofd.Close()
	bofd := bufio.NewWriter(ofd)
	defer bofd.Flush()

	// read the file line by line
	var lastId int64
	r := csv.NewReader(bfd)
	r.Read() // skip the header
FILE_READING_LOOP:
	for {
		record, err := r.Read()
		if err != nil {
			switch {
			case err == io.EOF:
				break FILE_READING_LOOP
			default:
				return fmt.Errorf("failed to read source file %s: %w", c.InputFilepath, err)
			}
		}
		c.logger.Debug("read a record", "record", record)
		name := record[0]
		lastId++
		// write the transformed data to the destination file
		_, err = bofd.WriteString(fmt.Sprintf("%d,%s\n", lastId, name))
		if err != nil {
			return fmt.Errorf("failed to write to destination file %s: %w", c.OutputFilepath, err)
		}
		c.logger.Debug("transformed a record", "id", lastId, "name", name)
	}

	return nil
}
