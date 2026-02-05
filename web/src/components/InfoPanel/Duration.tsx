import { createMemo } from "solid-js";

function Duration(props: { durationInSec: number }) {
  const formatTime = createMemo(() => {
    const totalSec = Math.floor(props.durationInSec);
    const h = Math.floor(totalSec / 3600);
    const m = Math.floor((totalSec % 3600) / 60);
    const s = totalSec % 60;

    const str = [h, m, s].map((n) => n.toString().padStart(2, "0")).join(":");

    let firstNonZeroFound = false;

    return str.split("").map((ch) => {
      if (!firstNonZeroFound && (ch === "0" || ch === ":")) {
        return <span class="text-gray-400">{ch}</span>;
      } else if (!firstNonZeroFound && ch !== "0" && ch !== ":") {
        firstNonZeroFound = true;
        return <span>{ch}</span>;
      } else {
        return <span>{ch}</span>;
      }
    });
  });

  return <div class="flex">{formatTime()}</div>;
}

export default Duration;
