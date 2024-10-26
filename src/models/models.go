package models

type Step struct {
	Name    string
	Command string
}

type Stage struct {
	Name  string
	Steps []Step
}

type Pipeline struct {
	Stages []Stage
}
