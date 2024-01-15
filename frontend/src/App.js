import { Routes, Route, BrowserRouter } from "react-router-dom";
import NavigationBar from "./components/NavigationBar";
import ProfileEditorPage from "./pages/ProfileEditorPage";
import PostsPage from "./pages/PostsPage";
import Signup from "./pages/SignupPage";
import Login from "./pages/LoginPage";
import NoPage from "./pages/NoPage";
import RequireAuth from "./components/RequireAuth";
import ProfileInfoPage from "./pages/ProfileInfoPage";
import GroupPage from "./pages/GroupPage";
import Logout from "./components/Logout";
import EventPage from "./pages/EventPage.js";
import RequireGuest from "./components/RequireGuest";

import "./style.css";

function App() {
  return (
    <BrowserRouter>
      <NavigationBar />
      <Routes>
        <Route element={<RequireGuest />}>
          <Route path="login" element={<Login />} />
          <Route path="signup" element={<Signup />} />
        </Route>

        <Route element={<RequireAuth />}>
          <Route path="/" element={<PostsPage />} />
          <Route path="profile" element={<ProfileEditorPage />} />
          <Route path="profile/:id" element={<ProfileInfoPage />} />
          <Route path="groups/:id" element={<GroupPage />} />
          <Route path="event/:id" element={<EventPage />} />

          <Route path="logout" element={<Logout />} />
          <Route path="*" element={<NoPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
