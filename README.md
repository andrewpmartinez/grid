# grid
GRID - Go Routine Inspect Dump parses Go routine dumps. It provides both a go module to parse dump files (`dump/`) and 
an executable GUI to show the results of the parse (`cmd/grid`).

The parser will provide broken down information from the go routine frames (such as line numbers, file names, etc.) as
well as provide de-duped go routines based on stack trace signatures.

![GRID Screenshot](images/screenshot.png)

## The Parser: Dump

The dump module can take in a clean go routine dump file and produces a `Dump` struct. It can be invoked by the 
`dump.ParseFile()` function which has the following signature:

```
ParseFile(filePath string, logger Logger) (*Dump, error)
```

With an example invocation of:

```
    dump, err := dump.ParseFile(pathToDumpFile, logrus.New())
```

The first argument, `filePath`, is simply the path to the go routine dump file. The second argument, `logger`, is
any logger you wish that satisfies the interface requirements (e.g. logrus).

The result of a parse is either error or a `Dump` struct.

```
type Dump struct {
	Routines []*Routine
	Stats    *Stats
}
```

`Dump` structs has the raw parsed `Routine`s as well as a `Stats` struct that contain structs that have been organized
by function and type as well as having duplicate stacks reduced to 1 entry. All the raw routines still are available
inside the `Stats` struct if needed.


## The GUI: GRID

The GRID GUI can be invoked by building this project and running:

```
> grid gui <path.to.dump>
```


## Building This Project

This project uses [AllenDang's GIU](https://github.com/AllenDang/giu) library and is subject to that project's build
requirements. See your OS's requirements [here](https://github.com/AllenDang/giu#install).

Once complete, building GRID can be accomplished as follows:

1) Check out the repository
2) Navigate to checkout repository
3) Install via `go`
```
go install ./...
```