import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import "../styles/ProjectManagementScreen.scss";
import {
  getAuthHeaders,
  getProfile,
  isAdmin,
  isModerator,
  parseRoles,
} from "../utils/auth";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";
const TASK_STATUSES = [
  "Новая",
  "В работе",
  "На проверке",
  "Выполнена",
  "Отклонена",
];
const BASE_ROLES = ["Руководитель", "Разработчик", "Участник"];

function normalizeUsername(value) {
  return value.replace(/^@/, "").trim().toLowerCase();
}

function parseMemberRole(role) {
  const roles = parseRoles(role);
  const isAdminRole = roles.includes("админ") || roles.includes("admin");
  const isModeratorRole = roles.includes("модератор") || roles.includes("moderator");
  const baseRole = BASE_ROLES.find((item) => roles.includes(item.toLowerCase())) ||
    role.split(",")[0]?.trim() ||
    "Участник";

  return {
    baseRole,
    isAdmin: isAdminRole,
    isModerator: isModeratorRole || baseRole.toLowerCase() === "руководитель",
  };
}

function buildRoleString(baseRole, moderator, admin) {
  const roles = [baseRole];
  if (baseRole.toLowerCase() === "руководитель") {
    roles.push("Модератор");
  } else if (moderator) {
    roles.push("Модератор");
  }
  if (admin && baseRole.toLowerCase() !== "разработчик") {
    roles.push("Админ");
  }
  return roles.join(", ");
}

