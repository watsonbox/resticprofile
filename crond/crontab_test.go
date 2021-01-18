//+build !darwin,!windows

package crond

import (
	"strings"
	"testing"

	"github.com/creativeprojects/resticprofile/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateEmptyCrontab(t *testing.T) {
	crontab := NewCrontab(nil)
	buffer := &strings.Builder{}
	err := crontab.Generate(buffer)
	require.NoError(t, err)
	assert.Equal(t, "", buffer.String())
}

func TestGenerateSimpleCrontab(t *testing.T) {
	crontab := NewCrontab([]Entry{NewEntry(calendar.NewEvent(func(event *calendar.Event) {
		event.Minute.MustAddValue(1)
		event.Hour.MustAddValue(1)
	}), "", "", "", "resticprofile backup")})
	buffer := &strings.Builder{}
	err := crontab.Generate(buffer)
	require.NoError(t, err)
	assert.Equal(t, "01 01 * * *\tresticprofile backup\n", buffer.String())
}

func TestCleanupCrontab(t *testing.T) {
	crontab := `# DO NOT EDIT THIS FILE - edit the master and reinstall.
# (/tmp/crontab.pMvuGY/crontab installed on Wed Jan 13 12:08:43 2021)
# (Cron version -- $Id: crontab.c,v 2.13 1994/01/17 03:20:37 vixie Exp $)
# m h  dom mon dow   command
`
	assert.Equal(t, "# m h  dom mon dow   command\n", cleanupCrontab(crontab))
}

func TestCleanCrontab(t *testing.T) {
	crontab := `#
#
#
# m h  dom mon dow   command
`
	assert.Equal(t, "#\n#\n#\n# m h  dom mon dow   command\n", cleanupCrontab(crontab))
}

func TestDeleteLine(t *testing.T) {
	testData := []struct {
		source   string
		expected bool
	}{
		{"#\n#\n#\n# 00,30 * * * *	/home/resticprofile --no-ansi --config config.yaml --name profile --log backup.log backup\n", false},
		{"#\n#\n#\n00,30 * * * *	/home/resticprofile --no-ansi --config config.yaml --name profile --log backup.log backup\n", true},
	}

	for _, testRun := range testData {
		t.Run("", func(t *testing.T) {
			_, found, err := deleteLine(testRun.source, Entry{configFile: "config.yaml", profileName: "profile", commandName: "backup"})
			require.NoError(t, err)
			assert.Equal(t, testRun.expected, found)
		})
	}
}

func TestVirginCrontab(t *testing.T) {
	crontab := "#\n#\n#\n# m h  dom mon dow   command\n"
	result, _, _, found := extractOwnSection(crontab)
	assert.False(t, found)
	assert.Equal(t, crontab, result)
}

func TestOwnSection(t *testing.T) {
	own := "-- 1\n#\n2\n3\n# --\n"
	before := "#\n#\n#\n# m h  dom mon dow   command\n"
	after := "# blah blah\n"
	crontab := before + startMarker + own + endMarker + after
	beforeResult, result, afterResult, found := extractOwnSection(crontab)
	assert.True(t, found)
	assert.Equal(t, own, result)
	assert.Equal(t, before, beforeResult)
	assert.Equal(t, after, afterResult)
}

func TestSectionOnItsOwn(t *testing.T) {
	own := "-- 1\n#\n2\n3\n# --\n"
	crontab := startMarker + own + endMarker
	beforeResult, result, afterResult, found := extractOwnSection(crontab)
	assert.True(t, found)
	assert.Equal(t, own, result)
	assert.Equal(t, "", beforeResult)
	assert.Equal(t, "", afterResult)
}

func TestUpdateEmptyCrontab(t *testing.T) {
	crontab := NewCrontab(nil)
	buffer := &strings.Builder{}
	err := crontab.Update("", true, buffer)
	require.NoError(t, err)
	assert.Equal(t, "\n"+startMarker+endMarker, buffer.String())
}

func TestUpdateSimpleCrontab(t *testing.T) {
	crontab := NewCrontab([]Entry{NewEntry(calendar.NewEvent(func(event *calendar.Event) {
		event.Minute.MustAddValue(1)
		event.Hour.MustAddValue(1)
	}), "", "", "", "resticprofile backup")})
	buffer := &strings.Builder{}
	err := crontab.Update("", true, buffer)
	require.NoError(t, err)
	assert.Equal(t, "\n"+startMarker+"01 01 * * *\tresticprofile backup\n"+endMarker, buffer.String())
}

func TestUpdateExistingCrontab(t *testing.T) {
	crontab := NewCrontab([]Entry{NewEntry(calendar.NewEvent(func(event *calendar.Event) {
		event.Minute.MustAddValue(1)
		event.Hour.MustAddValue(1)
	}), "", "", "", "resticprofile backup")})
	buffer := &strings.Builder{}
	err := crontab.Update("something\n"+startMarker+endMarker, true, buffer)
	require.NoError(t, err)
	assert.Equal(t, "something\n"+startMarker+"01 01 * * *\tresticprofile backup\n"+endMarker, buffer.String())
}
