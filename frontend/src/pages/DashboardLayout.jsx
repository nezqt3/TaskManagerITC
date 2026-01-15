import { NavLink, Outlet, useLocation, useSearchParams } from "react-router-dom";
import "../styles/DashboardLayout.scss";

const NAV_ITEMS = [
  { id: "main", label: "Главная", to: "/main" },
  { id: "projects", label: "Все проекты", to: "/projects" },
  { id: "services", label: "Сервисы", to: "/services" },
  { id: "members", label: "Участники", to: "/members" },
  { id: "profile", label: "Профиль", to: "/profile" },
];

export default function DashboardLayout() {
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const isProfileContext =
    location.pathname.startsWith("/profile") ||
    (location.pathname.startsWith("/projects/") &&
      searchParams.get("from") === "profile");
  const isProjectsContext =
    location.pathname === "/projects" || location.pathname.startsWith("/projects/");

  return (
    <div className="dashboard-layout">
      <header className="dashboard-header">
        <nav className="main-nav">
          {NAV_ITEMS.map((item) => (
            <NavLink
              key={item.id}
              to={item.to}
              className={({ isActive }) =>
                `main-nav__item${
                  (item.id === "profile" && isProfileContext) ||
                  (item.id === "projects" && !isProfileContext && isProjectsContext) ||
                  (!isProfileContext && !isProjectsContext && isActive)
                    ? " main-nav__item--active"
                    : ""
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
      </header>
      <main className="main-content" key={location.pathname}>
        <Outlet />
      </main>
    </div>
  );
}
