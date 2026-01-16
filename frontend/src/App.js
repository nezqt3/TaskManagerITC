import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import StartScreen from "./pages/StartScreen";
import MainScreen from "./pages/MainScreen";
import DashboardLayout from "./pages/DashboardLayout";
import ProjectsScreen from "./pages/ProjectsScreen";
import ProjectDetailsScreen from "./pages/ProjectDetailsScreen";
import ProjectManagementScreen from "./pages/ProjectManagementScreen";
import ServicesScreen from "./pages/ServicesScreen";
import MembersScreen from "./pages/MembersScreen";
import ProfileScreen from "./pages/ProfileScreen";
import { getJwt, getProfile } from "./utils/auth";

function App() {
  const isAuthorized = Boolean(getJwt() && getProfile());
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<StartScreen />} />
        <Route
          element={
            isAuthorized ? <DashboardLayout /> : <Navigate to="/" replace />
          }
        >
          <Route path="/main" element={<MainScreen />} />
          <Route path="/projects" element={<ProjectsScreen />} />
          <Route path="/projects/:id" element={<ProjectDetailsScreen />} />
          <Route
            path="/projects/:id/manage"
            element={<ProjectManagementScreen />}
          />
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
