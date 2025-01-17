package cron

type Jobs struct {
	FailCron *FailCron
}

func (jobs *Jobs) StartAll() {
	jobs.FailCron.Start()
}
