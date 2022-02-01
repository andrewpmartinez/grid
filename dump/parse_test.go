package dump

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

var routine1 = `intentional garbage
intentional garbage
goroutine 733794 [semacquire, 792 minutes]:
sync.runtime_SemacquireMutex(0xc00034b3ec, 0xc00380ac00, 0x1)
	/opt/hostedtoolcache/go/1.16.4/x64/src/runtime/sema.go:71 +0x47
sync.(*Mutex).lockSlow(0xc00034b3e8)
	/opt/hostedtoolcache/go/1.16.4/x64/src/sync/mutex.go:138 +0x105
sync.(*Mutex).Lock(...)
	/opt/hostedtoolcache/go/1.16.4/x64/src/sync/mutex.go:81
go.etcd.io/bbolt.(*DB).beginTx(0xc00034b200, 0xc001d16860, 0x411f5b, 0xc0036d8cc0)
	/home/runner/go/pkg/mod/github.com/openziti/bbolt@v1.3.6-0.20210317142109-547da822475e/db.go:559 +0x2b3
go.etcd.io/bbolt.(*DB).Begin(0xc00034b200, 0x0, 0x30, 0x7f652175ab70, 0x58)
	/home/runner/go/pkg/mod/github.com/openziti/bbolt@v1.3.6-0.20210317142109-547da822475e/db.go:552 +0x6f
go.etcd.io/bbolt.(*DB).View(0xc00034b200, 0xc0036d8cc0, 0x0, 0x0)
	/home/runner/go/pkg/mod/github.com/openziti/bbolt@v1.3.6-0.20210317142109-547da822475e/db.go:726 +0x4c
github.com/openziti/fabric/controller/db.(*Db).View(0xc000011050, 0xc0036d8cc0, 0xc000011050, 0x30)
	/home/runner/go/pkg/mod/github.com/openziti/fabric@v0.16.62/controller/db/db.go:73 +0x38
github.com/openziti/edge/controller/model.(*baseHandler).readEntityWithIndex(0xc000e1a440, 0x16f2fc5, 0x5, 0xc00130dec0, 0x24, 0x30, 0x1a23840, 0xc000d8c450, 0x1a4e528, 0xc000fa95f0, ...)
	/home/runner/go/pkg/mod/github.com/openziti/edge@v0.19.113/controller/model/base_handler.go:320 +0x114
github.com/openziti/edge/controller/model.(*ApiSessionHandler).ReadByToken(0xc000e1a440, 0xc00130de90, 0x24, 0x29, 0xc001d16a08, 0x1)
	/home/runner/go/pkg/mod/github.com/openziti/edge@v0.19.113/controller/model/api_session_handlers.go:69 +0xfa
github.com/openziti/edge/controller/env.(*AppEnv).FillRequestContext(0xc000db9d40, 0xc00334c000, 0xc0036d8b40, 0xc00075ad00)
	/home/runner/go/pkg/mod/github.com/openziti/edge@v0.19.113/controller/env/appenv.go:245 +0x652
github.com/openziti/edge/controller/server.ClientApiHandler.newHandler.func1(0x1a3c7a0, 0xc0036d8b40, 0xc00075ad00)
	/home/runner/go/pkg/mod/github.com/openziti/edge@v0.19.113/controller/server/client-api.go:168 +0x27d
net/http.HandlerFunc.ServeHTTP(0xc001ccfb60, 0x1a3c7a0, 0xc0036d8b40, 0xc00075ad00)
	/opt/hostedtoolcache/go/1.16.4/x64/src/net/http/server.go:2069 +0x44
github.com/gorilla/handlers.(*cors).ServeHTTP(0xc0026978c0, 0x1a3c7a0, 0xc0036d8b40, 0xc00075ad00)
	/home/runner/go/pkg/mod/github.com/gorilla/handlers@v1.5.1/cors.go:54 +0x103e
github.com/openziti/edge/controller/timeout.(*timeoutHandler).ServeHTTP.func1(0xc0036d8ba0, 0xc0011c2cc0, 0xc0036d8b40, 0xc00075ad00, 0xc00235f980)
	/home/runner/go/pkg/mod/github.com/openziti/edge@v0.19.113/controller/timeout/timeouts.go:72 +0x7f
created by github.com/openziti/edge/controller/timeout.(*timeoutHandler).ServeHTTP
	/home/runner/go/pkg/mod/github.com/openziti/edge@v0.19.113/controller/timeout/timeouts.go:65 +0x27b

goroutine 1 [select]:
github.com/openziti/fabric/controller/network.(*Network).Run(0xc00016c420)
	/home/runner/go/pkg/mod/github.com/openziti/fabric@v0.16.62/controller/network/network.go:669 +0x225
github.com/openziti/fabric/controller.(*Controller).Run(0xc000275c30, 0x18c8fc8, 0xc000275c30)
	/home/runner/go/pkg/mod/github.com/openziti/fabric@v0.16.62/controller/controller.go:138 +0x845
github.com/openziti/ziti/ziti-controller/subcmd.run(0x26be720, 0xc000733bf0, 0x1, 0x1)
	/home/runner/work/ziti/ziti/ziti-controller/subcmd/run.go:76 +0x838
github.com/spf13/cobra.(*Command).execute(0x26be720, 0xc000733bd0, 0x1, 0x1, 0x26be720, 0xc000733bd0)
	/home/runner/go/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:856 +0x2c2
github.com/spf13/cobra.(*Command).ExecuteC(0x26be4a0, 0xc00006c740, 0x471205, 0xc000000180)
	/home/runner/go/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:960 +0x375
github.com/spf13/cobra.(*Command).Execute(...)
	/home/runner/go/pkg/mod/github.com/spf13/cobra@v1.1.3/command.go:897
github.com/openziti/ziti/ziti-controller/subcmd.Execute()
	/home/runner/work/ziti/ziti/ziti-controller/subcmd/root.go:64 +0x31
main.main()
	/home/runner/work/ziti/ziti/ziti-controller/main.go:49 +0x25`

