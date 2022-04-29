package patchers

type PatcherToolDarwin struct {
	PatcherTool
}

func (p PatcherToolDarwin) FindDirectories() ([]string, []string) {
	panic("not implemented")
}
