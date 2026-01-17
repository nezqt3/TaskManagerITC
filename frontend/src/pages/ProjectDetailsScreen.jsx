import { useEffect, useLayoutEffect, useMemo, useRef, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";
import "../styles/ProjectDetailsScreen.scss";
import { getAuthHeaders, getProfile, isAdmin, isModerator, parseRoles } from "../utils/auth";
import { apiFetch } from "../utils/api";

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
  const [completionNotes, setCompletionNotes] = useState({});
  const [usersDirectory, setUsersDirectory] = useState({});
  const profile = useMemo(() => getProfile(), []);
  const [effectiveRole, setEffectiveRole] = useState(profile?.role || "");
  const [taskForm, setTaskForm] = useState({
    title: "",
    description: "",
    deadline: "",
    user: "",
  });
  const fromProfile = searchParams.get("from") === "profile";
  const highlightedTaskId = useMemo(() => {
    const value = Number(searchParams.get("taskId"));
    return Number.isNaN(value) ? null : value;
  }, [searchParams]);
  const backPath = fromProfile ? "/profile" : "/projects";
  const taskListRef = useRef(null);

  const authHeaders = useMemo(() => getAuthHeaders(), []);

  useEffect(() => {
    let isActive = true;
    setIsLoading(true);
    setProjectError("");

    apiFetch(`${API_BASE}/projects/${id}`, {
      headers: authHeaders,
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
  }, [authHeaders, id]);

  useLayoutEffect(() => {
    if (isLoadingTasks || !highlightedTaskId) {
      return;
    }
    const container = taskListRef.current;
    if (!container) {
      return;
    }
    const target = container.querySelector(
      `[data-task-id="${highlightedTaskId}"]`
    );
    if (!target) {
      return;
    }

    const scrollToTask = () => {
      const list = container.querySelector(".project-details__task-list");
      const scrollTargets = [];

      if (list && list.scrollHeight > list.clientHeight) {
        scrollTargets.push(list);
      }
      if (container.scrollHeight > container.clientHeight) {
        scrollTargets.push(container);
      }

      if (scrollTargets.length === 0) {
        target.scrollIntoView({ block: "center", behavior: "smooth" });
        return true;
      }

      scrollTargets.forEach((scrollTarget) => {
        const parentRect = scrollTarget.getBoundingClientRect();
        const targetRect = target.getBoundingClientRect();
        const offset = targetRect.top - parentRect.top - 16;
        scrollTarget.scrollTo({
          top: scrollTarget.scrollTop + offset,
          behavior: "smooth",
        });
      });
      return true;
    };

    const raf = requestAnimationFrame(scrollToTask);
    const timer = setTimeout(scrollToTask, 150);

    return () => {
      cancelAnimationFrame(raf);
      clearTimeout(timer);
    };
  }, [isLoadingTasks, highlightedTaskId, tasks]);

  useEffect(() => {
    let isActive = true;
    setIsLoadingTasks(true);
    setTasksError("");

    apiFetch(`${API_BASE}/tasks?id_project=${id}`, {
      headers: authHeaders,
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
  }, [authHeaders, id]);

  useEffect(() => {
    let isActive = true;
    apiFetch(`${API_BASE}/get_users`, { headers: authHeaders })
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
        const users = Array.isArray(data) ? data : [];
        const directory = {};
        users.forEach((user) => {
          if (!user.username) {
            return;
          }
          directory[user.username.toLowerCase()] = {
            telegramId: user.telegram_id,
            fullName:
              user.full_name ||
              [user.first_name, user.last_name].filter(Boolean).join(" "),
          };
        });
        const profileId = profile?.telegram_id;
        const profileUsername = profile?.username?.replace(/^@/, "").toLowerCase();
        const matchedUser = users.find(
          (user) =>
            (profileId && user.telegram_id === profileId) ||
            (profileUsername &&
              user.username?.toLowerCase() === profileUsername)
        );
        if (matchedUser?.role) {
          setEffectiveRole(matchedUser.role);
        }
        setUsersDirectory(directory);
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setUsersDirectory({});
      });

    return () => {
      isActive = false;
    };
  }, [authHeaders, profile]);

  const members = useMemo(() => {
    if (!project?.members) {
      return [];
    }
    return project.members.map((member) => ({
      name: member.full_name || member.username || "Участник",
      role: member.role || "",
      username: member.username || "",
    }));
  }, [project]);

  const isProjectMember = useMemo(() => {
    const profileUsername = profile?.username
      ? profile.username.replace(/^@/, "").toLowerCase()
      : "";
    if (!profileUsername || !project?.members) {
      return false;
    }
    return project.members.some(
      (member) =>
        (member.username || "").replace(/^@/, "").toLowerCase() ===
        profileUsername
    );
  }, [profile, project]);

  const memberRole = useMemo(() => {
    const profileUsername = profile?.username
      ? profile.username.replace(/^@/, "").toLowerCase()
      : "";
    if (!profileUsername || !project?.members) {
      return "";
    }
    const member = project.members.find(
      (entry) =>
        (entry.username || "").replace(/^@/, "").toLowerCase() ===
        profileUsername
    );
    return member?.role || "";
  }, [profile, project]);

  const isProjectLeader = useMemo(() => {
    const roles = parseRoles(memberRole);
    return roles.includes("руководитель");
  }, [memberRole]);

  const canEditTasks =
    isAdmin(effectiveRole) ||
    (isProjectMember && (isModerator(effectiveRole) || isProjectLeader));

  const taskAuthor = useMemo(() => {
    const fullName =
      profile?.full_name ||
      [profile?.first_name, profile?.last_name].filter(Boolean).join(" ");

    return fullName || profile?.username || "";
  }, [profile]);

  const canSubmitCompletion = () => Boolean(profile);

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

    const assignedUser = taskForm.user.trim();
    const directoryEntry = assignedUser
      ? usersDirectory[assignedUser.replace(/^@/, "").toLowerCase()]
      : null;

    const payload = {
      title: trimmedTitle,
      description: taskForm.description.trim(),
      deadline: taskForm.deadline,
      status: "Новая",
      user: assignedUser,
      author: taskAuthor,
      id_project: Number(id),
      id_user: directoryEntry?.telegramId ? Number(directoryEntry.telegramId) : 0,
    };

    apiFetch(`${API_BASE}/tasks`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
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
        setTasks((prev) => [data, ...prev]);
        setTaskForm({
          title: "",
          description: "",
          deadline: "",
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

  const handleCompletionSubmit = (taskId) => {
    const message = (completionNotes[taskId] || "").trim();
    apiFetch(`${API_BASE}/tasks/${taskId}/complete`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
      },
      body: JSON.stringify({ message }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return response.text();
      })
      .then(() => {
        setTasks((prev) =>
          prev.map((task) =>
            task.id === taskId
              ? {
                  ...task,
                  status: "На проверке",
                  completion_message: message,
                }
              : task
          )
        );
        setCompletionNotes((prev) => ({ ...prev, [taskId]: "" }));
      })
      .catch(() => {
        setTasksError("Не удалось отправить отчет по задаче");
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
            {canEditTasks && (
              <Link
                className="project-details__action"
                to={`/projects/${id}/manage`}
              >
                Управление
              </Link>
            )}
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
            <section className="project-details__tasks" ref={taskListRef}>
              <div className="project-details__tasks-header">
                <h3>Задачи</h3>
                <span className="project-details__tasks-count">
                  {tasks.length}
                </span>
                {canEditTasks && (
                  <button
                    className="project-details__tasks-toggle"
                    type="button"
                    onClick={() => setShowTaskForm((prev) => !prev)}
                    aria-expanded={showTaskForm}
                  >
                    {showTaskForm ? "—" : "+"}
                  </button>
                )}
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
                    <p className="project-details__task-title">{tasksError}</p>
                  </div>
                )}
                {!isLoadingTasks && !tasksError && tasks.length === 0 && (
                  <div className="project-details__task">
                    <p className="project-details__task-title">Задач пока нет</p>
                    <p className="project-details__task-subtitle">
                      Здесь появятся реальные задачи проекта.
                    </p>
                  </div>
                )}
                {!isLoadingTasks &&
                  !tasksError &&
                  tasks.map((task) => (
                    <article
                      className={`project-details__task-card${
                        highlightedTaskId === task.id
                          ? " project-details__task-card--highlight"
                          : ""
                      }`}
                      key={`${task.id}-${task.title}`}
                      data-task-id={task.id}
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
                        {task.deadline && <span>Срок: {task.deadline}</span>}
                        {task.user && <span>Исп.: {task.user}</span>}
                        {task.author && <span>Автор: {task.author}</span>}
                      </div>
                      {task.completion_message && (
                        <div className="project-details__task-note">
                          <strong>Отчет:</strong> {task.completion_message}
                        </div>
                      )}
                      {task.review_message && (
                        <div className="project-details__task-note">
                          <strong>Ответ модератора:</strong> {task.review_message}
                        </div>
                      )}
                      {task.status !== "На проверке" &&
                        task.status !== "Выполнена" &&
                        canSubmitCompletion(task) && (
                          <div className="project-details__task-completion">
                            <textarea
                              placeholder="Сообщение о выполненной работе"
                              rows={2}
                              value={completionNotes[task.id] || ""}
                              onChange={(event) =>
                                setCompletionNotes((prev) => ({
                                  ...prev,
                                  [task.id]: event.target.value,
                                }))
                              }
                            />
                            <button
                              type="button"
                              onClick={() => handleCompletionSubmit(task.id)}
                            >
                              Отправить на проверку
                            </button>
                          </div>
                        )}
                    </article>
                  ))}
              </div>
              {showTaskForm && canEditTasks && (
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
                        list="project-members"
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
                  <datalist id="project-members">
                    {members.map((member) => (
                      <option
                        key={member.username || member.name}
                        value={member.username}
                      >
                        {member.name}
                      </option>
                    ))}
                  </datalist>
                  {taskFormError && (
                    <p className="project-details__form-error">{taskFormError}</p>
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
                <div className="project-details__member">{projectError}</div>
              )}
              {!isLoading && !projectError && members.length === 0 && (
                <div className="project-details__member">Нет участников</div>
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
