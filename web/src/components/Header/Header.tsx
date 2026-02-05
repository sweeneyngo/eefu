import { Asterisk } from "components/Icons";

export default function Header() {
  return (
    <header class="flex w-full justify-between">
      <div class="flex">
        <Asterisk />
        <p class="font-bold">eefu</p>
        <p class="mx-2">—</p>
        <p>explore everything, free & unlimited</p>
      </div>
    </header>
  );
}
