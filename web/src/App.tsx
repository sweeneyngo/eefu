import { Router, Route } from "@solidjs/router";
import EmptyPage from "pages/EmptyPage";
import HomePage from "pages/HomePage";
import SongPage from "pages/SongPage";
import MainLayout from "./layouts/MainLayout";

function App() {
  return (
    <Router>
      <Route path="/" component={MainLayout}>
        <Route path="/" component={HomePage} />
        <Route path="/song/:song_group_hash_id" component={SongPage} />
        <Route path="*" component={EmptyPage} />
      </Route>
    </Router>
  );
}

export default App;
