import { useEffect, useMemo, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";
import "../styles/ProjectDetailsScreen.scss";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";

export default function ProjectDetailsScreen() {
  const { id } = useParams();
  const [searchParams] = useSearchParams();
  const [project, setProject] = useState(null);
  const [projectError, setProjectError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [tasks, setTasks] = useState([]);
  const [tasksError, setTasksError] = useState("");
  const [isLoadingTasks, setIsLoadingTasks] = useState(true);
  const [isSavingTask, setIsSavingTask] = useState(false);
  const [taskFormError, setTaskFormError] = useState("");
  const [showTaskForm, setShowTaskForm] = useState(false);
  const [profile] = useState(() => {
    const stored = localStorage.getItem("profile");
    if (!stored) {
      return null;
    }

    try {
      return JSON.parse(stored);
    } catch (error) {
      return null;
    }
  });
  const [taskForm, setTaskForm] = useState({
    title: "",
    description: "",
    deadline: "",
    status: "Новая",
    user: "",
  });
  const fromProfile = searchParams.get("from") === "profile";
  const backPath = fromProfile ? "/profile" : "/projects";

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

  useEffect(() => {
    let isActive = true;
    setIsLoadingTasks(true);
    setTasksError("");

    fetch(`${API_BASE}/tasks?id_project=${id}`)
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
        setTasks(Array.isArray(data) ? data : []);
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setTasksError("Не удалось загрузить задачи");
      })
      .finally(() => {
        if (!isActive) {
          return;
        }
        setIsLoadingTasks(false);
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

  const taskAuthor = useMemo(() => {
    const fullName =
      profile?.full_name ||
      [profile?.first_name, profile?.last_name].filter(Boolean).join(" ");

    return fullName || profile?.username || "";
  }, [profile]);

  const handleTaskChange = (event) => {
    const { name, value } = event.target;
    setTaskForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleTaskSubmit = (event) => {
    event.preventDefault();
    const trimmedTitle = taskForm.title.trim();
    if (!trimmedTitle) {
      setTaskFormError("Добавьте название задачи");
      return;
    }

    setIsSavingTask(true);
    setTaskFormError("");

    const payload = {
      title: trimmedTitle,
      description: taskForm.description.trim(),
      deadline: taskForm.deadline,
      status: taskForm.status || "Новая",
      user: taskForm.user.trim(),
      author: taskAuthor,
      id_project: Number(id),
    };

    fetch(`${API_BASE}/tasks`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return response.json();
      })
      .then((data) => {
        setTasks((prev) => [...prev, data]);
        setTaskForm({
          title: "",
          description: "",
          deadline: "",
          status: "Новая",
          user: "",
        });
        setShowTaskForm(false);
      })
      .catch(() => {
        setTaskFormError("Не удалось добавить задачу");
      })
      .finally(() => {
        setIsSavingTask(false);
      });
  };

  return (
    <section className="project-details">
      <div className="project-details__layout">
        <div className="project-details__content">
          <header className="project-details__header">
            <Link
              className="project-details__back"
              to={backPath}
              aria-label="Назад"
            >
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
            <section className="project-details__tasks">
              <div className="project-details__tasks-header">
                <h3>Задачи</h3>
                <span className="project-details__tasks-count">
                  {tasks.length}
                </span>
                <button
                  className="project-details__tasks-toggle"
                  type="button"
                  onClick={() => setShowTaskForm((prev) => !prev)}
                  aria-expanded={showTaskForm}
                >
                  {showTaskForm ? "—" : "+"}
                </button>
              </div>
              <div className="project-details__task-list">
                {isLoadingTasks && (
                  <div className="project-details__task">
                    <p className="project-details__task-title">
                      Загрузка задач...
                    </p>
                  </div>
                )}
                {!isLoadingTasks && tasksError && (
                  <div className="project-details__task">
                    <p className="project-details__task-title">
                      {tasksError}
                    </p>
                  </div>
                )}
                {!isLoadingTasks &&
                  !tasksError &&
                  tasks.length === 0 && (
                    <div className="project-details__task">
                      <p className="project-details__task-title">
                        Задач пока нет
                      </p>
                      <p className="project-details__task-subtitle">
                        Здесь появятся реальные задачи проекта.
                      </p>
                    </div>
                  )}
                {!isLoadingTasks &&
                  !tasksError &&
                  tasks.map((task) => (
                    <article
                      className="project-details__task-card"
                      key={`${task.id}-${task.title}`}
                    >
                      <div className="project-details__task-main">
                        <p className="project-details__task-title">
                          {task.title || "Без названия"}
                        </p>
                        {task.description && (
                          <p className="project-details__task-subtitle">
                            {task.description}
                          </p>
                        )}
                      </div>
                      <div className="project-details__task-meta">
                        <span>{task.status || "Без статуса"}</span>
                        {task.deadline && (
                          <span>Срок: {task.deadline}</span>
                        )}
                        {task.user && <span>Исп.: {task.user}</span>}
                        {task.author && <span>Автор: {task.author}</span>}
                      </div>
                    </article>
                  ))}
              </div>
              {showTaskForm && (
                <form
                  className="project-details__task-form"
                  onSubmit={handleTaskSubmit}
                >
                  <h4>Добавить задачу</h4>
                  <div className="project-details__form-grid">
                    <label className="project-details__field">
                      Название
                      <input
                        name="title"
                        value={taskForm.title}
                        onChange={handleTaskChange}
                        placeholder="Кратко о задаче"
                        required
                      />
                    </label>
                    <label className="project-details__field">
                      Ответственный
                      <input
                        name="user"
                        value={taskForm.user}
                        onChange={handleTaskChange}
                        placeholder="username"
                      />
                    </label>
                    <label className="project-details__field">
                      Дедлайн
                      <input
                        type="date"
                        name="deadline"
                        value={taskForm.deadline}
                        onChange={handleTaskChange}
                      />
                    </label>
                    <label className="project-details__field">
                      Статус
                      <select
                        name="status"
                        value={taskForm.status}
                        onChange={handleTaskChange}
                      >
                        <option value="Новая">Новая</option>
                        <option value="В работе">В работе</option>
                        <option value="Готово">Готово</option>
                      </select>
                    </label>
                    <label className="project-details__field project-details__field--full">
                      Описание
                      <textarea
                        name="description"
                        value={taskForm.description}
                        onChange={handleTaskChange}
                        placeholder="Что нужно сделать"
                        rows={3}
                      />
                    </label>
                  </div>
                  {taskFormError && (
                    <p className="project-details__form-error">
                      {taskFormError}
                    </p>
                  )}
                  <div className="project-details__form-actions">
                    <button
                      className="project-details__submit"
                      type="submit"
                      disabled={isSavingTask}
                    >
                      {isSavingTask ? "Сохранение..." : "Добавить"}
                    </button>
                  </div>
                </form>
              )}
            </section>
          )}
        </div>

        <aside className="project-details__side">
          <div className="project-details__status">
            <p>Статус: {project?.status || "—"}</p>
          </div>
          <div className="project-details__members">
            <h3>Участники</h3>
            <div className="project-details__members-list">
              {isLoading && (
                <div className="project-details__member">
                  Загрузка участников...
                </div>
              )}
              {!isLoading && projectError && (
                <div className="project-details__member">
                  {projectError}
                </div>
              )}
              {!isLoading &&
                !projectError &&
                members.length === 0 && (
                  <div className="project-details__member">
                    Нет участников
                  </div>
                )}
              {!isLoading &&
                !projectError &&
                members.map((member) => (
                  <div className="project-details__member" key={member.name}>
                    <span>{member.name}</span>
                    {member.role && <em>{member.role}</em>}
                  </div>
                ))}
            </div>
          </div>
        </aside>
      </div>
    </section>
  );
}
