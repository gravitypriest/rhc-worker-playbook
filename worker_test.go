package main

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/subpop/go-log"
)

// seed uuid.New function for deterministic tests
var seededUuidString string = "080027c2-7382-b2cc-1967-000000000001"
var seededUuid uuid.UUID = uuid.MustParse(seededUuidString)

func mockUuid() uuid.UUID {
	return seededUuid
}

func readFile(t *testing.T, file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("cannot read file %v: %v", file, err)
	}
	return data
}

func TestVerifyPlaybook(t *testing.T) {
	_, err := exec.LookPath("/usr/libexec/rhc-playbook-verifier")
	if err != nil {
		t.Skip("rhc-playbook-verifier is not installed")
	}

	tests := []struct {
		description string
		input       struct {
			playbook []byte
		}
		want []byte
	}{
		{
			description: "insights_remove.yml",
			input: struct {
				playbook []byte
			}{
				playbook: readFile(t, "./testdata/insights_remove.yml"),
			},
			want: []byte(`- name: Insights Disable
  hosts: localhost
  become: yes
  vars:
    insights_signature_exclude: /hosts,/vars/insights_signature
  tasks:
  - name: Disable the insights-client
    command: insights-client --disable-schedule
`),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			log.SetLevel(log.LevelDebug)
			got, err := verifyPlaybook(test.input.playbook)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(got, test.want) {
				t.Errorf("\ngot:\n%v\nwant:\n%v", string(got), string(test.want))
			}
		})
	}
}

func TestGenerateExecutorOnFailedEvent(t *testing.T) {

	expectedFailureEvent := map[string]any{
		"event":      "executor_on_failed",
		"uuid":       seededUuidString,
		"counter":    -1,
		"start_line": 0,
		"end_line":   0,
		"event_data": map[string]any{
			"crc_dispatcher_correlation_id": "dcdc7b28-6800-4af9-983a-60fda58a7156",
			"crc_dispatcher_error_code":     "TEST_ERROR",
			"crc_dispatcher_error_details":  "playbook run failed",
		},
	}
	receivedFailureEvent := generateExecutorOnFailedEvent(
		"dcdc7b28-6800-4af9-983a-60fda58a7156",
		"TEST_ERROR",
		errors.New("playbook run failed"),
		mockUuid,
	)

	if !reflect.DeepEqual(expectedFailureEvent, receivedFailureEvent) {
		t.Errorf(
			"EXPECTED: %v\nRECEIVED: %v",
			expectedFailureEvent,
			receivedFailureEvent,
		)
	}
}

func TestGenerateExecutorOnStartEvent(t *testing.T) {

	expectedStartEvent := map[string]any{
		"event":      "executor_on_start",
		"uuid":       seededUuidString,
		"counter":    -1,
		"stdout":     "",
		"start_line": 0,
		"end_line":   0,
		"event_data": map[string]any{
			"crc_dispatcher_correlation_id": "dcdc7b28-6800-4af9-983a-60fda58a7156",
		},
	}
	receivedStartEvent := generateExecutorOnStartEvent(
		"dcdc7b28-6800-4af9-983a-60fda58a7156",
		mockUuid,
	)

	if !reflect.DeepEqual(expectedStartEvent, receivedStartEvent) {
		t.Errorf(
			"EXPECTED: %v\nRECEIVED: %v",
			expectedStartEvent,
			receivedStartEvent,
		)
	}
}
