import { Router, Route } from "@solidjs/router";
import EmptyPage from "pages/EmptyPage";
import HomePage from "pages/HomePage";
import SongPage from "pages/SongPage";
import MainLayout from "./layouts/MainLayout";

const isProd = import.meta.env.MODE === "production";
const basename = isProd ? "/eefu" : "/";

function App() {
  return (
    <Router base={basename}>
      <Route path="/" component={MainLayout}>
        <Route path="/" component={HomePage} />
        <Route path="/song/:song_group_hash_id" component={SongPage} />
        <Route path="*" component={EmptyPage} />
      </Route>
    </Router>
  );
}

export default App;
