import { createEffect, createSignal, onCleanup } from "solid-js";

export function useAudio(src: () => string) {
  const [isPlaying, setIsPlaying] = createSignal(false);
  const [progress, setProgress] = createSignal(0);
  const [duration, setDuration] = createSignal(0);

  let audio: HTMLAudioElement;
  let animationFrame: number;

  createEffect(() => {
    const s = src();
    if (!s || !audio) return;

    // Reset playstate
    setIsPlaying(false);
    audio.pause();
    audio.currentTime = 0;
    setProgress(0);
  });

  const ratio = () => (duration() ? progress() / duration() : 0);

  const updateProgress = () => {
    if (audio) setProgress(audio.currentTime);
    animationFrame = requestAnimationFrame(updateProgress);
  };

  const togglePlay = () => {
    if (!audio) return;

    if (audio.paused) {
      audio.play();
      setIsPlaying(true);
      requestAnimationFrame(updateProgress);
    } else {
      audio.pause();
      setIsPlaying(false);
      cancelAnimationFrame(animationFrame);
    }
  };

  const handleSeek = (r: number) => {
    if (!audio) return;
    audio.currentTime = r * audio.duration;

    if (audio.paused) {
      audio.play();
      setIsPlaying(true);
      requestAnimationFrame(updateProgress);
    }
  };

  const handleLoadedMetadata = () => setDuration(audio.duration);

  onCleanup(() => {
    audio?.pause();
    cancelAnimationFrame(animationFrame);
  });

  return {
    audioRef: (el: HTMLAudioElement) => (audio = el),
    isPlaying,
    progress,
    duration,
    ratio,
    togglePlay,
    handleSeek,
    handleLoadedMetadata,
    src,
  };
}
