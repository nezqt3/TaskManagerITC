import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import "../styles/ProjectsScreen.scss";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";
const TASKS_LINK_TEXT = "просмотреть задачу";

export default function ProjectsScreen() {
  const [projects, setProjects] = useState([]);
  const [projectsError, setProjectsError] = useState("");
  const [isLoadingProjects, setIsLoadingProjects] = useState(true);

  useEffect(() => {
    let isActive = true;
    setIsLoadingProjects(true);
    setProjectsError("");

    fetch(`${API_BASE}/projects`)
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
  }, []);

  return (
    <section className="projects-screen">
      <header className="projects-screen__header">
        <h2 className="projects-screen__title">Все проекты</h2>
      </header>
      <div className="projects-screen__list">
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
                  пока нет загруженных карточек
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
            return (
              <article className="project-card" key={project.id}>
                <div className="project-card__main">
                  <p className="project-card__title">{project.title}</p>
                  <Link
                    className="project-card__subtitle"
                    to={`/projects/${project.id}`}
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
                </div>
              </article>
            );
          })}
      </div>
    </section>
  );
}
