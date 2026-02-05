import { createSignal, createEffect } from "solid-js";
import { audioCtx } from "resources/audioContext";

// Cache waveform outside mounts
const waveformCache = new Map<string, Float32Array>();

export default function Waveform(props: {
  src: string;
  progress: () => number;
  onSeek?: (ratio: number) => void;
  onLoaded?: () => void;
  width: number;
  height: number;
}) {
  const [loaded, setLoaded] = createSignal(false);

  let canvasRef: HTMLCanvasElement | undefined;
  let offscreen: HTMLCanvasElement | null = null;
  let audioBuffer: AudioBuffer | null = null;

  const drawFromFilteredData = (
    filteredData: Float32Array,
    ctx: CanvasRenderingContext2D,
  ) => {
    const max = Math.max(...filteredData);

    offscreen = document.createElement("canvas");
    offscreen.width = props.width;
    offscreen.height = props.height;
    const offCtx = offscreen.getContext("2d");
    if (!offCtx) return;

    offCtx.clearRect(0, 0, props.width, props.height);
    offCtx.fillStyle = "#e0e0e0";
    filteredData.forEach((v, i) => {
      const y = (v / max) * (props.height / 2);
      offCtx.fillRect(i, props.height / 2 - y, 1, y * 2);
    });

    ctx.clearRect(0, 0, props.width, props.height);
    ctx.drawImage(offscreen, 0, 0);
  };

  const drawStaticWaveform = async (src: string) => {
    if (!canvasRef) return;
    const ctx = canvasRef.getContext("2d");
    if (!ctx) return;

    let filteredData: Float32Array;

    if (waveformCache.has(src)) {
      filteredData = waveformCache.get(src)!;
      drawFromFilteredData(filteredData, ctx);
      setLoaded(true);
      props.onLoaded?.();
      return;
    }

    if (!audioBuffer) {
      const arrayBuffer = await fetch(src).then((r) => r.arrayBuffer());
      audioBuffer = await audioCtx.decodeAudioData(arrayBuffer);
    }

    const rawData = audioBuffer.getChannelData(0);
    const samples = props.width;
    const blockSize = Math.floor(rawData.length / samples);
    filteredData = new Float32Array(samples);

    for (let i = 0; i < samples; i++) {
      let sum = 0;
      for (let j = 0; j < blockSize; j++)
        sum += Math.abs(rawData[i * blockSize + j]);
      filteredData[i] = sum / blockSize;
    }
    waveformCache.set(src, filteredData);

    drawFromFilteredData(filteredData, ctx);
    setLoaded(true);
    props.onLoaded?.();
  };

  const drawPlayhead = () => {
    if (!canvasRef || !offscreen) return;
    const ctx = canvasRef.getContext("2d");
    if (!ctx) return;

    ctx.clearRect(0, 0, props.width, props.height);
    const playheadX = props.progress() * props.width;

    ctx.drawImage(
      offscreen,
      0,
      0,
      playheadX,
      props.height,
      0,
      0,
      playheadX,
      props.height,
    );
    ctx.fillStyle = "rgba(160,160,160,0.2)";
    ctx.fillRect(0, 0, playheadX, props.height);

    ctx.drawImage(
      offscreen,
      playheadX,
      0,
      props.width - playheadX,
      props.height,
      playheadX,
      0,
      props.width - playheadX,
      props.height,
    );

    ctx.strokeStyle = "rgba(160,160,160)";
    ctx.beginPath();
    ctx.moveTo(playheadX, 0);
    ctx.lineTo(playheadX, props.height);
    ctx.stroke();
  };

  const handleSeek = (event: MouseEvent) => {
    if (!canvasRef) return;
    const rect = canvasRef.getBoundingClientRect();
    const clickX = event.clientX - rect.left;
    let ratio = clickX / rect.width;
    ratio = Math.min(Math.max(ratio, 0), 1);
    props.onSeek?.(ratio);
  };

  createEffect(() => {
    const src = props.src;
    if (!src) return;

    // reset state
    setLoaded(false);
    audioBuffer = null;
    drawStaticWaveform(src);
  });

  createEffect(() => {
    if (!loaded() || !offscreen) return;
    drawPlayhead();
  });

  return (
    <div class="relative w-full cursor-pointer" onClick={handleSeek}>
      {!loaded() && (
        <div class="absolute inset-0 animate-pulse rounded-md bg-gray-200" />
      )}
      <canvas
        ref={canvasRef!}
        width={props.width}
        height={props.height}
        class={`mb-2 w-full rounded-md border border-gray-200 transition-opacity duration-300 ${loaded() ? "opacity-100" : "opacity-0"}`}
      />
    </div>
  );
}
