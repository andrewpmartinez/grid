package dump

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Logger interface {
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}

type context struct {
	dump             *Dump
	currentRoutine   *Routine
	currentFrame     *Frame
	lineNumber       int
	logger           Logger
	lastLineWasBlank bool
	foundStartOfDump bool
}

func (c *context) Done() {
	if c.currentFrame != nil {
		if c.currentFrame.FileEndLine == -1 {
			c.currentFrame.FileEndLine = c.lineNumber
		}
		c.currentRoutine.Frames = append(c.currentRoutine.Frames, c.currentFrame)
	}

	if c.currentRoutine != nil {
		c.dump.Routines = append(c.dump.Routines, c.currentRoutine)
		if c.currentRoutine.FileEndLine == -1 {
			c.currentRoutine.FileEndLine = c.lineNumber
		}
	}
}

func (c *context) NextLine(line string) {
	c.lineNumber = c.lineNumber + 1
	log.Default()

	trimmedLine := strings.TrimSpace(line)
	isGoRoutineLine := strings.HasPrefix(trimmedLine, "goroutine")

	if !c.foundStartOfDump {
		if isGoRoutineLine {
			c.foundStartOfDump = true
			c.lastLineWasBlank = true
		} else {
			return
		}
	}

	if isGoRoutineLine {
		c.handleGoRoutine(line)
	} else {
		if LocationLineMatch.MatchString(line) {
			c.handleLocation(line)
		} else if FunctionLineMatch.MatchString(line) {
			c.handleFunction(line)
		} else {
			if trimmedLine == "" {
				if c.currentFrame != nil {
					c.currentFrame.FileEndLine = c.lineNumber - 1
				}
			} else {
				c.logger.Errorf("unexpected line [%d] did not match any known go routine line types: %s", c.lineNumber, line)
			}
		}

		c.lastLineWasBlank = trimmedLine == ""
	}
}

// handleGoRoutine deals with go routine header dumps and primes the context to parse functions and locations. Go routine
// header examples are as follows:
// - goroutine 1 [select, 1 minutes]:
// - goroutine 9 [syscall, 6829 minutes]:
// - goroutine 10 [chan receive, 6829 minutes]:
func (c *context) handleGoRoutine(line string) {
	if c.currentRoutine != nil {
		if !c.lastLineWasBlank {
			c.logger.Warnf("unexpected go routine on line [%d] already in go routine found at [%d]", c.lineNumber, c.currentRoutine.FileStartLine)
			if c.lineNumber > 0 {
				c.currentRoutine.FileEndLine = c.lineNumber - 1 //current line only
			}
		} else if c.lineNumber > 0 {
			c.currentRoutine.FileEndLine = c.lineNumber - 2 //current line then blank line
		}

		c.goRoutineComplete(c.currentRoutine)
		c.currentRoutine = nil
		c.currentFrame = nil
	}

	routineLine, err := ParseRoutineLine(line)
	if err != nil {
		c.logger.Errorf("could not parse go routine on line [%d]: %v", c.lineNumber, err)
	} else {
		c.currentRoutine = &Routine{
			RoutineLine:   routineLine,
			Frames:        nil,
			FileStartLine: c.lineNumber,
			FileEndLine:   -1,
			LineText:      line,
			allBuilder:    strings.Builder{},
		}

		c.currentRoutine.allBuilder.Write([]byte(line))
		c.currentRoutine.allBuilder.Write([]byte("\n"))
	}
}

// handleLocation deals with locations that provide more information to function lines. Example location lines are as
// follows:
// - go.etcd.io/bbolt@v1.3.5-0.20200615073812-232d8fc87f50/db.go:726 +0x45
// - github.com/go-openapi/runtime@v0.21.0/middleware/spec.go:46 +0x18c
//
// Location lines start with tabs and contain a specific file, line numbers, and offset addresses. Locations augment
// the current "frame" which is some function being executed. The function the location belongs to is immediately
// above it in the dump.
func (c *context) handleLocation(line string) {
	if c.currentFrame == nil {
		c.logger.Warnf("unexpected location line [%d] not in a function, expected blank line or new function", c.lineNumber)
	} else if c.currentFrame.Line != -1 {
		c.logger.Warnf("unexpected location line [%d] location already parsed for function starting at [%d]", c.lineNumber, c.currentFrame.FileStartLine)
	} else {
		if location, err := ParseLocationLine(line); err != nil {
			c.logger.Errorf("could not parse location at [%d]: %v", c.lineNumber, err)
		} else {
			c.currentFrame.FileEndLine = c.lineNumber
			c.currentFrame.UniqueId = location.UniqueId
			c.currentFrame.Path = location.File
			c.currentFrame.Offset = location.OffsetAddress
			c.currentFrame.Line = location.Line
			c.currentRoutine.allBuilder.Write([]byte(line))
			c.currentRoutine.allBuilder.Write([]byte("\n"))
			c.currentFrame.LocationLineText = line
			c.currentRoutine.Frames = append(c.currentRoutine.Frames, c.currentFrame)
			c.currentFrame = nil
		}
	}
}

