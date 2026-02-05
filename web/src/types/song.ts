type AudioMetadata = {
  sample_rate: number;
  channels: number;
  duration: number;
};

type VideoMetadata = {
  duration: number;
};

type ImageMetadata = {
  width: number;
  height: number;
};

export type Genre = {
  name: string;
};

export type Tag = {
  name: string;
  type: string;
  description?: string;
};

export type Singer = {
  name: string;
  aliases: string[];
};

export type SongSinger = {
  singer: Singer;
  role: "main" | "featured" | string;
};

export type MediaSource = {
  storage_type: "local" | "s3" | "youtube" | "soundcloud";
  url: string;
  file_type: "audio" | "video" | "art";
  format_type: "mp3" | "wav" | "flac" | "mp4" | "webm" | "png" | "jpg" | "jpeg";
  checksum?: string;
  audio_metadata?: AudioMetadata;
  video_metadata?: VideoMetadata;
  image_metadata?: ImageMetadata;
};

export type SongAlias = {
  name: string;
  language: string;
};

export type Song = {
  hash_id: string;
  title: string;
  description: string;
  type: "cover" | "original" | "remix";
  released_at: string | null;
  genres: Genre[];
  tags: Tag[];
  singers: SongSinger[];
  media_sources: MediaSource[];
  aliases: SongAlias[];
  version: number;
  song_group_hash_id: string;
};

export type SongMedia = {
  file_type: string;
  url: string;
};

export type SongThumbnail = {
  song_group_hash_id: string;
  title: string;
  cover_url: string;
};

export type SongThumbnailWindow = {
  total_songs: number;
  songs: SongThumbnail[];
};
