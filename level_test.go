package l4g

import (
	"encoding/json"
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelTrace, "trace"},
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
		{LevelPanic, "panic"},
		{LevelFatal, "fatal"},
		{LevelTrace - 1, "trace"},
		{LevelFatal + 1, "fatal"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_Int(t *testing.T) {
	tests := []struct {
		level Level
		want  int
	}{
		{LevelTrace, 1},
		{LevelDebug, 2},
		{LevelInfo, 3},
		{LevelWarn, 4},
		{LevelError, 5},
		{LevelPanic, 6},
		{LevelFatal, 7},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := tt.level.Int(); got != tt.want {
				t.Errorf("Level.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_Real(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		want  Level
	}{
		{"normal trace", LevelTrace, LevelTrace},
		{"normal info", LevelInfo, LevelInfo},
		{"normal fatal", LevelFatal, LevelFatal},
		{"below trace", LevelTrace - 10, LevelTrace},
		{"above fatal", LevelFatal + 10, LevelFatal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Real(); got != tt.want {
				t.Errorf("Level.Real() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_MarshalJSON(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelTrace, `"trace"`},
		{LevelDebug, `"debug"`},
		{LevelInfo, `"info"`},
		{LevelWarn, `"warn"`},
		{LevelError, `"error"`},
		{LevelPanic, `"panic"`},
		{LevelFatal, `"fatal"`},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			got, err := json.Marshal(tt.level)
			if err != nil {
				t.Errorf("Level.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Level.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Level
		wantErr bool
	}{
		{"trace", `"trace"`, LevelTrace, false},
		{"debug", `"debug"`, LevelDebug, false},
		{"info", `"info"`, LevelInfo, false},
		{"warn", `"warn"`, LevelWarn, false},
		{"error", `"error"`, LevelError, false},
		{"panic", `"panic"`, LevelPanic, false},
		{"fatal", `"fatal"`, LevelFatal, false},
		{"uppercase", `"INFO"`, LevelInfo, false},
		{"mixed case", `"DeBuG"`, LevelDebug, false},
		{"unknown", `"unknown"`, LevelTrace, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Level
			err := json.Unmarshal([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Level.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Level.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_MarshalText(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelTrace, "trace"},
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
		{LevelPanic, "panic"},
		{LevelFatal, "fatal"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			got, err := tt.level.MarshalText()
			if err != nil {
				t.Errorf("Level.MarshalText() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Level.MarshalText() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestLevel_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    Level
		wantErr bool
	}{
		{"trace", "trace", LevelTrace, false},
		{"debug", "debug", LevelDebug, false},
		{"info", "info", LevelInfo, false},
		{"warn", "warn", LevelWarn, false},
		{"error", "error", LevelError, false},
		{"panic", "panic", LevelPanic, false},
		{"fatal", "fatal", LevelFatal, false},
		{"uppercase", "WARN", LevelWarn, false},
		{"mixed case", "ErRoR", LevelError, false},
		{"unknown", "unknown", LevelTrace, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Level
			err := got.UnmarshalText([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("Level.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Level.UnmarshalText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_Level(t *testing.T) {
	tests := []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal}
	for _, level := range tests {
		t.Run(level.String(), func(t *testing.T) {
			if got := level.Level(); got != level {
				t.Errorf("Level.Level() = %v, want %v", got, level)
			}
		})
	}
}

func TestLevelVar_NewLevelVar(t *testing.T) {
	t.Run("from level", func(t *testing.T) {
		lv := NewLevelVar(LevelWarn)
		if lv.Level() != LevelWarn {
			t.Errorf("NewLevelVar() level = %v, want %v", lv.Level(), LevelWarn)
		}
	})

	t.Run("from levelvar", func(t *testing.T) {
		lv1 := NewLevelVar(LevelError)
		lv2 := NewLevelVar(lv1)
		if lv1 != lv2 {
			t.Errorf("NewLevelVar() should return same instance")
		}
	})
}

func TestLevelVar_SetAndGet(t *testing.T) {
	lv := &LevelVar{}

	tests := []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal}
	for _, want := range tests {
		t.Run(want.String(), func(t *testing.T) {
			lv.Set(want)
			if got := lv.Level(); got != want {
				t.Errorf("LevelVar.Level() = %v, want %v", got, want)
			}
			if gotInt := lv.Int(); gotInt != want.Int() {
				t.Errorf("LevelVar.Int() = %v, want %v", gotInt, want.Int())
			}
		})
	}
}

func TestLevelVar_String(t *testing.T) {
	lv := NewLevelVar(LevelInfo)
	want := "LevelVar(info)"
	if got := lv.String(); got != want {
		t.Errorf("LevelVar.String() = %v, want %v", got, want)
	}
}

func TestLevelVar_MarshalText(t *testing.T) {
	lv := NewLevelVar(LevelWarn)
	got, err := lv.MarshalText()
	if err != nil {
		t.Errorf("LevelVar.MarshalText() error = %v", err)
		return
	}
	want := "warn"
	if string(got) != want {
		t.Errorf("LevelVar.MarshalText() = %v, want %v", string(got), want)
	}
}

func TestLevelVar_UnmarshalText(t *testing.T) {
	lv := &LevelVar{}
	err := lv.UnmarshalText([]byte("error"))
	if err != nil {
		t.Errorf("LevelVar.UnmarshalText() error = %v", err)
		return
	}
	if lv.Level() != LevelError {
		t.Errorf("LevelVar.Level() = %v, want %v", lv.Level(), LevelError)
	}
}

func TestLevelVar_Concurrent(t *testing.T) {
	lv := NewLevelVar(LevelInfo)

	done := make(chan bool)
	for i := range 10 {
		go func(level Level) {
			lv.Set(level)
			_ = lv.Level()
			_ = lv.Int()
			done <- true
		}(Level(i % 7))
	}

	for range 10 {
		<-done
	}
}
