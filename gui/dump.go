package gui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	g "github.com/AllenDang/giu"
	"github.com/andrewpmartinez/grid/dump"
	"github.com/sirupsen/logrus"
)

type DumpWindow struct {
	Dump         *dump.Dump
	masterWindow *g.MasterWindow

	buildFunctionRowsOnce sync.Once
	functionRows          []*g.TableRowWidget
	path                  string
	routineText           string
}

func NewDumpWindow() *DumpWindow {
	return &DumpWindow{
		buildFunctionRowsOnce: sync.Once{},
	}
}

func (dumpWindow *DumpWindow) LoadFile(path string) {
	var err error
	dumpWindow.path = path
	dumpWindow.Dump, err = dump.ParseFile(path, nil)

	if err != nil {
		errStr := fmt.Sprintf("error parsing file [%s]: %v", path, err)
		logrus.Error(errStr)
	}

}

func (dumpWindow *DumpWindow) buildRows() []*g.TableRowWidget {
	dumpWindow.buildFunctionRowsOnce.Do(func() {

		maxFuncName := ""
		maxCount := 0
		dumpWindow.functionRows = make([]*g.TableRowWidget, dumpWindow.Dump.Stats.RoutinesByFunction.Len())

		i := 0
		for el := dumpWindow.Dump.Stats.RoutinesByFunction.Front(); el != nil; el = el.Next() {
			routineStats, _ := el.Value.(*dump.RoutineStats)

			funcName := el.Key.(string)
			numTotalRoutines := len(routineStats.Routines)
			numUniqueRoutines := len(routineStats.RoutinesBySignature)

			if len(funcName) > len(maxFuncName) {
				maxFuncName = funcName
			}

			if numTotalRoutines > maxCount {
				maxCount = numTotalRoutines
			}

			button := g.Button("Open").OnClick(func() {
				dumpWindow.OpenFunctionDetail(funcName)
			})

			numTotalRoutinesLabel := g.Label(strconv.Itoa(numTotalRoutines))
			numUniqueRoutinesLabel := g.Label(strconv.Itoa(numUniqueRoutines))

			funcLabel := g.Label(funcName)

			dumpWindow.functionRows[i] = g.TableRow(
				button,
				numTotalRoutinesLabel,
				numUniqueRoutinesLabel,
				funcLabel,
			)

			i++
		}

		//		masterWindow.functionRows[0].BgColor(&(color.RGBA{200, 100, 100, 255}))
	})

	return dumpWindow.functionRows
}

func (dumpWindow *DumpWindow) loop() {
	masterWidth, masterHeight := dumpWindow.masterWindow.GetSize()

	g.Window("DumpStatus").
		Flags(g.WindowFlagsNoResize|g.WindowFlagsNoCollapse|g.WindowFlagsNoMove|g.WindowFlagsNoTitleBar).
		Size(float32(masterWidth), 30).
		Layout(
			g.Label("File: " + dumpWindow.path),
		)

	g.Window("DumpNav").
		Flags(g.WindowFlagsNoResize|g.WindowFlagsNoCollapse|g.WindowFlagsNoMove|g.WindowFlagsNoTitleBar).
		Pos(0, 31).
		Size(650, float32(math.Max(float64(masterHeight-31), 50))).
		Layout(
			g.Table().
				Columns(
					g.TableColumn("").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(50),
					g.TableColumn("Total").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(70),
					g.TableColumn("Unique").Flags(g.TableColumnFlagsWidthFixed).InnerWidthOrWeight(70),
					g.TableColumn("Function"),
				).
				Freeze(0, 1).
				FastMode(true).
				Rows(dumpWindow.buildRows()...),
		)

	g.Window("RoutineDetails").
		Flags(g.WindowFlagsNoResize|g.WindowFlagsNoCollapse|g.WindowFlagsNoMove|g.WindowFlagsNoTitleBar).
		Pos(650, 31).
		Size(-1, float32(math.Max(float64(masterHeight-31), 50))).
		Layout(
			g.InputTextMultiline(&dumpWindow.routineText).
				Size(float32(math.Max(float64(masterWidth-650), 400)), float32(math.Max(float64(masterHeight-31), 50))),
		)
}

func (dumpWindow *DumpWindow) Run() {
	dumpWindow.masterWindow = g.NewMasterWindow("Dump", 1700, 800, 0)
	dumpWindow.masterWindow.Run(dumpWindow.loop)
}

func (dumpWindow *DumpWindow) OpenFunctionDetail(funcName string) {
	routineStats := dumpWindow.Dump.Stats.GetRoutinesByFunction(funcName)

	builder := strings.Builder{}
	start := time.Now()
	for signature, routines := range routineStats.RoutinesBySignature {
		builder.WriteString(fmt.Sprintf("Signature: %s\nOccurences: %d \n\n", signature, len(routines)))

		builder.WriteString(routines[0].Raw())
		builder.WriteString("\n")
		builder.WriteString("go routine ids: ")
		isFirst := true
		idsPerLineCount := 0
		for _, routine := range routines {
			if !isFirst {
				builder.WriteString(",")
				if idsPerLineCount > 50 {
					idsPerLineCount = 0
					builder.WriteString("\n")
				}
			} else {
				isFirst = false
			}
			builder.WriteString(strconv.Itoa(routine.Id))
			idsPerLineCount++
		}
		builder.WriteString("\n\n--------------------------------------------------------------------------------------------------------------------\n\n")

	}

	dumpWindow.routineText = builder.String()

	duration := time.Now().Sub(start)
	println(duration.String())
}
