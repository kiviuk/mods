package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindCacheOpsDetails(t *testing.T) {
	newMods := func(t *testing.T) *Mods {
		db := testDB(t)
		return &Mods{
			db:     db,
			Config: &Config{},
		}
	}

	t.Run("all empty", func(t *testing.T) {
		msg := newMods(t).findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Empty(t, dets.ReadID)
		require.NotEmpty(t, dets.WriteID)
		require.Empty(t, dets.Title)
	})

	t.Run("show id", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message", "openai", "gpt-4"))
		mods.Config.Show = id[:8]
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
	})

	t.Run("show title", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.Show = "message 1"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
	})

	t.Run("continue id", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message", "openai", "gpt-4"))
		mods.Config.Continue = id[:5]
		mods.Config.Prefix = "prompt"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.Equal(t, id, dets.WriteID)
	})

	t.Run("continue with no prompt", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.ContinueLast = true
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.Equal(t, id, dets.WriteID)
		require.Empty(t, dets.Title)
	})

	t.Run("continue title", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.Continue = "message 1"
		mods.Config.Prefix = "prompt"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.Equal(t, id, dets.WriteID)
	})

	t.Run("continue last", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.ContinueLast = true
		mods.Config.Prefix = "prompt"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.Equal(t, id, dets.WriteID)
		require.Empty(t, dets.Title)
	})

	t.Run("continue last with name", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.Continue = "message 2"
		mods.Config.Prefix = "prompt"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.Equal(t, "message 2", dets.Title)
		require.NotEmpty(t, dets.WriteID)
		require.Equal(t, id, dets.WriteID)
	})

	t.Run("write", func(t *testing.T) {
		mods := newMods(t)
		mods.Config.Title = "some title"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Empty(t, dets.ReadID)
		require.NotEmpty(t, dets.WriteID)
		require.NotEqual(t, "some title", dets.WriteID)
		require.Equal(t, "some title", dets.Title)
	})

	t.Run("continue id and write with title", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.Title = "some title"
		mods.Config.Continue = id[:10]
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.NotEmpty(t, dets.WriteID)
		require.NotEqual(t, id, dets.WriteID)
		require.NotEqual(t, "some title", dets.WriteID)
		require.Equal(t, "some title", dets.Title)
	})

	t.Run("continue title and write with title", func(t *testing.T) {
		mods := newMods(t)
		id := newConversationID()
		require.NoError(t, mods.db.Save(id, "message 1", "openai", "gpt-4"))
		mods.Config.Title = "some title"
		mods.Config.Continue = "message 1"
		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)
		require.Equal(t, id, dets.ReadID)
		require.NotEmpty(t, dets.WriteID)
		require.NotEqual(t, id, dets.WriteID)
		require.NotEqual(t, "some title", dets.WriteID)
		require.Equal(t, "some title", dets.Title)
	})

	t.Run("show invalid", func(t *testing.T) {
		mods := newMods(t)
		mods.Config.Show = "aaa"
		msg := mods.findCacheOpsDetails()()
		err := msg.(modsError)
		require.Equal(t, "Could not find the conversation.", err.reason)
		require.EqualError(t, err, "no conversations found: aaa")
	})

	t.Run("uses config model and api not global config", func(t *testing.T) {
		mods := newMods(t)
		mods.Config.Model = "claude-3.7-sonnet"
		mods.Config.API = "anthropic"

		msg := mods.findCacheOpsDetails()()
		dets := msg.(cacheDetailsMsg)

		require.Equal(t, "claude-3.7-sonnet", dets.Model)
		require.Equal(t, "anthropic", dets.API)
		require.Empty(t, dets.ReadID)
		require.NotEmpty(t, dets.WriteID)
	})
}

func TestRemoveWhitespace(t *testing.T) {
	t.Run("only whitespaces", func(t *testing.T) {
		require.Equal(t, "", removeWhitespace(" \n"))
	})

	t.Run("regular text", func(t *testing.T) {
		require.Equal(t, " regular\n ", removeWhitespace(" regular\n "))
	})
}

var cutPromptTests = map[string]struct {
	msg      string
	prompt   string
	expected string
}{
	"bad error": {
		msg:      "nope",
		prompt:   "the prompt",
		expected: "the prompt",
	},
	"crazy error": {
		msg:      tokenErrMsg(10, 93),
		prompt:   "the prompt",
		expected: "the prompt",
	},
	"cut prompt": {
		msg:      tokenErrMsg(10, 3),
		prompt:   "this is a long prompt I have no idea if its really 10 tokens",
		expected: "this is a long prompt ",
	},
	"missmatch of token estimation vs api result": {
		msg:      tokenErrMsg(30000, 100),
		prompt:   "tell me a joke",
		expected: "tell me a joke",
	},
}

func tokenErrMsg(l, ml int) string {
	return fmt.Sprintf(`This model's maximum context length is %d tokens. However, your messages resulted in %d tokens`, ml, l)
}

func TestCutPrompt(t *testing.T) {
	for name, tc := range cutPromptTests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, cutPrompt(tc.msg, tc.prompt))
		})
	}
}

func TestExpandStatusText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		config   *Config
		expected string
	}{
		{
			name:     "nil config",
			text:     "Generating with {model_alias_name}",
			config:   nil,
			expected: "Generating with {model_alias_name}",
		},
		{
			name:     "replace model_alias_name",
			text:     "Generating with {model_alias_name}",
			config:   &Config{Model: "4o"},
			expected: "Generating with 4o",
		},
		{
			name:     "replace model_full_name",
			text:     "Generating with {model_full_name}",
			config:   &Config{ModelFullName: "gpt-4o"},
			expected: "Generating with gpt-4o",
		},
		{
			name:     "replace all placeholders",
			text:     "{model_alias_name}/{model_full_name} at {temp} as {role}",
			config:   &Config{Model: "4o", ModelFullName: "gpt-4o", Temperature: 0.7, Role: "assistant"},
			expected: "4o/gpt-4o at 0.7 as assistant",
		},
		{
			name:     "zero temperature",
			text:     "{model_alias_name} at {temp}",
			config:   &Config{Model: "gpt-4", Temperature: 0},
			expected: "gpt-4 at ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, expandStatusText(tc.text, tc.config))
		})
	}
}

func TestResolveModelAlias(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		wantAlias    string
		wantFullName string
	}{
		{
			name:         "exact match",
			config:       &Config{Model: "gpt-4o", API: "openai", APIs: APIs{{Name: "openai", Models: map[string]Model{"gpt-4o": {Aliases: []string{"4o"}}}}}},
			wantAlias:    "gpt-4o",
			wantFullName: "gpt-4o",
		},
		{
			name:         "alias match",
			config:       &Config{Model: "4o", API: "openai", APIs: APIs{{Name: "openai", Models: map[string]Model{"gpt-4o": {Aliases: []string{"4o"}}}}}},
			wantAlias:    "4o",
			wantFullName: "gpt-4o",
		},
		{
			name:         "no match",
			config:       &Config{Model: "unknown", API: "openai", APIs: APIs{{Name: "openai", Models: map[string]Model{"gpt-4o": {Aliases: []string{"4o"}}}}}},
			wantAlias:    "unknown",
			wantFullName: "unknown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			alias, fullName := resolveModelAlias(tc.config)
			require.Equal(t, tc.wantAlias, alias)
			require.Equal(t, tc.wantFullName, fullName)
		})
	}
}
