export default function Searchbar(props: {
  keyword: string;
  onInput: (keyword: string) => void;
}) {
  return (
    <div class="grid w-full grid-rows-[auto_1fr] gap-4">
      <input
        type="text"
        placeholder="Search for a song..."
        value={props.keyword}
        onInput={(e) => props.onInput(e.currentTarget.value)}
        class="w-full rounded-md border border-gray-300 p-2"
      />
    </div>
  );
}
