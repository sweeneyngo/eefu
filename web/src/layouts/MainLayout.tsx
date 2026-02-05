import Header from "components/Header/Header";

export default function MainLayout(props: { children?: any }) {
  return (
    <div class="mx-auto max-w-[1200px] px-3 py-3 lg:px-8">
      <Header />
      {props.children}
    </div>
  );
}
