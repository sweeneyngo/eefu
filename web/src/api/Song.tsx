import type { Song, SongMedia } from "types/song";

const BASE_URL =
  window.location.hostname === "localhost"
    ? "http://localhost:8080"
    : "https://eefu-api.fly.dev";

export async function fetchSongs(): Promise<Song[]> {
  const res = await fetch(`${BASE_URL}/songs`);
  if (!res.ok) {
    throw new Error(`Failed to fetch song: ${res.statusText}`);
  }
  const data: Song[] = await res.json();
  return data;
}

export async function fetchSong(songGroupHashID: string): Promise<Song[]> {
  const res = await fetch(
    `${BASE_URL}/songs/group/${songGroupHashID}/versions`,
  );
  if (!res.ok) {
    throw new Error(`Failed to fetch song: ${res.statusText}`);
  }
  const data: Song[] = await res.json();
  return data;
}

export async function fetchSongMedia(hash_id: string): Promise<SongMedia[]> {
  const res = await fetch(`${BASE_URL}/songs/${hash_id}/download`);
  if (!res.ok) {
    throw new Error(`Failed to fetch song ${hash_id}: ${res.statusText}`);
  }
  const data: SongMedia[] = await res.json();
  return data;
}
