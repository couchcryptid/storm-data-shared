package config

import (
	"testing"
	"time"
)

func TestEnvOrDefault_EnvSet(t *testing.T) {
	t.Setenv("TEST_KEY", "fromenv")
	if got := EnvOrDefault("TEST_KEY", "fallback"); got != "fromenv" {
		t.Errorf("got %q, want %q", got, "fromenv")
	}
}

func TestEnvOrDefault_Fallback(t *testing.T) {
	if got := EnvOrDefault("UNSET_KEY_12345", "fallback"); got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func TestParseBrokers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"single", "localhost:9092", []string{"localhost:9092"}},
		{"multiple", "a:9092,b:9092", []string{"a:9092", "b:9092"}},
		{"whitespace", " a:9092 , b:9092 ", []string{"a:9092", "b:9092"}},
		{"empty parts", "a:9092,,b:9092", []string{"a:9092", "b:9092"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseBrokers(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseBatchSize_Default(t *testing.T) {
	n, err := ParseBatchSize()
	if err != nil {
		t.Fatal(err)
	}
	if n != 50 {
		t.Errorf("got %d, want 50", n)
	}
}

func TestParseBatchSize_Custom(t *testing.T) {
	t.Setenv("BATCH_SIZE", "100")
	n, err := ParseBatchSize()
	if err != nil {
		t.Fatal(err)
	}
	if n != 100 {
		t.Errorf("got %d, want 100", n)
	}
}

func TestParseBatchSize_Invalid(t *testing.T) {
	t.Setenv("BATCH_SIZE", "0")
	_, err := ParseBatchSize()
	if err == nil {
		t.Error("expected error for BATCH_SIZE=0")
	}
}

func TestParseBatchSize_TooLarge(t *testing.T) {
	t.Setenv("BATCH_SIZE", "1001")
	_, err := ParseBatchSize()
	if err == nil {
		t.Error("expected error for BATCH_SIZE=1001")
	}
}

func TestParseBatchFlushInterval_Default(t *testing.T) {
	d, err := ParseBatchFlushInterval()
	if err != nil {
		t.Fatal(err)
	}
	if d != 500*time.Millisecond {
		t.Errorf("got %v, want 500ms", d)
	}
}

func TestParseShutdownTimeout_Default(t *testing.T) {
	d, err := ParseShutdownTimeout()
	if err != nil {
		t.Fatal(err)
	}
	if d != 10*time.Second {
		t.Errorf("got %v, want 10s", d)
	}
}

func TestParseShutdownTimeout_Invalid(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "-1s")
	_, err := ParseShutdownTimeout()
	if err == nil {
		t.Error("expected error for negative timeout")
	}
}
