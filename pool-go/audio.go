package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Audio is synthesized at startup into raw PCM buffers, so the game ships no
// sound files. Ebiten's context plays signed 16-bit little-endian stereo at the
// context sample rate.
const audioSampleRate = 44100

var (
	audioCtx      *audio.Context
	sndClick      []byte // ball-to-ball contact
	sndRail       []byte // cushion / jaw
	sndPocket     []byte // ball dropping
	sndCue        []byte // cue strike
	activePlayers []*audio.Player
)

// initAudio creates the audio context and bakes the sound buffers. Safe to call
// once; further calls are ignored.
func initAudio() {
	if audioCtx != nil {
		return
	}
	audioCtx = audio.NewContext(audioSampleRate)
	sndClick = synthTone(0.06, 880, 55, 0.25)
	sndRail = synthTone(0.10, 190, 30, 0.12)
	sndPocket = synthSweep(0.28, 520, 150, 9)
	sndCue = mix(synthTone(0.05, 1200, 80, 0.2), synthTone(0.07, 130, 40, 0))
}

// play starts a fresh player for a buffer at the given volume, allowing sounds
// to overlap. Finished players are pruned so they can be collected.
func play(buf []byte, volume float64) {
	if audioCtx == nil || buf == nil {
		return
	}
	kept := activePlayers[:0]
	for _, p := range activePlayers {
		if p.IsPlaying() {
			kept = append(kept, p)
		}
	}
	activePlayers = kept

	p := audioCtx.NewPlayerFromBytes(buf)
	p.SetVolume(math.Max(0, math.Min(1, volume)))
	p.Play()
	activePlayers = append(activePlayers, p)
}

func playClick(strength float64)  { play(sndClick, math.Min(1, 0.15+strength*0.08)) }
func playRail(strength float64)   { play(sndRail, math.Min(0.8, 0.1+strength*0.06)) }
func playPocket()                 { play(sndPocket, 0.7) }
func playCueStrike(power float64) { play(sndCue, math.Min(1, 0.3+power*0.6)) }

// synthTone renders a decaying sine (optionally dusted with noise).
func synthTone(dur, freq, decay, noise float64) []byte {
	n := int(dur * audioSampleRate)
	buf := make([]byte, n*4)
	for i := 0; i < n; i++ {
		t := float64(i) / audioSampleRate
		env := math.Exp(-t * decay)
		v := math.Sin(2*math.Pi*freq*t) * env
		if noise > 0 {
			v += (rand.Float64()*2 - 1) * noise * env
		}
		writeStereo(buf, i, v)
	}
	return buf
}

// synthSweep renders a tone that glides from f0 to f1 over its duration.
func synthSweep(dur, f0, f1, decay float64) []byte {
	n := int(dur * audioSampleRate)
	buf := make([]byte, n*4)
	phase := 0.0
	for i := 0; i < n; i++ {
		t := float64(i) / audioSampleRate
		f := f0 + (f1-f0)*(t/dur)
		phase += 2 * math.Pi * f / audioSampleRate
		writeStereo(buf, i, math.Sin(phase)*math.Exp(-t*decay))
	}
	return buf
}

// mix sums two equal-rate buffers sample-for-sample (length of the longer one).
func mix(a, b []byte) []byte {
	if len(b) > len(a) {
		a, b = b, a
	}
	out := make([]byte, len(a))
	copy(out, a)
	for i := 0; i+1 < len(b); i += 2 {
		va := int16(uint16(out[i]) | uint16(out[i+1])<<8)
		vb := int16(uint16(b[i]) | uint16(b[i+1])<<8)
		s := int32(va) + int32(vb)
		if s > 32767 {
			s = 32767
		} else if s < -32768 {
			s = -32768
		}
		out[i] = byte(s)
		out[i+1] = byte(s >> 8)
	}
	return out
}

// writeStereo clamps v to [-1,1] and writes it to both channels of frame i.
func writeStereo(buf []byte, i int, v float64) {
	if v > 1 {
		v = 1
	} else if v < -1 {
		v = -1
	}
	s := int16(v * 32767)
	lo, hi := byte(uint16(s)), byte(uint16(s)>>8)
	buf[i*4] = lo
	buf[i*4+1] = hi
	buf[i*4+2] = lo
	buf[i*4+3] = hi
}
