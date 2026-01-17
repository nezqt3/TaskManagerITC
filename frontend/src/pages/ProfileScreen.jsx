import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import "../styles/ProfileScreen.scss";
import { getAuthHeaders, getProfile } from "../utils/auth";
import { apiFetch } from "../utils/api";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";
const TASKS_LINK_TEXT = "просмотреть задачу";

export default function ProfileScreen() {
  const [profile] = useState(() => getProfile());

  const derivedProfile = useMemo(() => {
    const firstName =
      profile?.first_name ||
      (profile?.full_name ? profile.full_name.split(" ")[0] : "");
    const lastName =
      profile?.last_name ||
      (profile?.full_name ? profile.full_name.split(" ").slice(1).join(" ") : "");
    const fullName =
      profile?.full_name ||
      [profile?.first_name, profile?.last_name].filter(Boolean).join(" ");

    return {
      avatar:
        profile?.photo_url ||
        "https://images.unsplash.com/photo-1504593811423-6dd665756598?auto=format&fit=crop&w=400&q=70",
      handle: profile?.username ? `@${profile.username}` : "—",
      fullName: fullName || "—",
      firstName: firstName || "—",
      lastName: lastName || "—",
      dateOfBirthday: profile?.date_of_birthday || "—",
      role: profile?.role || "—",
      telegramID: profile?.telegram_id || "—",
    };
  }, [profile]);

  const [projects, setProjects] = useState([]);
  const [projectsError, setProjectsError] = useState("");
  const [isLoadingProjects, setIsLoadingProjects] = useState(true);

  useEffect(() => {
    const username = profile?.username?.replace(/^@/, "");
    if (!username) {
      setProjects([]);
      setIsLoadingProjects(false);
      return;
    }

    let isActive = true;
    setIsLoadingProjects(true);
    setProjectsError("");

    apiFetch(`${API_BASE}/projects?username=${encodeURIComponent(username)}`, {
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
        setProjects(Array.isArray(data) ? data : []);
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setProjectsError("Не удалось загрузить проекты");
      })
      .finally(() => {
        if (!isActive) {
          return;
        }
        setIsLoadingProjects(false);
      });

    return () => {
      isActive = false;
    };
  }, [profile]);

  const profileFields = useMemo(
    () => [
      { label: "Дата рождения", value: derivedProfile.dateOfBirthday },
      { label: "Должность", value: derivedProfile.role },
      { label: "TelegramID", value: derivedProfile.telegramID },
    ],
    [derivedProfile]
  );

  const normalizedUsername = profile?.username
    ? profile.username.replace(/^@/, "").toLowerCase()
    : "";

  return (
    <section className="profile-screen">
      <div className="profile-screen__inner">
        <aside className="profile-card">
          <div className="profile-card__avatar">
            <img
              src={derivedProfile.avatar}
              alt="Портрет пользователя"
              loading="lazy"
            />
          </div>
          <p className="profile-card__handle">{derivedProfile.handle}</p>
          <p className="profile-card__name">{derivedProfile.fullName}</p>
          <div className="profile-card__fields" aria-label="личные данные">
            {profileFields.map((field) => (
              <article className="profile-card__field" key={field.label}>
                <span className="profile-card__label">{field.label}:</span>
                <span className="profile-card__value">{field.value}</span>
              </article>
            ))}
          </div>
          <button type="button" className="profile-card__button">
            Изменить
          </button>
        </aside>

        <section className="profile-projects">
          <div className="profile-projects__header">
            <span className="profile-projects__badge">Мои проекты</span>
          </div>

          <div className="profile-projects__list">
            {isLoadingProjects && (
              <article className="project-card" aria-busy="true">
                <div className="project-card__main">
                  <p className="project-card__title">Загрузка...</p>
                  <p className="project-card__subtitle">подождите немного</p>
                </div>
              </article>
            )}
            {!isLoadingProjects && projectsError && (
              <article className="project-card">
                <div className="project-card__main">
                  <p className="project-card__title">Ошибка</p>
                  <p className="project-card__subtitle">{projectsError}</p>
                </div>
              </article>
            )}
            {!isLoadingProjects &&
              !projectsError &&
              projects.length === 0 && (
                <article className="project-card">
                  <div className="project-card__main">
                    <p className="project-card__title">Проекты не найдены</p>
                    <p className="project-card__subtitle">
                      нет связанных проектов для профиля
                    </p>
                  </div>
                </article>
              )}
            {!isLoadingProjects &&
              !projectsError &&
              projects.length > 0 &&
              projects.map((project) => {
                const members = Array.isArray(project.members)
                  ? project.members
                  : [];
                const currentMember = members.find(
                  (member) =>
                    member.username?.toLowerCase() === normalizedUsername
                );
                const role = currentMember?.role || "участник";
                return (
                  <article className="project-card" key={project.id}>
                    <div className="project-card__main">
                      <p className="project-card__title">{project.title}</p>
                      <Link
                        className="project-card__subtitle"
                        to={`/projects/${project.id}?from=profile`}
                      >
                        {TASKS_LINK_TEXT}
                      </Link>
                    </div>
                    <div className="project-card__meta">
                      <p>
                        <span>статус:</span> {project.status || "—"}
                      </p>
                      <p>
                        <span>участники:</span> {members.length}
                      </p>
                      <p>
                        <span>роль:</span> {role}
                      </p>
                    </div>
                  </article>
                );
              })}
          </div>
        </section>
      </div>
    </section>
  );
}
