export default function NotFound(props: { message: string }) {
  return (
    <div class="flex flex-col items-center justify-center py-20 text-center">
      <h1 class="text-6xl">404</h1>
      <h2 class="mb-6 text-xl text-gray-500">
        {props.message || "Page not found"}
      </h2>
      <p class="mx-auto max-w-md text-gray-400">
        Sorry, we couldn&apos;t find what you were looking for. Maybe it
        wandered off into the void, or perhaps it never existed. Either way,
        let&apos;s get you back on track.
      </p>
      <a
        href="/"
        class="my-4 rounded bg-blue-400 px-4 py-2 text-white transition hover:bg-purple-600"
      >
        Go Home
      </a>
    </div>
  );
}
