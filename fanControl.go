package pwmfan

import (
	"time"

	"github.com/eZioPan/pwmfan-go/common"
	rpio "github.com/stianeikeland/go-rpio"
)

//NewFan initialize a Fan object
func NewFan(cfg common.Config) (fan *common.Fan) {
	fan = new(common.Fan)
	fan.Trigger.State = common.Stop
	fan.Current.Cycle = 0
	fan.Cfg = cfg
	fan.Pin = fan.Cfg.Pin
	rpio.Pin(fan.Pin).Pwm()
	rpio.Pin(fan.Pin).Freq(int(fan.Cfg.PwmFreq))
	return fan
}

// Monitor is the function to control fan's real-time stat
func Monitor(fan *common.Fan) {
	for {
		fan.Current.Temp = common.ReadCPUTemperature(fan.Cfg.CPUTempPath, 1000)
		switch fan.Trigger.State {
		case common.Stop:
			if fan.Current.Temp >= fan.Cfg.Start.Temp && fan.Trigger.Count < fan.Cfg.Start.Count {
				fan.Trigger.Count++
			}
			if fan.Current.Temp < fan.Cfg.Start.Temp && fan.Trigger.Count > 0 {
				fan.Trigger.Count--
			}
			if fan.Trigger.Count >= fan.Cfg.Start.Count {
				fan.Trigger.State = common.Start
				fan.Trigger.Count = 0
				fan.Current.Count = 0
			}
		case common.Start:
			fan.Trigger.State = common.Run
			fan.Trigger.Count = 0
			fan.Current.Count = 0
		case common.Run:
			if fan.Current.Temp <= fan.Cfg.Low.Temp && fan.Trigger.Count < fan.Cfg.Low.Count {
				fan.Trigger.Count++
			}
			if fan.Current.Temp > fan.Cfg.Low.Temp && fan.Trigger.Count > 0 {
				fan.Trigger.Count--
			}
			if fan.Trigger.Count >= fan.Cfg.Low.Count {
				fan.Trigger.State = common.Stop
				fan.Trigger.Count = 0
				fan.Current.Count = 0
			}
		}
		UpdateCycleFromState(fan, common.LinearClampRemap)
		fan.Current.Count++
		rpio.Pin(fan.Pin).DutyCycle(fan.Current.Cycle, fan.Cfg.FullCycle)
		time.Sleep(time.Second / time.Duration(fan.Cfg.SampleRate))

		// Don't pour rubbish into system log
		// TODO: try use level classified log
	}
}

// UpdateCycleFromState update pwm fan's Cycle information from State information
func UpdateCycleFromState(fan *common.Fan, remapper common.RemapFunc) {
	switch fan.Trigger.State {
	case common.Stop:
		fan.Current.Cycle = fan.Cfg.StopCycle
	case common.Start:
		fan.Current.Cycle = fan.Cfg.Start.Cycle
	case common.Run:
		cycle := remapper([]float64{fan.Current.Temp},
			fan.Cfg.Low.Temp,
			fan.Cfg.High.Temp,
			float64(fan.Cfg.Low.Cycle),
			float64(fan.Cfg.High.Cycle),
		)[0]
		fan.Current.Cycle = uint32(cycle)
	}
}
