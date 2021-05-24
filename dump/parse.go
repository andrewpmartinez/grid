package dump

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Dump struct {
	Routines []*Routine
	Stats    Stats
}

type Stats struct {
	RoutinesByType     map[string]*Routine
	RoutinesByLocation map[string][]*Routine
}

type Routine struct {
	*RoutineLine
	Frames []*Frame

	FileStartLine int
	FileEndLine   int
}

type Frame struct {
	Function              string
	ArgumentAddresses     []string
	StructContext         string
	StructContextFunction string
	Path                  string
	Offset                string
	Line                  int

	FileStartLine int
	FileEndLine   int
}

func ParseFile(filePath string, logger Logger) (*Dump, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("could not open file [%s]: %v", filePath, err)
	}

	fileScanner := bufio.NewScanner(file)
	return ParseScanner(fileScanner, logger)
}

func ParseScanner(scanner *bufio.Scanner, logger Logger) (*Dump, error) {
	ctx := context{
		dump:           &Dump{},
		currentRoutine: nil,
		logger:         logger,
	}

	if ctx.logger == nil {
		ctx.logger = logrus.New()
	}

	for scanner.Scan() {
		ctx.NextLine(scanner.Text())
	}

	ctx.Done()

	return ctx.dump, nil
}