// handleFunction deals with function lines in a dump file. Example function lines are as follows:
// - github.com/go-openapi/runtime/middleware.Spec.func1({0x1a87090, 0xc002ac13e0}, 0xc001381c20)
// - github.com/gorilla/handlers.(*cors).ServeHTTP(0xc000b265a0, {0x1a87090, 0xc002ac13e0}, 0xc000791600)
//
// Function lines are found under a go routine header and precede a location line. Function lines represent one frame
// in the call stack of a go routine. They are lacking some execution info which is provided in the dump files
// // next line, which is a location line.
func (c *context) handleFunction(line string) {
	if c.currentFrame != nil {
		c.logger.Warnf("unexpected function on line [%d] already in a function, expected a location for function from line [%d]", c.lineNumber, c.currentFrame.FileStartLine)
	} else {
		function, err := ParseFunctionLine(line)

		if err != nil {
			c.logger.Errorf("could not parse function at line [%d]: %v", c.lineNumber, err)
		} else {
			c.currentRoutine.allBuilder.Write([]byte(line))
			c.currentRoutine.allBuilder.Write([]byte("\n"))
			c.currentFrame = &Frame{
				Function:              function.FullyQualifiedName,
				ArgumentAddresses:     function.ArgAddresses,
				StructContext:         function.StructContext,
				StructContextFunction: function.StructFunction,
				Path:                  "",
				Offset:                "",
				Line:                  -1,
				FileStartLine:         c.lineNumber,
				FileEndLine:           -1,
				FunctionLineText:      line,
			}
		}
	}
}

// goRoutineComplete handles cataloging the routines for later consumption.
func (c *context) goRoutineComplete(routine *Routine) {
	routine.calculateStackSignature()
	c.dump.Routines = append(c.dump.Routines, c.currentRoutine)

	c.dump.Stats.Collect(routine)
}

type RoutineLine struct {
	Id           int
	Type         string
	Duration     string
	DurationUnit string
}

// RoutineLineMatch is a regex that matches a goroutine line.
// Check the `waitReasonStrings` variable in `src/runtime/runtime2.go` for the possible wait reasons.
var RoutineLineMatch = regexp.MustCompile(`^goroutine (\d+) \[([\w \.\(\)]+)(, (\d+) (\w+))?]:$`)

var (
	RoutineLineMatchId           = 1
	RoutineLineMatchType         = 2
	RoutineLineMatchDuration     = 4
	RoutineLineMatchDurationUnit = 5
	RoutineLineMatchLength       = 6
)

func ParseRoutineLine(line string) (*RoutineLine, error) {
	matches := RoutineLineMatch.FindStringSubmatch(line)

	if len(matches) != RoutineLineMatchLength {
		return nil, fmt.Errorf("could not parse go routine, invalid submatches [%d] expected [%d]", len(matches), RoutineLineMatchLength)
	}

	id, _ := strconv.ParseInt(matches[RoutineLineMatchId], 10, 32)

	return &RoutineLine{
		Id:           int(id),
		Type:         matches[RoutineLineMatchType],
		Duration:     matches[RoutineLineMatchDuration],
		DurationUnit: matches[RoutineLineMatchDurationUnit],
	}, nil
}

type FunctionLine struct {
	FullyQualifiedName string
	ArgAddresses       []string
	Location           *LocationLine
	StructContext      string
	StructFunction     string
}

var FunctionLineMatch = regexp.MustCompile(`([ \w\.\/]+)(\((.*?)\))?([\w\.\/]+)?(\((.*?)\))?$`)
var (
	FunctionLineMatchFullyQualifiedName = 1
	FunctionLineStructContext           = 3
	FunctionLineStructFunction          = 4
	FunctionLineArgAddresses            = 6
	FunctionLineMatchLength             = 7
)

func ParseFunctionLine(line string) (*FunctionLine, error) {
	matches := FunctionLineMatch.FindStringSubmatch(line)

	if len(matches) != FunctionLineMatchLength {
		return nil, fmt.Errorf("could not parse function line, invalid submatches [%d] expected [%d]", len(matches), FunctionLineMatchLength)
	}

	if matches[6] == "" {
		args := strings.Split(matches[3], ", ")
		return &FunctionLine{
			FullyQualifiedName: matches[FunctionLineMatchFullyQualifiedName],
			ArgAddresses:       args,
			Location:           nil,
		}, nil
	}

	args := strings.Split(matches[FunctionLineArgAddresses], ", ")
	return &FunctionLine{
		FullyQualifiedName: matches[FunctionLineMatchFullyQualifiedName],
		StructContext:      matches[FunctionLineStructContext],
		StructFunction:     matches[FunctionLineStructFunction],
		ArgAddresses:       args,
		Location:           nil,
	}, nil

}

type LocationLine struct {
	UniqueId      string
	File          string
	Line          int
	OffsetAddress string
}

var LocationLineMatch = regexp.MustCompile(`^\t(.*?\.go):(\d+)( \+(.*?))?$`)
var (
	LocationLineMatchedFile          = 1
	LocationLineMatchedLine          = 2
	LocationLineMatchedOffsetAddress = 4
	LocationLineMatchedLength        = 5
)

func ParseLocationLine(line string) (*LocationLine, error) {
	matches := LocationLineMatch.FindStringSubmatch(line)

	if len(matches) != LocationLineMatchedLength {
		return nil, fmt.Errorf("could not parse location line, invalid submatches [%d] expected [%d]", len(matches), LocationLineMatchedLength)
	}

	var lineNumber = -1

	if num, err := strconv.ParseInt(matches[LocationLineMatchedLine], 10, 64); err == nil {
		lineNumber = int(num)
	}

	return &LocationLine{
		UniqueId:      matches[LocationLineMatchedFile] + ":" + matches[LocationLineMatchedLine],
		File:          matches[LocationLineMatchedFile],
		Line:          lineNumber,
		OffsetAddress: matches[LocationLineMatchedOffsetAddress],
	}, nil
}
