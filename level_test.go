package l4g

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

func TestLevelString(t *testing.T) {
	for _, test := range []struct {
		in   Level
		want string
	}{
		{VERBOSE, "VERBOSE"},
		{INFO, "INFO"},
		{ERROR, "ERROR"},
		{ERROR + 2, "ERROR+2"},
		{ERROR - 2, "WARN+2"},
		{WARN, "WARN"},
		{WARN - 1, "INFO+3"},
		{INFO, "INFO"},
		{INFO + 1, "INFO+1"},
		{INFO - 3, "DEBUG+1"},
		{DEBUG, "DEBUG"},
		{DEBUG - 2, "DEBUG-2"},
	} {
		got := test.in.String()
		if got != test.want {
			t.Errorf("%d: got %s, want %s", test.in, got, test.want)
		}
	}
}

func TestLevelVar(t *testing.T) {
	var al LevelVar
	if got, want := al.Level(), INFO; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	al.Set(WARN)
	if got, want := al.Level(), WARN; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	al.Set(INFO)
	if got, want := al.Level(), INFO; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

}

func TestLevelMarshalJSON(t *testing.T) {
	want := WARN - 3
	wantData := []byte(`"INFO+1"`)
	data, err := want.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, wantData) {
		t.Errorf("got %s, want %s", string(data), string(wantData))
	}
	var got Level
	if err := got.UnmarshalJSON(data); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestLevelMarshalText(t *testing.T) {
	want := WARN - 3
	wantData := []byte("INFO+1")
	data, err := want.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, wantData) {
		t.Errorf("got %s, want %s", string(data), string(wantData))
	}
	var got Level
	if err := got.UnmarshalText(data); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestLevelAppendText(t *testing.T) {
	buf := make([]byte, 4, 16)
	want := WARN - 3
	wantData := []byte("\x00\x00\x00\x00INFO+1")
	data, err := want.AppendText(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, wantData) {
		t.Errorf("got %s, want %s", string(data), string(wantData))
	}
}

func TestLevelParse(t *testing.T) {
	for _, test := range []struct {
		in   string
		want Level
	}{
		{"DEBUG", DEBUG},
		{"INFO", INFO},
		{"WARN", WARN},
		{"ERROR", ERROR},
		{"debug", DEBUG},
		{"iNfo", INFO},
		{"INFO+87", INFO + 87},
		{"Error-18", ERROR - 18},
		{"Error-8", INFO},
	} {
		var got Level
		if err := got.parse(test.in); err != nil {
			t.Fatalf("%q: %v", test.in, err)
		}
		if got != test.want {
			t.Errorf("%q: got %s, want %s", test.in, got, test.want)
		}
	}
}

func TestLevelParseError(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string // error string should contain this
	}{
		{"", "unknown name"},
		{"dbg", "unknown name"},
		{"INFO+", "invalid syntax"},
		{"INFO-", "invalid syntax"},
		{"ERROR+23x", "invalid syntax"},
	} {
		var l Level
		err := l.parse(test.in)
		if err == nil || !strings.Contains(err.Error(), test.want) {
			t.Errorf("%q: got %v, want string containing %q", test.in, err, test.want)
		}
	}
}

func TestLevelFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	lf := INFO
	fs.TextVar(&lf, "level", lf, "set level")
	err := fs.Parse([]string{"-level", "WARN+3"})
	if err != nil {
		t.Fatal(err)
	}
	if g, w := lf, WARN+3; g != w {
		t.Errorf("got %v, want %v", g, w)
	}
}

func TestLevelVarMarshalText(t *testing.T) {
	var v LevelVar
	v.Set(WARN)
	data, err := v.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	var v2 LevelVar
	if err := v2.UnmarshalText(data); err != nil {
		t.Fatal(err)
	}
	if g, w := v2.Level(), WARN; g != w {
		t.Errorf("got %s, want %s", g, w)
	}
}

func TestLevelVarAppendText(t *testing.T) {
	var v LevelVar
	v.Set(WARN)
	buf := make([]byte, 4, 16)
	data, err := v.AppendText(buf)
	if err != nil {
		t.Fatal(err)
	}
	var v2 LevelVar
	if err := v2.UnmarshalText(data[4:]); err != nil {
		t.Fatal(err)
	}
	if g, w := v2.Level(), WARN; g != w {
		t.Errorf("got %s, want %s", g, w)
	}
}

func TestLevelVarFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	v := &LevelVar{}
	v.Set(WARN + 3)
	fs.TextVar(v, "level", v, "set level")
	err := fs.Parse([]string{"-level", "WARN+3"})
	if err != nil {
		t.Fatal(err)
	}
	if g, w := v.Level(), WARN+3; g != w {
		t.Errorf("got %v, want %v", g, w)
	}
}

func TestLevelVarString(t *testing.T) {
	var v LevelVar
	v.Set(ERROR)
	got := v.String()
	want := "LevelVar(ERROR)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
