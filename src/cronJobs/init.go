package cronJobs

type Jobs struct {
	FailCron *FailCronImpl
}

func (jobs *Jobs) StartAll() {
	jobs.FailCron.Start()
}
