package dump

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"github.com/elliotchance/orderedmap"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Dump struct {
	Routines []*Routine
	Stats    *Stats
}

type Routine struct {
	*RoutineLine
	Frames []*Frame

	FileStartLine  int
	FileEndLine    int
	LineText       string
	AllLines       []string
	allBuilder     strings.Builder
	StackSignature string
}

func (r *Routine) Raw() string {
	return r.allBuilder.String()
}

func (r *Routine) calculateStackSignature() {
	stringBuilder := strings.Builder{}

	for _, frame := range r.Frames {
		stringBuilder.Write([]byte(frame.UniqueId))
	}

	r.StackSignature = fmt.Sprintf("%x", sha256.Sum256([]byte(stringBuilder.String())))
}

type Frame struct {
	Function              string
	ArgumentAddresses     []string
	StructContext         string
	StructContextFunction string
	Path                  string
	Offset                string
	Line                  int
	FunctionLineText      string
	LocationLineText      string

	FileStartLine int
	FileEndLine   int
	UniqueId      string
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
		dump: &Dump{
			Routines: nil,
			Stats: &Stats{
				RoutinesByType:     orderedmap.NewOrderedMap(),
				RoutinesByFunction: orderedmap.NewOrderedMap(),
			},
		},
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
