package gomodule

import "github.com/google/blueprint"

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		// TODO: Визначте поля структури, щоб отримати дані з визначень у файлі build.bood
	}
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	// TODO: Імплементууйте генерацію правил збірки для ninja-файла.
}
