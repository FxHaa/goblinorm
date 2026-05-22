package hooks

import (
	"fmt"

	"goblinorm/internal/goblinerr"
	"goblinorm/internal/logger"
)

type BeforeSummoner interface {
	BeforeSummon() error
}

type AfterSummoner interface {
	AfterSummon() error
}

type BeforeDiviner interface {
	BeforeDivine() error
}

type AfterDiviner interface {
	AfterDivine() error
}

type BeforeSacrificer interface {
	BeforeSacrifice() error
}

type AfterSacrificer interface {
	AfterSacrifice() error
}

func BeforeSummon(model any, log *logger.Logger) error {
	hook, ok := model.(BeforeSummoner)
	if !ok {
		return nil
	}

	return cast("BeforeSummon", log, hook.BeforeSummon)
}

func AfterSummon(model any, log *logger.Logger) error {
	hook, ok := model.(AfterSummoner)
	if !ok {
		return nil
	}

	return cast("AfterSummon", log, hook.AfterSummon)
}

func BeforeDivine(model any, log *logger.Logger) error {
	hook, ok := model.(BeforeDiviner)
	if !ok {
		return nil
	}

	return cast("BeforeDivine", log, hook.BeforeDivine)
}

func AfterDivine(model any, log *logger.Logger) error {
	hook, ok := model.(AfterDiviner)
	if !ok {
		return nil
	}

	return cast("AfterDivine", log, hook.AfterDivine)
}

func BeforeSacrifice(model any, log *logger.Logger) error {
	hook, ok := model.(BeforeSacrificer)
	if !ok {
		return nil
	}

	return cast("BeforeSacrifice", log, hook.BeforeSacrifice)
}

func AfterSacrifice(model any, log *logger.Logger) error {
	hook, ok := model.(AfterSacrificer)
	if !ok {
		return nil
	}

	return cast("AfterSacrifice", log, hook.AfterSacrifice)
}

func cast(hookName string, log *logger.Logger, hook func() error) error {
	if log != nil {
		log.CastingHook(hookName)
	}

	if err := hook(); err != nil {
		err = fmt.Errorf("%w: %s: %v", goblinerr.ErrHookFailed, hookName, err)
		if log != nil {
			log.SpellFailed(err)
		}
		return err
	}

	if log != nil {
		log.HookSuccessful(hookName)
	}

	return nil
}
