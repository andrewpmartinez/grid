package dump

import (
	"fmt"
	"github.com/elliotchance/orderedmap"
	"time"
)

type Stats struct {
	RoutinesByType     *orderedmap.OrderedMap // type -> []*RoutineStats
	RoutinesByFunction *orderedmap.OrderedMap // funcName -> []*RoutineStats
}

type RoutineStats struct {
	Routines            []*Routine
	RoutinesBySignature map[string][]*Routine //sig -> example routine
}

type LocationStats struct {
	Count        int
	ShortestWait time.Duration
	LongestWait  time.Duration
}

func (stats *Stats) Collect(routine *Routine) {
	if typeVal, _ := stats.RoutinesByType.Get(routine.Type); typeVal == nil {
		routineStats := &RoutineStats{
			Routines: []*Routine{routine},
			RoutinesBySignature: map[string][]*Routine{
				routine.StackSignature: {routine},
			},
		}
		stats.RoutinesByType.Set(routine.Type, routineStats)
	} else {
		routineStats, ok := typeVal.(*RoutineStats)

		if !ok {
			panic(fmt.Sprintf("could not convert %T to %T", typeVal, routineStats))
		}

		routineStats.Routines = append(routineStats.Routines, routine)
		routineStats.RoutinesBySignature[routine.StackSignature] = append(routineStats.RoutinesBySignature[routine.StackSignature], routine)

		stats.RoutinesByType.Set(routine.Type, routineStats)
	}

	if typeVal, _ := stats.RoutinesByFunction.Get(routine.Frames[0].Function); typeVal == nil {
		routineStats := &RoutineStats{
			Routines: []*Routine{routine},
			RoutinesBySignature: map[string][]*Routine{
				routine.StackSignature: {routine},
			},
		}
		stats.RoutinesByFunction.Set(routine.Frames[0].Function, routineStats)
	} else {
		routineStats, ok := typeVal.(*RoutineStats)

		if !ok {
			panic(fmt.Sprintf("could not convert %T to %T", typeVal, routineStats))
		}

		routineStats.Routines = append(routineStats.Routines, routine)
		routineStats.RoutinesBySignature[routine.StackSignature] = append(routineStats.RoutinesBySignature[routine.StackSignature], routine)

		stats.RoutinesByFunction.Set(routine.Frames[0].Function, routineStats)
	}
}

func (stats *Stats) GetRoutinesByFunction(name string) *RoutineStats {
	var ret *RoutineStats

	if typeVal, _ := stats.RoutinesByFunction.Get(name); typeVal != nil {
		ret, _ = typeVal.(*RoutineStats)
	}

	return ret
}
