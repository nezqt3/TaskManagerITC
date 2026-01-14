import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import StartScreen from "./pages/StartScreen";
import MainScreen from "./pages/MainScreen";
import DashboardLayout from "./pages/DashboardLayout";
import ProjectsScreen from "./pages/ProjectsScreen";
import ServicesScreen from "./pages/ServicesScreen";
import MembersScreen from "./pages/MembersScreen";
import ProfileScreen from "./pages/ProfileScreen";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<StartScreen />} />
        <Route element={<DashboardLayout />}>
          <Route path="/main" element={<MainScreen />} />
          <Route path="/projects" element={<ProjectsScreen />} />
          <Route path="/services" element={<ServicesScreen />} />
          <Route path="/members" element={<MembersScreen />} />
          <Route path="/profile" element={<ProfileScreen />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
