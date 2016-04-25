// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"time"

	"github.com/limetext/backend"
	"github.com/limetext/backend/log"
	"github.com/limetext/gopy/lib"
)

func sublime_ErrorMessage(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for backend.DummyFrontend.ErrorMessage() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	backend.GetEditor().Frontend().ErrorMessage(arg1)
	return toPython(nil)
}

func sublime_MessageDialog(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for backend.DummyFrontend.MessageDialog() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	backend.GetEditor().Frontend().MessageDialog(arg1)
	return toPython(nil)
}

func sublime_OkCancelDialog(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
		arg2 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for backend.DummyFrontend.OkCancelDialog() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	if v, err := tu.GetItem(1); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for backend.DummyFrontend.OkCancelDialog() arg2, not %s", v.Type())
			} else {
				arg2 = v2
			}
		}
	}
	ret0 := backend.GetEditor().Frontend().OkCancelDialog(arg1, arg2)
	var err error
	var pyret0 py.Object

	pyret0, err = toPython(ret0)
	if err != nil {
		return nil, err
	}
	return pyret0, err
}

func sublime_StatusMessage(tu *py.Tuple) (py.Object, error) {
	var (
		arg1 string
	)
	if v, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		if v3, err2 := fromPython(v); err2 != nil {
			return nil, err2
		} else {
			if v2, ok := v3.(string); !ok {
				return nil, fmt.Errorf("Expected type string for backend.DummyFrontend.StatusMessage() arg1, not %s", v.Type())
			} else {
				arg1 = v2
			}
		}
	}
	backend.GetEditor().Frontend().StatusMessage(arg1)
	return toPython(nil)
}

func sublime_Console(tu *py.Tuple, kwargs *py.Dict) (py.Object, error) {
	if tu.Size() != 1 {
		return nil, fmt.Errorf("Unexpected argument count: %d", tu.Size())
	}
	if i, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		log.Info("Python sez: %s", i)
	}
	return toPython(nil)
}

func sublime_SetTimeOut(tu *py.Tuple, kwargs *py.Dict) (py.Object, error) {
	var (
		pyarg py.Object
	)
	if tu.Size() != 2 {
		return nil, fmt.Errorf("Unexpected argument count: %d", tu.Size())
	}
	if i, err := tu.GetItem(0); err != nil {
		return nil, err
	} else {
		pyarg = i
	}
	if i, err := tu.GetItem(1); err != nil {
		return nil, err
	} else if v, err := fromPython(i); err != nil {
		return nil, err
	} else if v2, ok := v.(int); !ok {
		return nil, fmt.Errorf("Expected int not %s", i.Type())
	} else {
		pyarg.Incref()
		go func() {
			time.Sleep(time.Millisecond * time.Duration(v2))
			l := py.NewLock()
			defer l.Unlock()
			defer pyarg.Decref()
			if ret, err := pyarg.Base().CallFunctionObjArgs(); err != nil {
				log.Debug("Error in callback: %v", err)
			} else {
				ret.Decref()
			}
		}()
	}
	return toPython(nil)
}

var manual_methods = []py.Method{
	{Name: "console", Func: sublime_Console},
	{Name: "set_timeout", Func: sublime_SetTimeOut},
	{Name: "packages_path", Func: sublime_PackagesPath},
	{Name: "error_message", Func: sublime_ErrorMessage},
	{Name: "message_dialog", Func: sublime_MessageDialog},
	{Name: "ok_cancel_dialog", Func: sublime_OkCancelDialog},
	{Name: "status_message", Func: sublime_StatusMessage},
}
