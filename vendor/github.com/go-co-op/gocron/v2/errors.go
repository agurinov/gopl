package gocron

import "fmt"

// Public error definitions
var (
	ErrCronJobInvalid                = fmt.Errorf("gocron: CronJob: invalid crontab")
	ErrCronJobParse                  = fmt.Errorf("gocron: CronJob: crontab parse failure")
	ErrDailyJobAtTimeNil             = fmt.Errorf("gocron: DailyJob: atTime within atTimes must not be nil")
	ErrDailyJobAtTimesNil            = fmt.Errorf("gocron: DailyJob: atTimes must not be nil")
	ErrDailyJobHours                 = fmt.Errorf("gocron: DailyJob: atTimes hours must be between 0 and 23 inclusive")
	ErrDailyJobMinutesSeconds        = fmt.Errorf("gocron: DailyJob: atTimes minutes and seconds must be between 0 and 59 inclusive")
	ErrDurationJobIntervalZero       = fmt.Errorf("gocron: DurationJob: time interval is 0")
	ErrDurationRandomJobMinMax       = fmt.Errorf("gocron: DurationRandomJob: minimum duration must be less than maximum duration")
	ErrEventListenerFuncNil          = fmt.Errorf("gocron: eventListenerFunc must not be nil")
	ErrJobNotFound                   = fmt.Errorf("gocron: job not found")
	ErrJobRunNowFailed               = fmt.Errorf("gocron: Job: RunNow: scheduler unreachable")
	ErrMonthlyJobDays                = fmt.Errorf("gocron: MonthlyJob: daysOfTheMonth must be between 31 and -31 inclusive, and not 0")
	ErrMonthlyJobAtTimeNil           = fmt.Errorf("gocron: MonthlyJob: atTime within atTimes must not be nil")
	ErrMonthlyJobAtTimesNil          = fmt.Errorf("gocron: MonthlyJob: atTimes must not be nil")
	ErrMonthlyJobDaysNil             = fmt.Errorf("gocron: MonthlyJob: daysOfTheMonth must not be nil")
	ErrMonthlyJobHours               = fmt.Errorf("gocron: MonthlyJob: atTimes hours must be between 0 and 23 inclusive")
	ErrMonthlyJobMinutesSeconds      = fmt.Errorf("gocron: MonthlyJob: atTimes minutes and seconds must be between 0 and 59 inclusive")
	ErrNewJobTaskNil                 = fmt.Errorf("gocron: NewJob: Task must not be nil")
	ErrNewJobTaskNotFunc             = fmt.Errorf("gocron: NewJob: Task.Function must be of kind reflect.Func")
	ErrNewJobWrongNumberOfParameters = fmt.Errorf("gocron: NewJob: Number of provided parameters does not match expected")
	ErrNewJobWrongTypeOfParameters   = fmt.Errorf("gocron: NewJob: Type of provided parameters does not match expected")
	ErrOneTimeJobStartDateTimePast   = fmt.Errorf("gocron: OneTimeJob: start must not be in the past")
	ErrStopExecutorTimedOut          = fmt.Errorf("gocron: timed out waiting for executor to stop")
	ErrStopJobsTimedOut              = fmt.Errorf("gocron: timed out waiting for jobs to finish")
	ErrStopSchedulerTimedOut         = fmt.Errorf("gocron: timed out waiting for scheduler to stop")
	ErrWeeklyJobAtTimeNil            = fmt.Errorf("gocron: WeeklyJob: atTime within atTimes must not be nil")
	ErrWeeklyJobAtTimesNil           = fmt.Errorf("gocron: WeeklyJob: atTimes must not be nil")
	ErrWeeklyJobDaysOfTheWeekNil     = fmt.Errorf("gocron: WeeklyJob: daysOfTheWeek must not be nil")
	ErrWeeklyJobHours                = fmt.Errorf("gocron: WeeklyJob: atTimes hours must be between 0 and 23 inclusive")
	ErrWeeklyJobMinutesSeconds       = fmt.Errorf("gocron: WeeklyJob: atTimes minutes and seconds must be between 0 and 59 inclusive")
	ErrPanicRecovered                = fmt.Errorf("gocron: panic recovered")
	ErrWithClockNil                  = fmt.Errorf("gocron: WithClock: clock must not be nil")
	ErrWithDistributedElectorNil     = fmt.Errorf("gocron: WithDistributedElector: elector must not be nil")
	ErrWithDistributedLockerNil      = fmt.Errorf("gocron: WithDistributedLocker: locker must not be nil")
	ErrWithDistributedJobLockerNil   = fmt.Errorf("gocron: WithDistributedJobLocker: locker must not be nil")
	ErrWithIdentifierNil             = fmt.Errorf("gocron: WithIdentifier: identifier must not be nil")
	ErrWithLimitConcurrentJobsZero   = fmt.Errorf("gocron: WithLimitConcurrentJobs: limit must be greater than 0")
	ErrWithLocationNil               = fmt.Errorf("gocron: WithLocation: location must not be nil")
	ErrWithLoggerNil                 = fmt.Errorf("gocron: WithLogger: logger must not be nil")
	ErrWithMonitorNil                = fmt.Errorf("gocron: WithMonitor: monitor must not be nil")
	ErrWithNameEmpty                 = fmt.Errorf("gocron: WithName: name must not be empty")
	ErrWithStartDateTimePast         = fmt.Errorf("gocron: WithStartDateTime: start must not be in the past")
	ErrWithStopDateTimePast          = fmt.Errorf("gocron: WithStopDateTime: end must not be in the past")
	ErrStartTimeLaterThanEndTime     = fmt.Errorf("gocron: WithStartDateTime: start must not be later than end")
	ErrStopTimeEarlierThanStartTime  = fmt.Errorf("gocron: WithStopDateTime: end must not be earlier than start")
	ErrWithStopTimeoutZeroOrNegative = fmt.Errorf("gocron: WithStopTimeout: timeout must be greater than 0")
)

// internal errors
var (
	errAtTimeNil    = fmt.Errorf("errAtTimeNil")
	errAtTimesNil   = fmt.Errorf("errAtTimesNil")
	errAtTimeHours  = fmt.Errorf("errAtTimeHours")
	errAtTimeMinSec = fmt.Errorf("errAtTimeMinSec")
)
