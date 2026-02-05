import { Switch, Match } from "solid-js";
import { Audio, Video, Image } from "components/Icons";

export default function MediaTypeIcon(props: { type?: string | null }) {
  return (
    <Switch>
      <Match when={props.type === "audio"}>
        <Audio />
      </Match>
      <Match when={props.type === "video"}>
        <Video />
      </Match>
      <Match when={props.type === "image"}>
        <Image />
      </Match>
    </Switch>
  );
}
