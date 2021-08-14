package shell

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyConversionToArgs(t *testing.T) {
	args := NewArgs().GetAll()
	assert.Equal(t, []string{}, args)
}

func TestConversionToArgsFromFlags(t *testing.T) {
	args := NewArgs()
	args.AddFlags("aaa", []string{"one", "two"}, ArgConfigEscape)
	args.AddFlag("bbb", "three", ArgConfigEscape)
	assert.Equal(t, []string{"--aaa", "one", "--aaa", "two", "--bbb", "three"}, args.GetAll())
}

func TestConversionToArgsNoFlag(t *testing.T) {
	args := NewArgs()
	args.AddArgs([]string{"one", "two"}, ArgConfigEscape)
	args.AddArg("three", ArgConfigEscape)
	assert.Equal(t, []string{"one", "two", "three"}, args.GetAll())
}

func TestConversionToArgs(t *testing.T) {
	args := NewArgs()
	args.AddFlags("aaa", []string{"simple", "with space", "with\"quote"}, ArgConfigEscape)
	args.AddFlags("bbb", []string{"simple", "with space", "with\"quote"}, ArgConfigKeepGlobQuote)
	args.AddArgs([]string{"with space", "with\"quote", "with$variable"}, ArgConfigEscape)
	args.AddArg("with space\"quote", ArgConfigKeepGlobQuote)
	args.AddArg("with$variable", ArgConfigKeepGlobQuote)

	expected := []string{
		"--aaa",
		"simple",
		"--aaa",
		`with\ space`,
		"--aaa",
		`with\"quote`,
		"--bbb",
		"simple",
		"--bbb",
		`"with space"`,
		"--bbb",
		`"with\"quote"`,
		`with\ space`,
		`with\"quote`,
		"with\\$variable",
		`"with space\"quote"`,
		"\"with$variable\"",
	}
	if runtime.GOOS == "windows" {
		expected = []string{
			"--aaa",
			"simple",
			"--aaa",
			"with space",
			"--aaa",
			"with\"quote",
			"--bbb",
			"simple",
			"--bbb",
			"with space",
			"--bbb",
			"with\"quote",
			"with space",
			"with\"quote",
			"with$variable",
			"with space\"quote",
			"with$variable",
		}
	}
	assert.Equal(t, expected, args.GetAll())
}