/**
 * AudioManager class
 * Handles all audio playback using Web Audio API
 */

export class AudioManager {
  private context: AudioContext | null = null;
  private isMuted: boolean = false;

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

  /** Play pocket sound effect */
  playPocketSound(): void {
    if (!this.context || this.isMuted) return;

    const now = this.context.currentTime;
    const oscillator = this.context.createOscillator();
    const gainNode = this.context.createGain();

    oscillator.type = "triangle";
    oscillator.frequency.setValueAtTime(740, now);
    oscillator.frequency.exponentialRampToValueAtTime(370, now + 0.08);

    gainNode.gain.setValueAtTime(0.0001, now);
    gainNode.gain.exponentialRampToValueAtTime(0.055, now + 0.01);
    gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.085);

    oscillator.connect(gainNode);
    gainNode.connect(this.context.destination);

    oscillator.start(now);
    oscillator.stop(now + 0.09);
  }

  /** Play cue ball hit sound */
  playCueHitSound(power: number): void {
    if (!this.context || this.isMuted) return;

    const now = this.context.currentTime;
    const oscillator = this.context.createOscillator();
    const gainNode = this.context.createGain();

    // Higher power = higher pitch and louder
    const baseFreq = 200 + power * 10;
    oscillator.type = "sine";
    oscillator.frequency.setValueAtTime(baseFreq, now);
    oscillator.frequency.exponentialRampToValueAtTime(baseFreq * 0.5, now + 0.05);

    const volume = Math.min(0.1 + power * 0.005, 0.3);
    gainNode.gain.setValueAtTime(0.0001, now);
    gainNode.gain.exponentialRampToValueAtTime(volume, now + 0.01);
    gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.08);

    oscillator.connect(gainNode);
    gainNode.connect(this.context.destination);

    oscillator.start(now);
    oscillator.stop(now + 0.1);
  }

  /** Play ball collision sound */
  playCollisionSound(impactForce: number): void {
    if (!this.context || this.isMuted) return;

    const now = this.context.currentTime;
    const oscillator = this.context.createOscillator();
    const gainNode = this.context.createGain();

    oscillator.type = "triangle";
    oscillator.frequency.setValueAtTime(300 + impactForce * 20, now);
    oscillator.frequency.exponentialRampToValueAtTime(150, now + 0.05);

    const volume = Math.min(0.05 + impactForce * 0.01, 0.15);
    gainNode.gain.setValueAtTime(0.0001, now);
    gainNode.gain.exponentialRampToValueAtTime(volume, now + 0.005);
    gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.06);

    oscillator.connect(gainNode);
    gainNode.connect(this.context.destination);

    oscillator.start(now);
    oscillator.stop(now + 0.07);
  }

  /** Play cushion bounce sound */
  playCushionSound(): void {
    if (!this.context || this.isMuted) return;

    const now = this.context.currentTime;
    const oscillator = this.context.createOscillator();
    const gainNode = this.context.createGain();

    oscillator.type = "sawtooth";
    oscillator.frequency.setValueAtTime(150, now);
    oscillator.frequency.exponentialRampToValueAtTime(80, now + 0.08);

    gainNode.gain.setValueAtTime(0.0001, now);
    gainNode.gain.exponentialRampToValueAtTime(0.04, now + 0.01);
    gainNode.gain.exponentialRampToValueAtTime(0.0001, now + 0.1);

    oscillator.connect(gainNode);
    gainNode.connect(this.context.destination);

    oscillator.start(now);
    oscillator.stop(now + 0.12);
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
