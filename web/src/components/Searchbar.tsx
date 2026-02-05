export default function SearchBar(props: {
  value: string;
  onInput: (value: string) => void;
  placeholder?: string;
}) {
  return (
    <input
      type="text"
      value={props.value}
      placeholder={props.placeholder || "Let's find some songs..."}
      onInput={(e) => props.onInput(e.currentTarget.value)}
      class="mb-2 w-full rounded-md border border-gray-400 px-2 py-1 text-sm focus:border-gray-200 focus:ring-1 focus:ring-gray-800 focus:outline-none"
    />
  );
}
