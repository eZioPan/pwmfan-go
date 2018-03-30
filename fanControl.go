package pwmfan

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

//NewFan initialize a Fan object
func NewFan(cfg Config) (fan *Fan) {
	fan = new(Fan)
	fan.SetState(Stop)
	fan.SetCycle(0)
	fan.SetCfg(cfg)
	fan.Pin = rpio.Pin(fan.GetCfg().Pin)
	fan.Pin.Pwm()
	fan.Pin.Freq(int(fan.GetCfg().PwmFreq))
	return fan
}

// Monitor is the function to control fan's real-time stat
func (fan *Fan) Monitor() {
	for {
		fan.SetTemp(ReadCPUTemperature(fan.GetCfg().CPUTempPath, 1000))
		switch fan.GetState() {
		case Stop:
			if fan.GetTemp() >= fan.GetCfg().StartTemp && fan.GetStartCounter() < fan.GetCfg().StartCount {
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
			if fan.GetTemp() <= fan.GetCfg().LowTemp && fan.GetStopCounter() < fan.GetCfg().StopCount {
				fan.SetStopCounter(fan.GetStopCounter() + 1)
			} else if fan.GetTemp() >= fan.GetCfg().LowTemp && fan.GetStopCounter() > 0 {
				fan.SetStopCounter(fan.GetStopCounter() - 1)
			}
			if fan.GetStopCounter() >= fan.GetCfg().StopCount {
				fan.SetState(Stop)
				fan.SetStopCounter(0)
			}
		}
		fan.UpdateCycleFromState(LinearClampRemap)
		fan.Pin.DutyCycle(uint32(fan.GetCycle()), uint32(fan.GetCfg().FullCycle))
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

// SetState set a fan's current running state
func (fan *Fan) SetState(state FanState) {
	fan.State = state
}

// GetState get a fan's current running state
func (fan Fan) GetState() (state FanState) {
	return fan.State
}

// SetCfg set a fan's configuration
func (fan *Fan) SetCfg(cfg Config) {
	fan.Cfg = cfg
}

// GetCfg get a fan's configuration
func (fan Fan) GetCfg() (cfg Config) {
	return fan.Cfg
}

// SetCycle set fan's current pwm duty cycle information
func (fan *Fan) SetCycle(cycle uint) {
	fan.Cycle = cycle
}

// GetCycle get fan's current pwm duty cycle information
func (fan Fan) GetCycle() (cycle uint) {
	return fan.Cycle
}

// SetTemp set fan's current loaded temperature information
func (fan *Fan) SetTemp(temp float64) {
	fan.Temp = temp
}

// GetTemp get fan's current loaded temperature information
func (fan Fan) GetTemp() (temp float64) {
	return fan.Temp
}

// SetStartCounter set fan's current start counter
func (fan *Fan) SetStartCounter(num uint) {
	fan.StartCounter = num
}

// GetStartCounter get fan's current start counter
func (fan Fan) GetStartCounter() (num uint) {
	return fan.StartCounter
}

// SetStopCounter set fan's current stop counter
func (fan *Fan) SetStopCounter(num uint) {
	fan.StopCounter = num
}

// GetStopCounter get fan's current stop counter
func (fan Fan) GetStopCounter() (num uint) {
	return fan.StopCounter
}