func Test_ParseEntireGoRoutine(t *testing.T) {
	stringBuffer := bytes.NewBufferString(routine1)
	buffReader := bufio.NewReader(stringBuffer)
	scanner := bufio.NewScanner(buffReader)

	t.Run("does not return errors on valid go routine", func(t *testing.T) {
		req := require.New(t)
		dump, err := ParseScanner(scanner, nil)
		req.NoError(err)

		t.Run("has two go routines", func(t *testing.T) {
			req := require.New(t)
			req.Len(dump.Routines, 2)
		})

		t.Run("first go routine 0 has 15 frames", func(t *testing.T) {
			req := require.New(t)
			req.Len(dump.Routines[0].Frames, 15, "expected %d frames got %d", 15, len(dump.Routines[0].Frames))
		})
		//goroutine 733794 [semacquire, 792 minutes
		t.Run("go routine 0 id to be 733794", func(t *testing.T) {
			req := require.New(t)
			req.Equal(733794, dump.Routines[0].Id)
		})

		t.Run("go routine 0 type to be semacquire", func(t *testing.T) {
			req := require.New(t)
			req.Equal("semacquire", dump.Routines[0].Type)
		})

		t.Run("go routine 0 duration is 792", func(t *testing.T) {
			req := require.New(t)
			req.Equal("792", dump.Routines[0].Duration)
		})

		t.Run("go routine 0 duration unit is minutes", func(t *testing.T) {
			req := require.New(t)
			req.Equal("minutes", dump.Routines[0].DurationUnit)
		})

		t.Run("go routine 0 file start line is 3", func(t *testing.T) {
			req := require.New(t)
			req.Equal(3, dump.Routines[0].FileStartLine)
		})

		t.Run("go routine 0 file end line is 33", func(t *testing.T) {
			req := require.New(t)
			req.Equal(33, dump.Routines[0].FileEndLine)
		})

		t.Run("go routine 0 frame 15 has start of line 32", func(t *testing.T) {
			req := require.New(t)
			req.Equal(32, dump.Routines[0].Frames[14].FileStartLine)
		})

		t.Run("go routine 0 frame 15 has end of line 33", func(t *testing.T) {
			req := require.New(t)
			req.Equal(33, dump.Routines[0].Frames[14].FileEndLine)
		})

		t.Run("go routine 0 frame 0 (no struct context)", func(t *testing.T) {
			t.Run("has the correct fully qualified function name", func(t *testing.T) {
				req := require.New(t)
				req.Equal("sync.runtime_SemacquireMutex", dump.Routines[0].Frames[0].Function)
			})

			t.Run("has the correct function argument addresses", func(t *testing.T) {
				req := require.New(t)
				req.Equal("0xc00034b3ec", dump.Routines[0].Frames[0].ArgumentAddresses[0])
				req.Equal("0xc00380ac00", dump.Routines[0].Frames[0].ArgumentAddresses[1])
				req.Equal("0x1", dump.Routines[0].Frames[0].ArgumentAddresses[2])
			})

			t.Run("has the correct location file path", func(t *testing.T) {
				req := require.New(t)
				req.Equal("/opt/hostedtoolcache/go/1.16.4/x64/src/runtime/sema.go", dump.Routines[0].Frames[0].Path)
			})

			t.Run("has the correct location file line number", func(t *testing.T) {
				req := require.New(t)
				req.Equal(71, dump.Routines[0].Frames[0].Line)
			})

			t.Run("has the correct location file offset", func(t *testing.T) {
				req := require.New(t)
				req.Equal("0x47", dump.Routines[0].Frames[0].Offset)
			})

			t.Run("has the correct location file start line", func(t *testing.T) {
				req := require.New(t)
				req.Equal(4, dump.Routines[0].Frames[0].FileStartLine)
			})

			t.Run("has the correct location file end line", func(t *testing.T) {
				req := require.New(t)
				req.Equal(5, dump.Routines[0].Frames[0].FileEndLine)
			})
		})

		t.Run("go routine 0 frame 3 (with struct context)", func(t *testing.T) {

			t.Run("has the correct fully qualified function name", func(t *testing.T) {
				req := require.New(t)
				req.Equal("go.etcd.io/bbolt.", dump.Routines[0].Frames[3].Function)
			})

			t.Run("has the correct struct context", func(t *testing.T) {
				req := require.New(t)
				req.Equal("*DB", dump.Routines[0].Frames[3].StructContext)
			})

			t.Run("has the correct struct function", func(t *testing.T) {
				req := require.New(t)
				req.Equal(".beginTx", dump.Routines[0].Frames[3].StructContextFunction)
			})

			t.Run("has the correct function argument addresses", func(t *testing.T) {
				req := require.New(t)
				req.Equal("0xc00034b200", dump.Routines[0].Frames[3].ArgumentAddresses[0])
				req.Equal("0xc001d16860", dump.Routines[0].Frames[3].ArgumentAddresses[1])
				req.Equal("0x411f5b", dump.Routines[0].Frames[3].ArgumentAddresses[2])
				req.Equal("0xc0036d8cc0", dump.Routines[0].Frames[3].ArgumentAddresses[3])
			})

			t.Run("has the correct location file path", func(t *testing.T) {
				req := require.New(t)
				req.Equal("/home/runner/go/pkg/mod/github.com/openziti/bbolt@v1.3.6-0.20210317142109-547da822475e/db.go", dump.Routines[0].Frames[3].Path)
			})

			t.Run("has the correct location file line number", func(t *testing.T) {
				req := require.New(t)
				req.Equal(559, dump.Routines[0].Frames[3].Line)
			})

			t.Run("has the correct location file offset", func(t *testing.T) {
				req := require.New(t)
				req.Equal("0x2b3", dump.Routines[0].Frames[3].Offset)
			})

			t.Run("has the correct location file start line", func(t *testing.T) {
				req := require.New(t)
				req.Equal(10, dump.Routines[0].Frames[3].FileStartLine)
			})

			t.Run("has the correct location file end line", func(t *testing.T) {
				req := require.New(t)
				req.Equal(11, dump.Routines[0].Frames[3].FileEndLine)
			})
		})
	})
}
