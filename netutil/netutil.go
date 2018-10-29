package netutil

import (
	"../gwlog"
	"fmt"
	"io"
	"net"
	"reflect"
	"runtime/debug"
)

func init() {

}

func IsTemporaryNetError(err error) bool {
	if err == nil {
		return false
	}
	netErr, ok := err.(net.Error)
	if !ok {
		return false
	}
	return netErr.Temporary() || netErr.Timeout()
}

func IsConnectionClosed(_err interface{}) bool {
	err, ok := _err.(error)
	if ok && err == io.EOF {
		return true
	}
	netErr, ok := _err.(net.Error)
	if !ok {
		return false
	}
	if netErr.Timeout() || netErr.Temporary() {
		return false
	}
	return true
}

func WriteAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Write(data)
		if n > 0 {
			data = data[n:]
		}
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

//maybe have some problem
func ReadAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Read(data)
		if n > 0 {
			data = data[n:]
		}
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func ReadLine(conn net.Conn) (string, error) {
	linebuff := make([]byte, 1)
	buff := [1]byte{0}
	for {
		n, err := conn.Read(buff[0:1])
		if err != nil {
			if IsTemporaryNetError(err) {
				continue
			} else {
				return "", err
			}
		}
		if n == 1 {
			c := buff[0]
			if c == '\n' {
				return string(linebuff), err
			} else {
				linebuff = append(linebuff, c)
			}
		}
	}
}

func ConnectTCP(host string, port int) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	return net.Dial("tcp", addr)
}

func ServeForever(f interface{}, args ...interface{}) {
	fval := reflect.ValueOf(f)
	argCnt := len(args)
	argVals := make([]reflect.Value, argCnt, argCnt)
	for i := 0; i < argCnt; i++ {
		argVals[i] = reflect.ValueOf(args[i])
	}
	for {
		runServe(fval, argVals)
	}
}

func runServe(f reflect.Value, args []reflect.Value) {
	defer func() {
		err := recover()
		if err != nil {
			gwlog.Error("ServeForever:func %v quited with error %v", f, err)
			debug.PrintStack()
		}
	}()
	rets := f.Call(args)
	gwlog.Debug("ServeForever:func %v returns %v", f, rets)
}
