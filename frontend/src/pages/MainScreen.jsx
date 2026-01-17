import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import "../styles/MainScreen.scss";
import { getAuthHeaders, getProfile } from "../utils/auth";
import { apiFetch } from "../utils/api";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";

const formatDate = (value) => {
  if (!value) {
    return "—";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleDateString("ru-RU");
};

const getPriorityLabel = (deadline) => {
  if (!deadline) {
    return "Обычно";
  }
  const date = new Date(deadline);
  if (Number.isNaN(date.getTime())) {
    return "Обычно";
  }
  const diff = date.getTime() - Date.now();
  const days = diff / (1000 * 60 * 60 * 24);
  return days <= 3 ? "Важно" : "Обычно";
};

export default function MainScreen() {
  const profile = useMemo(() => getProfile(), []);
  const [dashboard, setDashboard] = useState({
    tasks: [],
    projects: [],
    events: [],
  });
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const username = profile?.username?.replace(/^@/, "");
    if (!username) {
      setIsLoading(false);
      return;
    }

    let isActive = true;
    setIsLoading(true);
    setError("");

    apiFetch(`${API_BASE}/dashboard?username=${encodeURIComponent(username)}`, {
      headers: getAuthHeaders(),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return response.json();
      })
      .then((data) => {
        if (!isActive) {
          return;
        }
        setDashboard({
          tasks: Array.isArray(data.tasks) ? data.tasks : [],
          projects: Array.isArray(data.projects) ? data.projects : [],
          events: Array.isArray(data.events) ? data.events : [],
        });
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setError("Не удалось загрузить данные дашборда");
      })
      .finally(() => {
        if (!isActive) {
          return;
        }
        setIsLoading(false);
      });

    return () => {
      isActive = false;
    };
  }, [profile]);

  const normalizedUsername = profile?.username
    ? profile.username.replace(/^@/, "").toLowerCase()
    : "";

  const assignedTasks = useMemo(() => {
    if (!normalizedUsername) {
      return [];
    }
    return dashboard.tasks.filter(
      (task) =>
        (task.user || "").replace(/^@/, "").toLowerCase() === normalizedUsername
    );
  }, [dashboard, normalizedUsername]);
  const taskCards = useMemo(() => assignedTasks.slice(0, 6), [assignedTasks]);
  const activeProjects = useMemo(
    () =>
      dashboard.projects.filter(
        (project) => !(project.status || "").toLowerCase().includes("выполн")
      ),
    [dashboard]
  );
  const projectCards = useMemo(() => activeProjects.slice(0, 6), [activeProjects]);
  const eventCards = useMemo(() => dashboard.events.slice(0, 6), [dashboard]);

  return (
    <div className="dashboard">
      <section className="dashboard-panel">
        <header className="dashboard-panel__header">
          <h3>Текущие задачи:</h3>
          <span>{assignedTasks.length}</span>
        </header>
        <div className="dashboard-panel__list">
          {isLoading && <div className="dashboard-card">Загрузка...</div>}
          {!isLoading && error && <div className="dashboard-card">{error}</div>}
          {!isLoading && !error && taskCards.length === 0 && (
            <div className="dashboard-card">Нет активных задач</div>
          )}
          {!isLoading &&
            !error &&
            taskCards.map((task) => {
              const priority = getPriorityLabel(task.deadline);
              const taskLink = `/projects/${task.id_project}?from=profile&taskId=${task.id}`;
              return (
                <Link
                  className="dashboard-card dashboard-card--link"
                  key={task.id}
                  to={taskLink}
                >
                  <div className="dashboard-card__meta">
                    <span>Срок выполнения: {formatDate(task.deadline)}</span>
                    <span className={priority === "Важно" ? "badge badge--hot" : "badge"}>
                      {priority}
                    </span>
                  </div>
                  <h4>{task.title || "Без названия"}</h4>
                  <p className="dashboard-card__sub">
                    Создал: {task.author || "—"}
                  </p>
                  <p className="dashboard-card__sub">
                    Проект: {task.project_title || "—"}
                  </p>
                </Link>
              );
            })}
        </div>
      </section>

      <section className="dashboard-panel">
        <header className="dashboard-panel__header">
          <h3>Текущие проекты:</h3>
          <span>{activeProjects.length}</span>
        </header>
        <div className="dashboard-panel__list">
          {isLoading && <div className="dashboard-card">Загрузка...</div>}
          {!isLoading && error && <div className="dashboard-card">{error}</div>}
          {!isLoading && !error && projectCards.length === 0 && (
            <div className="dashboard-card">Нет активных проектов</div>
          )}
          {!isLoading &&
            !error &&
            projectCards.map((project) => {
              const members = Array.isArray(project.members)
                ? project.members
                : [];
              const member = members.find(
                (item) => item.username?.toLowerCase() === normalizedUsername
              );
              const role = member?.role || "участник";
              return (
                <article className="dashboard-card" key={project.id}>
                  <p className="dashboard-card__status">
                    Статус: {project.status || "в работе"}
                  </p>
                  <h4>{project.title}</h4>
                  <p className="dashboard-card__sub">Роль: {role}</p>
                </article>
              );
            })}
        </div>
      </section>

      <section className="dashboard-panel dashboard-panel--compact">
        <header className="dashboard-panel__header">
          <h3>События:</h3>
          <span>{dashboard.events.length}</span>
        </header>
        <div className="dashboard-panel__list">
          {isLoading && <div className="dashboard-card">Загрузка...</div>}
          {!isLoading && error && <div className="dashboard-card">{error}</div>}
          {!isLoading && !error && eventCards.length === 0 && (
            <div className="dashboard-card">Нет событий</div>
          )}
          {!isLoading &&
            !error &&
            eventCards.map((event) => (
              <article className="dashboard-card" key={event.id}>
                <div className="dashboard-card__meta">
                  <span>{event.date}</span>
                  <span>{event.time_range}</span>
                </div>
                <h4>{event.title}</h4>
                <p className="dashboard-card__sub">
                  Создан: {event.created_by}
                </p>
              </article>
            ))}
        </div>
      </section>
    </div>
  );
}
