import { createEffect, createSignal, Show } from "solid-js";
import ControlPanel from "components/Audio/ControlPanel";
import Waveform from "components/Audio/Waveform";
import { useAudio } from "resources/hooks/useAudio";

export default function AudioPlayer(props: { src: string; title?: string }) {
  const [waveformLoaded, setWaveformLoaded] = createSignal(false);
  const {
    audioRef,
    isPlaying,
    progress,
    duration,
    ratio,
    togglePlay,
    handleSeek,
    handleLoadedMetadata,
    src,
  } = useAudio(() => props.src);

  // When src reloads, reset waveform
  createEffect(() => {
    src();
    setWaveformLoaded(false);
  });

  return (
    <div class="mx-auto w-full rounded-md border border-gray-300 bg-white p-4 shadow-sm md:max-w-lg">
      <Show when={src()} keyed>
        {(currentSrc) => (
          <Waveform
            src={currentSrc}
            progress={ratio}
            onSeek={handleSeek}
            width={600}
            height={100}
            onLoaded={() => setWaveformLoaded(true)}
          />
        )}
      </Show>
      <ControlPanel
        isPlaying={isPlaying()}
        togglePlay={togglePlay}
        progress={ratio()}
        currentTime={progress()}
        duration={duration()}
        title={props.title}
        disabled={!waveformLoaded()}
      />
      <audio
        ref={audioRef}
        src={src()}
        onLoadedMetadata={handleLoadedMetadata}
        loop
        class="hidden"
      />
    </div>
  );
}
