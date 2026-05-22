package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	enabled bool
	debug   bool
	out     io.Writer
}

func New() *Logger {
	return &Logger{
		enabled: false,
		debug:   false,
		out:     os.Stdout,
	}
}

func (l *Logger) Enable() {
	l.enabled = true
}

func (l *Logger) Disable() {
	l.enabled = false
}

func (l *Logger) EnableDebug() {
	l.debug = true
	l.enabled = true
}

func (l *Logger) DisableDebug() {
	l.debug = false
}

func (l *Logger) SetOutput(out io.Writer) {
	if out == nil {
		return
	}

	l.out = out
}

func (l *Logger) Spell(format string, args ...any) {
	if !l.enabled {
		return
	}

	_, err := fmt.Fprintf(l.out, format+"\n", args...)
	if err != nil {
		return
	}
}

func (l *Logger) DebugSpell(format string, args ...any) {
	if !l.debug {
		return
	}

	_, err := fmt.Fprintf(l.out, format+"\n", args...)
	if err != nil {
		return
	}
}

func (l *Logger) ForbiddenSQL(sqlText string, args ...any) {
	l.DebugSpell("[📜] Forbidden SQL revealed:")
	l.DebugSpell("%s", sqlText)

	if len(args) > 0 {
		l.DebugSpell("[🎒] With offerings: %v", args)
	}
}

func (l *Logger) Summoning(modelName string) {
	l.Spell("[🧙] Summoning %s...", modelName)
}

func (l *Logger) ManaAccepted() {
	l.Spell("[🔥] Mana accepted.")
}

func (l *Logger) InsertSuccessful() {
	l.Spell("[☠️] INSERT successful.")
}

func (l *Logger) Raising(modelName string) {
	l.Spell("[🧟] Raising table for %s...", modelName)
}

func (l *Logger) MigrationSuccessful() {
	l.Spell("[🪦] Table awakened successfully.")
}

func (l *Logger) Divining(modelName string) {
	l.Spell("[🔮] Divining %s from the database...", modelName)
}

func (l *Logger) QuerySuccessful() {
	l.Spell("[✨] Divination successful.")
}

func (l *Logger) Sacrificing(modelName string) {
	l.Spell("[🩸] Sacrificing %s...", modelName)
}

func (l *Logger) DeleteSuccessful() {
	l.Spell("[☠️] DELETE successful.")
}

func (l *Logger) CastingHook(hookName string) {
	l.Spell("[🪄] Casting %s hook...", hookName)
}

func (l *Logger) HookSuccessful(hookName string) {
	l.Spell("[✨] %s hook completed.", hookName)
}

func (l *Logger) SpellFailed(err error) {
	if err == nil {
		return
	}

	l.Spell("[💀] Spell failed: %v", err)
}
