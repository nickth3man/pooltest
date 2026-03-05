import { describe, expect, it, vi, afterEach } from "vitest";
import { AudioManager } from "../audio/AudioManager.js";

class MockAudioContext {
  public currentTime = 0;
  public state: AudioContextState = "running";
  public destination = {} as AudioNode;

  createOscillator(): OscillatorNode {
    const frequency = {
      setValueAtTime: vi.fn(),
      exponentialRampToValueAtTime: vi.fn()
    } as unknown as AudioParam;

    return {
      type: "sine",
      frequency,
      connect: vi.fn(),
      start: vi.fn(),
      stop: vi.fn()
    } as unknown as OscillatorNode;
  }

  createGain(): GainNode {
    const gain = {
      setValueAtTime: vi.fn(),
      exponentialRampToValueAtTime: vi.fn()
    } as unknown as AudioParam;

    return {
      gain,
      connect: vi.fn()
    } as unknown as GainNode;
  }

  close(): Promise<void> {
    this.state = "closed";
    return Promise.resolve();
  }

  resume(): Promise<void> {
    this.state = "running";
    return Promise.resolve();
  }
}

describe("AudioManager", () => {
  afterEach(() => {
    delete (window as unknown as { AudioContext?: typeof AudioContext }).AudioContext;
  });

  it("creates audio context lazily and reports availability", () => {
    (window as unknown as { AudioContext: typeof AudioContext }).AudioContext = MockAudioContext as unknown as typeof AudioContext;

    const manager = new AudioManager();
    expect(manager.isAvailable()).toBe(false);

    manager.ensureContext();
    expect(manager.isAvailable()).toBe(true);
  });

  it("respects mute state when playing sounds", () => {
    (window as unknown as { AudioContext: typeof AudioContext }).AudioContext = MockAudioContext as unknown as typeof AudioContext;

    const manager = new AudioManager();
    manager.ensureContext();
    manager.setMuted(true);

    expect(() => manager.playPocketSound()).not.toThrow();
    expect(() => manager.playCueHitSound(0.6)).not.toThrow();
    expect(() => manager.playCollisionSound(0.5)).not.toThrow();
    expect(() => manager.playCushionSound()).not.toThrow();
    expect(manager.getMuted()).toBe(true);
  });
});
