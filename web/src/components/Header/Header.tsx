import { useNavigate } from "@solidjs/router";
import { Asterisk } from "components/Icons";

export default function Header() {
  const navigate = useNavigate();
  return (
    <header
      class="flex w-full cursor-pointer justify-between"
      onClick={() => navigate("/")}
    >
      <div class="flex">
        <Asterisk />
        <p class="font-bold">eefu</p>
        <p class="mx-2">—</p>
        <p>explore everything, free & unlimited</p>
      </div>
    </header>
  );
}
