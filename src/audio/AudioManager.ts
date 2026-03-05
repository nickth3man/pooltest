/**
 * AudioManager - Web Audio API sound synthesis
 * 
 * All sounds are procedurally generated (no external audio files).
 * Uses oscillators with envelope shaping for realistic billiard sounds.
 */
export class AudioManager {
  private context: AudioContext | null = null;
  private isMuted: boolean = false;

  private withAudioContext(effect: (context: AudioContext) => void): void {
    if (!this.context || this.isMuted) return;
    effect(this.context);
  }

  private playSynth(
    context: AudioContext,
    configure: (oscillator: OscillatorNode, gainNode: GainNode, now: number) => number
  ): void {
    const now = context.currentTime;
    const oscillator = context.createOscillator();
    const gainNode = context.createGain();

    // Allow caller to configure oscillator and envelope
    const duration = configure(oscillator, gainNode, now);

    // Chain: oscillator -> gain -> speakers
    oscillator.connect(gainNode);
    gainNode.connect(context.destination);

    oscillator.start(now);
    oscillator.stop(now + duration);
  }

  constructor() {
    // Audio context is created lazily on first user interaction
  }

  /** Ensure audio context is initialized */
  ensureContext(): void {
    if (this.context) {
      if (this.context.state === "suspended") {
        this.context.resume();
      }
      return;
    }

    const AC = window.AudioContext || (window as unknown as { webkitAudioContext: typeof AudioContext }).webkitAudioContext;
    if (AC) {
      this.context = new AC();
    }
  }

  /** Play pocket sound effect (ball falling in pocket) */
  playPocketSound(): void {
    this.withAudioContext((context) => {
      this.playSynth(context, (oscillator, gainNode, now) => {
        // Falling thud: triangle wave with pitch drop
        oscillator.type = "triangle";
        oscillator.frequency.setValueAtTime(740, now);              // Start high
        oscillator.frequency.exponentialRampToValueAtTime(370, now + 0.08);  // Drop pitch

        // Quick attack/decay envelope
        gainNode.gain.setValueAtTime(0.0001, now);                  // Silent start
        gainNode.gain.exponentialRampToValueAtTime(0.055, now + 0.01);  // Attack
        gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.085); // Decay

        return 0.09;  // Duration in seconds
      });
    });
  }

  /** Play cue ball hit sound */
  playCueHitSound(power: number): void {
    this.withAudioContext((context) => {
      this.playSynth(context, (oscillator, gainNode, now) => {
        const baseFreq = 200 + power * 10;
        oscillator.type = "sine";
        oscillator.frequency.setValueAtTime(baseFreq, now);
        oscillator.frequency.exponentialRampToValueAtTime(baseFreq * 0.5, now + 0.05);

        const volume = Math.min(0.1 + power * 0.005, 0.3);
        gainNode.gain.setValueAtTime(0.0001, now);
        gainNode.gain.exponentialRampToValueAtTime(volume, now + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.08);

        return 0.1;
      });
    });
  }

  /** Play ball collision sound */
  playCollisionSound(impactForce: number): void {
    this.withAudioContext((context) => {
      this.playSynth(context, (oscillator, gainNode, now) => {
        oscillator.type = "triangle";
        oscillator.frequency.setValueAtTime(300 + impactForce * 20, now);
        oscillator.frequency.exponentialRampToValueAtTime(150, now + 0.05);

        const volume = Math.min(0.05 + impactForce * 0.01, 0.15);
        gainNode.gain.setValueAtTime(0.0001, now);
        gainNode.gain.exponentialRampToValueAtTime(volume, now + 0.005);
        gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.06);

        return 0.07;
      });
    });
  }

  /** Play cushion bounce sound */
  playCushionSound(): void {
    this.withAudioContext((context) => {
      this.playSynth(context, (oscillator, gainNode, now) => {
        oscillator.type = "sawtooth";
        oscillator.frequency.setValueAtTime(150, now);
        oscillator.frequency.exponentialRampToValueAtTime(80, now + 0.08);

        gainNode.gain.setValueAtTime(0.0001, now);
        gainNode.gain.exponentialRampToValueAtTime(0.04, now + 0.01);
        gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.1);

        return 0.12;
      });
    });
  }

  /** Mute/unmute all audio */
  setMuted(muted: boolean): void {
    this.isMuted = muted;
  }

  /** Get mute state */
  getMuted(): boolean {
    return this.isMuted;
  }

  /** Toggle mute state */
  toggleMute(): boolean {
    this.isMuted = !this.isMuted;
    return this.isMuted;
  }

  /** Check if audio is available */
  isAvailable(): boolean {
    return this.context !== null;
  }

  /** Clean up resources */
  destroy(): void {
    if (this.context && this.context.state !== "closed") {
      this.context.close();
    }
    this.context = null;
  }
}
