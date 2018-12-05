package frontends

import (
	"fmt"
	"github.com/racerxdl/radioserver/protocol"
	"github.com/racerxdl/spy2go/airspy"
	"github.com/racerxdl/spy2go/spytypes"
	"math"
)

const airspyMaximumFrequency = 1768000
const airspyMinimumFrequency = 24000

type AirspyFrontend struct {
	device *airspy.Device
	cb SamplesCallback

	deviceSerial uint64
	maxSampleRate uint32
	maxDecimationStage uint32
}

func CreateAirspyFrontend(serial uint64) Frontend {
	airspy.Initialize()
	var f = &AirspyFrontend{
		device: airspy.MakeAirspyDevice(serial),
		deviceSerial: 0,
		maxSampleRate: 0,
	}

	f.device.SetSampleType(spytypes.SamplesComplex64)

	if f.deviceSerial == 0 {
		// Fetch device serial
		f.deviceSerial = f.device.GetSerial()
	}

	for _, v := range f.device.GetAvailableSampleRates() {
		if v > f.maxSampleRate {
			f.maxSampleRate = v
		}
	}

	var maxDecimationStage = uint32(0)
	var calcSR = f.maxSampleRate

	for calcSR >= minimumSampleRate {
		maxDecimationStage += 1
		var decim = uint32(math.Pow(float64(maxDecimationStage), 2))
		calcSR = f.maxSampleRate / decim
	}

	f.maxDecimationStage = maxDecimationStage

	return f
}

func (f *AirspyFrontend) GetUintDeviceSerial() uint32 {
	return uint32(f.deviceSerial & 0xFFFFFFFF)
}

func (f *AirspyFrontend) MinimumFrequency() uint32 {
	return airspyMinimumFrequency
}

func (f *AirspyFrontend) MaximumFrequency() uint32 {
	return airspyMaximumFrequency
}

func (f *AirspyFrontend) GetMaximumBandwidth() uint32 {
	return uint32(float32(f.maxSampleRate) * 0.8)
}

func (f *AirspyFrontend) MaximumGainIndex() uint32 {
	return 16
}

func (f *AirspyFrontend) MaximumDecimationStages() uint32 {
	return f.maxDecimationStage
}

func (f *AirspyFrontend) GetDeviceType() uint32 {
	return protocol.DeviceAirspyOne
}

func (f *AirspyFrontend) internalCb(dType int, data interface{}) {
	if dType != spytypes.SamplesComplex64 {
		panic("Spy2Go Library is sending different types than we asked!")
	}

	samples := data.([]complex64)

	if f.cb != nil {
		f.cb(samples)
	}
}

func (f *AirspyFrontend) GetDeviceSerial() string {
	return fmt.Sprintf("%08x", f.deviceSerial)
}
func (f *AirspyFrontend) GetMaximumSampleRate() uint32 {
	return f.maxSampleRate
}
func (f *AirspyFrontend) SetSampleRate(sampleRate uint32) uint32 {
	f.device.SetSampleRate(sampleRate)
	return f.device.GetSampleRate()
}
func (f *AirspyFrontend) SetCenterFrequency(centerFrequency uint32) uint32 {
	f.device.SetCenterFrequency(centerFrequency)
	return f.device.GetCenterFrequency()
}
func (f *AirspyFrontend) GetAvailableSampleRates() []uint32 {
	return f.device.GetAvailableSampleRates()
}
func (f *AirspyFrontend) Start() {
	f.device.Start()
}
func (f *AirspyFrontend) Stop() {
	f.device.Stop()
}
func (f *AirspyFrontend) SetAntenna(value string) {
	// Nothing
}
func (f *AirspyFrontend) SetAGC(agc bool) {
	f.device.SetAGC(agc)
}
func (f *AirspyFrontend) SetGain(value uint8) {
	f.device.SetLinearityGain(value)
}
func (f *AirspyFrontend) SetBiasT(value bool) {
	f.device.SetBiasT(value)
}
func (f *AirspyFrontend) GetCenterFrequency() uint32 {
	return f.device.GetSampleRate()
}
func (f *AirspyFrontend) GetName() string {
	return f.device.GetName()
}
func (f *AirspyFrontend) GetShortName() string {
	return "Airspy"
}
func (f *AirspyFrontend) GetSampleRate() uint32 {
	return f.device.GetSampleRate()
}
func (f *AirspyFrontend) SetSamplesAvailableCallback(cb SamplesCallback) {
	f.cb = cb
}
func (f *AirspyFrontend) Init() bool {
	return true
}

func (f *AirspyFrontend) Destroy() {
	// Nothing
}