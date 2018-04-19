package pwmfan

import (
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

//NewFan initialize a Fan object
func NewFan(cfg Config) (fan *Fan) {
	fan = new(Fan)
	fan.SetState(Stop)
	fan.SetCycle(0)
	fan.SetCfg(cfg)
	rpio.Pin(fan.Pin).Pwm()
	rpio.Pin(fan.Pin).Freq(int(fan.GetCfg().PwmFreq))
	return fan
}

// Monitor is the function to control fan's real-time stat
func (fan *Fan) Monitor() {
	for {
		fan.SetTemp(ReadCPUTemperature(fan.GetCfg().CPUTempPath, 1000))
		switch fan.GetState() {
		case Stop:
			if fan.GetTemp() >= fan.GetCfg().StartTemp && fan.StartCounter < fan.GetCfg().StartCount {
				fan.SetStartCounter(fan.GetStartCounter() + 1)
			} else if fan.GetTemp() <= fan.GetCfg().StartTemp && fan.GetStartCounter() > 0 {
				fan.SetStartCounter(fan.GetStartCounter() - 1)
			}
			if fan.GetStartCounter() >= fan.GetCfg().StartCount {
				fan.SetState(Start)
				fan.SetStartCounter(0)
			}
		case Start:
			fan.SetState(Run)
		case Run:
			if fan.GetTemp() <= fan.GetCfg().LowTemp && fan.StopCounter < fan.GetCfg().StopCount {
				fan.SetStopCounter(fan.GetStopCounter() + 1)
			} else if fan.GetTemp() >= fan.GetCfg().LowTemp && fan.StopCounter > 0 {
				fan.SetStopCounter(fan.GetStopCounter() - 1)
			}
			if fan.GetStopCounter() >= fan.GetCfg().StopCount {
				fan.SetState(Stop)
				fan.SetStopCounter(0)
			}
		}
		fan.UpdateCycleFromState(LinearClampRemap)
		rpio.Pin(fan.Pin).DutyCycle(uint32(fan.GetCycle()), uint32(fan.GetCfg().FullCycle))
		time.Sleep(time.Second / time.Duration(fan.GetCfg().SampleRate))

		// Don't pour rubbish into system log
		// TODO: try use level classified log

	}
}

// UpdateCycleFromState update pwm fan's Cycle information from State information
func (fan *Fan) UpdateCycleFromState(remapper RemapFunc) {
	switch fan.GetState() {
	case Stop:
		fan.SetCycle(fan.GetCfg().StopCycle)
	case Start:
		fan.SetCycle(fan.GetCfg().StartCycle)
	case Run:
		cycle := remapper([]float64{fan.GetTemp()},
			fan.GetCfg().LowTemp,
			fan.GetCfg().HighTemp,
			float64(fan.GetCfg().LowCycle),
			float64(fan.GetCfg().HighCycle),
		)[0]
		fan.SetCycle(uint(cycle))
	}
}
