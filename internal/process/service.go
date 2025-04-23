package process

type Action string

const (
	ActionNotify  Action = "notify"
	ActionUpgrade        = "upgrade"
)

type Options struct {
	Enable bool   `json:"autodock.enable"`
	Action Action `json:"autodock.action"`
	// TODO: Add extra dependencies to update before/stop after
}

func OptsFromLabels(labels map[string]string) Options {
	opts := &Options{
		Enable: false,
		Action: ActionNotify,
	}

	if enableStr, ok := labels["autodock.enable"]; ok {
		opts.Enable = enableStr == "true"
	}

	if actionStr, ok := labels["autodock.action"]; ok {
		opts.Action = Action(actionStr)
	}

	return *opts
}
