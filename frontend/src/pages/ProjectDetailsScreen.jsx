import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import "../styles/ProjectDetailsScreen.scss";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";

export default function ProjectDetailsScreen() {
  const { id } = useParams();
  const [project, setProject] = useState(null);
  const [projectError, setProjectError] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let isActive = true;
    setIsLoading(true);
    setProjectError("");

    fetch(`${API_BASE}/projects/${id}`)
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
        setProject(data);
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setProjectError("Не удалось загрузить проект");
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
  }, [id]);

  const members = useMemo(() => {
    if (!project?.members) {
      return [];
    }
    return project.members.map((member) => ({
      name: member.full_name || member.username || "Участник",
      role: member.role || "",
    }));
  }, [project]);

  return (
    <section className="project-details">
      <header className="project-details__header">
        <Link className="project-details__back" to="/projects" aria-label="Назад">
          ←
        </Link>
        <div className="project-details__title">
          <p className="project-details__label">Просмотр проекта</p>
          <h2 className="project-details__name">
            {project?.title || "Проект"}
          </h2>
          <p className="project-details__description">
            {project?.description || "Описание пока не добавлено."}
          </p>
        </div>
        <button className="project-details__action" type="button">
          Управление
        </button>
      </header>

      {isLoading && (
        <div className="project-details__placeholder">
          Загрузка данных проекта...
        </div>
      )}
      {!isLoading && projectError && (
        <div className="project-details__placeholder">{projectError}</div>
      )}

      {!isLoading && !projectError && project && (
        <div className="project-details__grid">
          <section className="project-details__tasks">
            <h3>Задачи</h3>
            <div className="project-details__task">
              <p className="project-details__task-title">Задач пока нет</p>
              <p className="project-details__task-subtitle">
                Здесь появятся реальные задачи проекта.
              </p>
            </div>
          </section>

          <aside className="project-details__side">
            <div className="project-details__status">
              <p>Статус: {project.status || "—"}</p>
            </div>
            <div className="project-details__members">
              <h3>Участники</h3>
              <div className="project-details__members-list">
                {members.length === 0 && (
                  <div className="project-details__member">
                    Нет участников
                  </div>
                )}
                {members.map((member) => (
                  <div className="project-details__member" key={member.name}>
                    <span>{member.name}</span>
                    {member.role && <em>{member.role}</em>}
                  </div>
                ))}
              </div>
            </div>
          </aside>
        </div>
      )}
    </section>
  );
}
