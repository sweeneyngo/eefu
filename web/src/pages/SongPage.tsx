import { createEffect, createMemo, createSignal, Show } from "solid-js";
import { useParams } from "@solidjs/router";
import { fetchSong, fetchSongMedia } from "api/Song";
import type { Song } from "types/song";
import AudioPlayer from "components/Audio/AudioPlayer";
import CoverArt from "components/CoverArt/CoverArt";
import InfoCard from "components/InfoPanel/InfoCard";
import NotFound from "components/NotFound";
import { VersionCards } from "components/Audio/VersionCards";

export default function SongPage() {
  const params = useParams();
  const songGroupHashID = () => params.song_group_hash_id;

  const [song, setSong] = createSignal<Song[]>([]);
  const [audioUrl, setAudioUrl] = createSignal<string>("");
  const [artUrl, setArtUrl] = createSignal<string>("");
  const [loaded, setLoaded] = createSignal(false);
  const [selectedVersion, setSelectedVersion] = createSignal(0);

  const selectedSong = createMemo(() => song()[selectedVersion()]);

  createEffect(async () => {
    const songGroupHashId = songGroupHashID();
    if (!songGroupHashId) return;

    setLoaded(false);

    try {
      const data = await fetchSong(songGroupHashId);
      if (!data || data.length === 0) return;

      const presignedData = await fetchSongMedia(data[0].hash_id);
      if (!presignedData || presignedData.length === 0) return;

      const audioMedia = presignedData.find((m) => m.file_type === "audio");
      const artMedia = presignedData.find(
        (m) => m.file_type === "art" && m.url.endsWith("webp"),
      );
      if (!audioMedia || !artMedia) return;

      setSong(data);
      setAudioUrl(audioMedia.url);
      setArtUrl(artMedia.url);
      setSelectedVersion(0);
    } catch (err) {
      console.error(err);
    } finally {
      setLoaded(true);
    }
  });

  // When updating version, re-fetch presignedData
  async function updateSelectedVersion(v: number) {
    const presignedData = await fetchSongMedia(song()[v].hash_id);
    const audioMedia = presignedData.find((m) => m.file_type === "audio");
    const artMedia = presignedData.find((m) => m.file_type === "art");

    if (!audioMedia || !artMedia) return;

    setAudioUrl(audioMedia.url);
    setArtUrl(artMedia.url);

    setSelectedVersion(() => {
      const length = song()?.length ?? 1;
      return Math.min(Math.max(v, 0), length - 1);
    });
  }

  const handleDownload = () => {
    if (!audioUrl()) return () => {};
    console.log("downloading", audioUrl());
    return () => {
      window.open(audioUrl(), "_blank");
    };
  };

  return (
    <div class="mt-2">
      <Show
        when={song()?.length > 0 && audioUrl() && artUrl()}
        fallback={loaded() ? <NotFound message="Song not found" /> : <></>}
      >
        <div class="grid grid-cols-1 gap-4 md:grid-cols-[32rem_1fr]">
          <div class="flex flex-col">
            <CoverArt src={artUrl()} />
            <Show when={audioUrl() && selectedSong()}>
              <AudioPlayer src={audioUrl()} title={selectedSong().title} />
            </Show>
            <VersionCards
              song={song()}
              version={selectedVersion()}
              onSelectVersion={(version) => updateSelectedVersion(version)}
            />
          </div>
          <Show when={selectedSong()}>
            <InfoCard song={selectedSong()} onDownload={handleDownload()} />
          </Show>
        </div>
      </Show>
    </div>
  );
}