export default function ProjectManagementScreen() {
  const { id } = useParams();
  const profile = useMemo(() => getProfile(), []);
  const [effectiveRole, setEffectiveRole] = useState(profile?.role || "");
  const isMemberAdmin = isAdmin(effectiveRole);

  const [project, setProject] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [actionError, setActionError] = useState("");
  const [memberForm, setMemberForm] = useState({
    username: "",
    baseRole: "Участник",
    moderator: false,
    admin: false,
  });
  const [memberEdits, setMemberEdits] = useState({});
  const [taskEdits, setTaskEdits] = useState({});
  const [reviewNotes, setReviewNotes] = useState({});

  const authHeaders = useMemo(() => getAuthHeaders(), []);

  useEffect(() => {
    let isActive = true;
    setIsLoading(true);
    setError("");

    Promise.all([
      fetch(`${API_BASE}/projects/${id}`, { headers: authHeaders }),
      fetch(`${API_BASE}/tasks?id_project=${id}`, { headers: authHeaders }),
      fetch(`${API_BASE}/get_users`, { headers: authHeaders }),
    ])
      .then(async ([projectResponse, tasksResponse, usersResponse]) => {
        if (!projectResponse.ok || !tasksResponse.ok || !usersResponse.ok) {
          throw new Error("failed");
        }
        const projectData = await projectResponse.json();
        const tasksData = await tasksResponse.json();
        const usersData = await usersResponse.json();
        return { projectData, tasksData, usersData };
      })
      .then(({ projectData, tasksData, usersData }) => {
        if (!isActive) {
          return;
        }
        setProject(projectData);
        setTasks(Array.isArray(tasksData) ? tasksData : []);
        const list = Array.isArray(usersData) ? usersData : [];
        setUsers(list);
        const profileId = profile?.telegram_id;
        const profileUsername = profile?.username
          ? normalizeUsername(profile.username)
          : "";
        const matchedUser = list.find(
          (user) =>
            (profileId && user.telegram_id === profileId) ||
            (profileUsername &&
              normalizeUsername(user.username || "") === profileUsername)
        );
        if (matchedUser?.role) {
          setEffectiveRole(matchedUser.role);
        }
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setError("Не удалось загрузить данные проекта");
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
  }, [authHeaders, id, profile?.telegram_id, profile?.username]);

  useEffect(() => {
    if (!project?.members) {
      return;
    }
    const nextEdits = {};
    project.members.forEach((member) => {
      const parsed = parseMemberRole(member.role || "");
      nextEdits[normalizeUsername(member.username || "")] = {
        baseRole: parsed.baseRole,
        moderator: parsed.isModerator,
        admin: parsed.isAdmin,
      };
    });
    setMemberEdits(nextEdits);
  }, [project]);

  useEffect(() => {
    if (!tasks.length) {
      return;
    }
    const nextEdits = {};
    tasks.forEach((task) => {
      nextEdits[task.id] = {
        title: task.title || "",
        description: task.description || "",
        deadline: task.deadline || "",
        status: task.status || "Новая",
        user: task.user || "",
        id_user: task.id_user || 0,
        id_project: task.id_project,
      };
    });
    setTaskEdits(nextEdits);
  }, [tasks]);

  const userOptions = useMemo(() => {
    return users
      .filter((user) => user.username)
      .map((user) => ({
        username: user.username,
        fullName:
          user.full_name ||
          [user.first_name, user.last_name].filter(Boolean).join(" "),
        telegramId: user.telegram_id,
      }));
  }, [users]);

  const isProjectMember = useMemo(() => {
    if (isMemberAdmin) {
      return true;
    }
    const profileUsername = profile?.username
      ? normalizeUsername(profile.username)
      : "";
    if (!profileUsername || !project?.members) {
      return false;
    }
    return project.members.some(
      (member) =>
        normalizeUsername(member.username || "") === profileUsername
    );
  }, [isMemberAdmin, profile, project]);

  const memberRole = useMemo(() => {
    const profileUsername = profile?.username
      ? normalizeUsername(profile.username)
      : "";
    if (!profileUsername || !project?.members) {
      return "";
    }
    const member = project.members.find(
      (entry) => normalizeUsername(entry.username || "") === profileUsername
    );
    return member?.role || "";
  }, [profile, project]);

  const isProjectLeader = useMemo(() => {
    const roles = parseRoles(memberRole);
    return roles.includes("руководитель");
  }, [memberRole]);

  const canManageProjectMembers = isMemberAdmin;
  const canManageProjectTasks =
    isMemberAdmin ||
    (isProjectMember && (isModerator(effectiveRole) || isProjectLeader));
  const canReviewTasks =
    isMemberAdmin ||
    (isProjectMember && (isModerator(effectiveRole) || isProjectLeader));

  const userDirectory = useMemo(() => {
    const directory = {};
    userOptions.forEach((user) => {
      directory[normalizeUsername(user.username)] = user.telegramId;
    });
    return directory;
  }, [userOptions]);

  const handleMemberFormChange = (event) => {
    const { name, value, type, checked } = event.target;
    setMemberForm((prev) => {
      const next = { ...prev, [name]: type === "checkbox" ? checked : value };
      if (name === "baseRole" && value.toLowerCase() === "разработчик") {
        next.admin = false;
      }
      if (name === "baseRole" && value.toLowerCase() === "руководитель") {
        next.moderator = true;
      }
      return next;
    });
  };

  const handleMemberEditChange = (username, field, value) => {
    setMemberEdits((prev) => {
      const current = prev[username] || {
        baseRole: "Участник",
        moderator: false,
        admin: false,
      };
      const next = { ...current, [field]: value };
      if (field === "baseRole" && value.toLowerCase() === "разработчик") {
        next.admin = false;
      }
      if (field === "baseRole" && value.toLowerCase() === "руководитель") {
        next.moderator = true;
      }
      return { ...prev, [username]: next };
    });
  };

  const handleAddMember = (event) => {
    event.preventDefault();
    if (!canManageProjectMembers) {
      return;
    }
    setActionError("");

    const normalized = normalizeUsername(memberForm.username);
    if (!normalized) {
      setActionError("Введите username участника");
      return;
    }

    const roleValue = buildRoleString(
      memberForm.baseRole,
      memberForm.moderator,
      memberForm.admin
    );
    const userInfo = userOptions.find(
      (user) => normalizeUsername(user.username || "") === normalized
    );

    fetch(`${API_BASE}/projects/${id}/members`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
      },
      body: JSON.stringify({
        username: normalized,
        full_name: userInfo?.fullName || "",
        role: roleValue,
        telegram_id: userInfo?.telegramId || "",
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return response.text();
      })
      .then(() =>
        fetch(`${API_BASE}/projects/${id}`, { headers: authHeaders })
      )
      .then((response) => response.json())
      .then((data) => {
        setProject(data);
        setMemberForm({
          username: "",
          baseRole: "Участник",
          moderator: false,
          admin: false,
        });
      })
      .catch(() => setActionError("Не удалось добавить участника"));
  };

  const handleUpdateMemberRole = (username) => {
    if (!canManageProjectMembers) {
      return;
    }
    setActionError("");
    const edit = memberEdits[username];
    if (!edit) {
      return;
    }

    const roleValue = buildRoleString(edit.baseRole, edit.moderator, edit.admin);
    fetch(`${API_BASE}/projects/${id}/members/${username}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
      },
      body: JSON.stringify({ role: roleValue }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return fetch(`${API_BASE}/projects/${id}`, { headers: authHeaders });
      })
      .then((response) => response.json())
      .then((data) => setProject(data))
      .catch(() => setActionError("Не удалось обновить роль"));
  };

  const handleRemoveMember = (username) => {
    if (!canManageProjectMembers) {
      return;
    }
    setActionError("");
    fetch(`${API_BASE}/projects/${id}/members/${username}`, {
      method: "DELETE",
      headers: authHeaders,
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return response.text();
      })
      .then(() =>
        fetch(`${API_BASE}/projects/${id}`, { headers: authHeaders })
      )
      .then((response) => response.json())
      .then((data) => setProject(data))
      .catch(() => setActionError("Не удалось удалить участника"));
  };

  const handleTaskEditChange = (taskId, field, value) => {
    setTaskEdits((prev) => ({
      ...prev,
      [taskId]: {
        ...prev[taskId],
        [field]: value,
      },
    }));
  };

  const handleTaskUpdate = (taskId) => {
    if (!canManageProjectTasks) {
      return;
    }
    setActionError("");
    const payload = taskEdits[taskId];
    if (!payload) {
      return;
    }
    const normalized = normalizeUsername(payload.user || "");
    const idUser = normalized ? Number(userDirectory[normalized] || 0) : 0;

    fetch(`${API_BASE}/tasks/${taskId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
      },
      body: JSON.stringify({ ...payload, id_user: idUser }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
      })
      .catch(() => setActionError("Не удалось сохранить задачу"));
  };

  const handleTaskDelete = (taskId) => {
    if (!canManageProjectTasks) {
      return;
    }
    setActionError("");
    fetch(`${API_BASE}/tasks/${taskId}`, {
      method: "DELETE",
      headers: authHeaders,
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
      })
      .then(() =>
        fetch(`${API_BASE}/tasks?id_project=${id}`, { headers: authHeaders })
      )
      .then((response) => response.json())
      .then((data) => setTasks(Array.isArray(data) ? data : []))
      .catch(() => setActionError("Не удалось удалить задачу"));
  };

  const handleReviewTask = (taskId, approved) => {
    if (!canReviewTasks) {
      return;
    }
    setActionError("");
    const message = reviewNotes[taskId] || "";

    fetch(`${API_BASE}/tasks/${taskId}/review`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...authHeaders,
      },
      body: JSON.stringify({ approved, message }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
        return fetch(`${API_BASE}/tasks?id_project=${id}`, { headers: authHeaders });
      })
      .then((response) => response.json())
      .then((data) => setTasks(Array.isArray(data) ? data : []))
      .catch(() => setActionError("Не удалось подтвердить задачу"));
  };

  if (isLoading) {
    return (
      <section className="project-management">
        <div className="project-management__panel">Загрузка проекта...</div>
      </section>
    );
  }

  if (error) {
    return (
      <section className="project-management">
        <div className="project-management__panel">{error}</div>
      </section>
    );
  }

  return (
    <section className="project-management">
      <header className="project-management__header">
        <Link className="project-management__back" to={`/projects/${id}`}>
          ←
        </Link>
        <div>
          <p className="project-management__label">Управление проектом</p>
          <h2 className="project-management__title">
            {project?.title || "Проект"}
          </h2>
          <p className="project-management__subtitle">
            {project?.description || "Описание пока не добавлено."}
          </p>
        </div>
        <div className="project-management__role">
          <span>Уровень:</span>
          <strong>
            {isMemberAdmin
              ? "Админ"
              : canManageProjectTasks
                ? "Модератор"
                : "Участник"}
          </strong>
        </div>
      </header>

      {actionError && (
        <div className="project-management__error">{actionError}</div>
      )}

      <div className="project-management__layout">
        <section className="project-management__panel">
          <div className="project-management__panel-header">
            <h3>Участники</h3>
            <span>{project?.members?.length || 0}</span>
          </div>
          {!canManageProjectMembers && (
            <p className="project-management__hint">
              Управлять участниками может только администратор проекта.
            </p>
          )}
          <div className="project-management__members">
            {project?.members?.map((member) => {
              const key = normalizeUsername(member.username || "");
              const edit = memberEdits[key] || {
                baseRole: "Участник",
                moderator: false,
                admin: false,
              };
              const fullName = member.full_name || member.username || "Участник";
              const displayUsername = (member.username || "").replace(/^@/, "");
              return (
                <article className="member-row" key={key}>
                  <div className="member-row__info">
                    <p>{fullName}</p>
                    <span>@{displayUsername || "—"}</span>
                  </div>
                  <div className="member-row__controls">
                    <select
                      value={edit.baseRole}
                      onChange={(event) =>
                        handleMemberEditChange(key, "baseRole", event.target.value)
                      }
                      disabled={!canManageProjectMembers}
                    >
                      {BASE_ROLES.map((roleOption) => (
                        <option key={roleOption} value={roleOption}>
                          {roleOption}
                        </option>
                      ))}
                    </select>
                    <label>
                      <input
                        type="checkbox"
                        checked={edit.moderator}
                        onChange={(event) =>
                          handleMemberEditChange(key, "moderator", event.target.checked)
                        }
                        disabled={
                          !canManageProjectMembers || edit.baseRole.toLowerCase() === "руководитель"
                        }
                      />
                      Модератор
                    </label>
                    <label>
                      <input
                        type="checkbox"
                        checked={edit.admin}
                        onChange={(event) =>
                          handleMemberEditChange(key, "admin", event.target.checked)
                        }
                        disabled={
                          !canManageProjectMembers || edit.baseRole.toLowerCase() === "разработчик"
                        }
                      />
                      Админ
                    </label>
                    {canManageProjectMembers && (
                      <div className="member-row__actions">
                        <button
                          type="button"
                          onClick={() => handleUpdateMemberRole(key)}
                        >
                          Сохранить
                        </button>
                        <button
                          type="button"
                          className="danger"
                          onClick={() => handleRemoveMember(key)}
                        >
                          Удалить
                        </button>
                      </div>
                    )}
                  </div>
                </article>
              );
            })}
          </div>

          <form className="member-form" onSubmit={handleAddMember}>
            <h4>Добавить участника</h4>
            <div className="member-form__grid">
              <label>
                Username
                <input
                  list="users-list"
                  name="username"
                  value={memberForm.username}
                  onChange={handleMemberFormChange}
                  placeholder="@username"
                  disabled={!canManageProjectMembers}
                />
              </label>
              <label>
                Базовая роль
                <select
                  name="baseRole"
                  value={memberForm.baseRole}
                  onChange={handleMemberFormChange}
                  disabled={!canManageProjectMembers}
                >
                  {BASE_ROLES.map((roleOption) => (
                    <option key={roleOption} value={roleOption}>
                      {roleOption}
                    </option>
                  ))}
                </select>
              </label>
              <label className="member-form__toggle">
                <input
                  type="checkbox"
                  name="moderator"
                  checked={memberForm.moderator}
                  onChange={handleMemberFormChange}
                  disabled={!canManageProjectMembers}
                />
                Модератор
              </label>
              <label className="member-form__toggle">
                <input
                  type="checkbox"
                  name="admin"
                  checked={memberForm.admin}
                  onChange={handleMemberFormChange}
                  disabled={!canManageProjectMembers || memberForm.baseRole.toLowerCase() === "разработчик"}
                />
                Админ
              </label>
            </div>
            <button type="submit" disabled={!canManageProjectMembers}>
              Добавить
            </button>
          </form>
          <datalist id="users-list">
            {userOptions.map((user) => (
              <option key={user.username} value={user.username}>
                {user.fullName || user.username}
              </option>
            ))}
          </datalist>
        </section>

        <section className="project-management__panel">
          <div className="project-management__panel-header">
            <h3>Задачи</h3>
            <span>{tasks.length}</span>
          </div>
          {!canManageProjectTasks && (
            <p className="project-management__hint">
              Управлять задачами могут модераторы и администраторы.
            </p>
          )}
          <div className="project-management__tasks">
            {tasks.length === 0 && (
              <p className="project-management__hint">Задачи пока не добавлены.</p>
            )}
            {tasks.map((task) => {
              const edit = taskEdits[task.id];
              if (!edit) {
                return null;
              }
              return (
                <article className="task-row" key={task.id}>
                  <div className="task-row__main">
                    <input
                      value={edit.title}
                      onChange={(event) =>
                        handleTaskEditChange(task.id, "title", event.target.value)
                      }
                      disabled={!canManageProjectTasks}
                    />
                    <textarea
                      value={edit.description}
                      onChange={(event) =>
                        handleTaskEditChange(task.id, "description", event.target.value)
                      }
                      rows={2}
                      disabled={!canManageProjectTasks}
                    />
                    <div className="task-row__meta">
                      <label>
                        Статус
                        <select
                          value={edit.status}
                          onChange={(event) =>
                            handleTaskEditChange(task.id, "status", event.target.value)
                          }
                          disabled={!canManageProjectTasks}
                        >
                          {TASK_STATUSES.map((status) => (
                            <option key={status} value={status}>
                              {status}
                            </option>
                          ))}
                        </select>
                      </label>
                      <label>
                        Дедлайн
                        <input
                          type="date"
                          value={edit.deadline}
                          onChange={(event) =>
                            handleTaskEditChange(task.id, "deadline", event.target.value)
                          }
                          disabled={!canManageProjectTasks}
                        />
                      </label>
                      <label>
                        Исполнитель
                        <input
                          value={edit.user}
                          onChange={(event) =>
                            handleTaskEditChange(task.id, "user", event.target.value)
                          }
                          disabled={!canManageProjectTasks}
                        />
                      </label>
                    </div>
                    {task.completion_message && (
                      <div className="task-row__note">
                        <strong>Отчет:</strong> {task.completion_message}
                      </div>
                    )}
                    {task.review_message && (
                      <div className="task-row__note">
                        <strong>Ответ модератора:</strong> {task.review_message}
                      </div>
                    )}
                  </div>
                  <div className="task-row__actions">
                    {task.status === "На проверке" && canReviewTasks && (
                      <div className="task-row__review">
                        <textarea
                          placeholder="Комментарий по проверке"
                          value={reviewNotes[task.id] || ""}
                          onChange={(event) =>
                            setReviewNotes((prev) => ({
                              ...prev,
                              [task.id]: event.target.value,
                            }))
                          }
                          rows={2}
                        />
                        <div className="task-row__review-actions">
                          <button
                            type="button"
                            onClick={() => handleReviewTask(task.id, true)}
                          >
                            Подтвердить
                          </button>
                          <button
                            type="button"
                            className="danger"
                            onClick={() => handleReviewTask(task.id, false)}
                          >
                            Отклонить
                          </button>
                        </div>
                      </div>
                    )}
                    {canManageProjectTasks && (
                      <div className="task-row__manage">
                        <button
                          type="button"
                          onClick={() => handleTaskUpdate(task.id)}
                        >
                          Сохранить
                        </button>
                        <button
                          type="button"
                          className="danger"
                          onClick={() => handleTaskDelete(task.id)}
                        >
                          Удалить
                        </button>
                      </div>
                    )}
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
