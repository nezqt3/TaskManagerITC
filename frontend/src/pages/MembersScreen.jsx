import { useEffect, useMemo, useState } from "react";
import "../styles/MembersScreen.scss";
import { getAuthHeaders, getProfile, isAdmin } from "../utils/auth";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";

export default function MembersScreen() {
  const profile = useMemo(() => getProfile(), []);
  const [effectiveRole, setEffectiveRole] = useState(profile?.role || "");
  const canEdit = isAdmin(effectiveRole);
  const [users, setUsers] = useState([]);
  const [usersError, setUsersError] = useState("");
  const [isLoadingUsers, setIsLoadingUsers] = useState(true);
  const [editingId, setEditingId] = useState("");
  const [editForm, setEditForm] = useState({});

  useEffect(() => {
    let isActive = true;
    setIsLoadingUsers(true);
    setUsersError("");

    fetch(`${API_BASE}/get_users`, {
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
        const list = Array.isArray(data) ? data : [];
        setUsers(list);
        const profileId = profile?.telegram_id;
        const profileUsername = profile?.username?.replace(/^@/, "").toLowerCase();
        const matchedUser = list.find(
          (user) =>
            (profileId && user.telegram_id === profileId) ||
            (profileUsername && user.username?.toLowerCase() === profileUsername)
        );
        if (matchedUser?.role) {
          setEffectiveRole(matchedUser.role);
        }
      })
      .catch(() => {
        if (!isActive) {
          return;
        }
        setUsersError("Не удалось загрузить пользователей");
      })
      .finally(() => {
        if (!isActive) {
          return;
        }
        setIsLoadingUsers(false);
      });

    return () => {
      isActive = false;
    };
  }, [profile]);

  useEffect(() => {
    if (!canEdit) {
      return;
    }
    const nextForm = {};
    users.forEach((user) => {
      const roleValue = user.role || "";
      const accessLevel = roleValue.toLowerCase().includes("админ")
        ? "admin"
        : roleValue.toLowerCase().includes("модератор")
          ? "moderator"
          : "member";
      nextForm[user.telegram_id] = {
        full_name: user.full_name || "",
        username: user.username || "",
        role: user.role || "",
        date_of_birthday: user.date_of_birthday || "",
        number_of_phone: user.number_of_phone || "",
        may_to_open: Boolean(user.may_to_open),
        access_level: accessLevel,
      };
    });
    setEditForm(nextForm);
  }, [users, canEdit]);

  const handleEditChange = (telegramId, field, value) => {
    setEditForm((prev) => ({
      ...prev,
      [telegramId]: {
        ...prev[telegramId],
        [field]: value,
      },
    }));
  };

  const handleSave = (telegramId) => {
    const payload = editForm[telegramId];
    if (!payload) {
      return;
    }
    const normalizedRole = payload.role
      .split(",")
      .map((part) => part.trim())
      .filter((part) => part && !["админ", "admin", "модератор", "moderator"].includes(part.toLowerCase()));
    if (payload.access_level === "moderator") {
      normalizedRole.push("Модератор");
    }
    if (payload.access_level === "admin") {
      normalizedRole.push("Модератор");
      normalizedRole.push("Админ");
    }
    const roleValue = normalizedRole.join(", ");
    fetch(`${API_BASE}/users/${telegramId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        ...getAuthHeaders(),
      },
      body: JSON.stringify({
        full_name: payload.full_name,
        username: payload.username,
        role: roleValue,
        date_of_birthday: payload.date_of_birthday,
        number_of_phone: payload.number_of_phone,
        may_to_open: payload.may_to_open,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("failed");
        }
      })
      .then(() =>
        fetch(`${API_BASE}/get_users`, {
          headers: getAuthHeaders(),
        })
      )
      .then((response) => response.json())
      .then((data) => {
        setUsers(Array.isArray(data) ? data : []);
        setEditingId("");
      })
      .catch(() => setUsersError("Не удалось обновить пользователя"));
  };

  return (
    <section className="members-screen">
      <header className="members-screen__header">
        <h2 className="members-screen__title">Участники</h2>
      </header>
      <div className="members-screen__list">
        {isLoadingUsers && (
          <article className="member-card" aria-busy="true">
            <div className="member-card__main">
              <p className="member-card__name">Загрузка...</p>
              <p className="member-card__hint">подождите немного</p>
            </div>
          </article>
        )}
        {!isLoadingUsers && usersError && (
          <article className="member-card">
            <div className="member-card__main">
              <p className="member-card__name">Ошибка</p>
              <p className="member-card__hint">{usersError}</p>
            </div>
          </article>
        )}
        {!isLoadingUsers && !usersError && users.length === 0 && (
          <article className="member-card">
            <div className="member-card__main">
              <p className="member-card__name">Пользователи не найдены</p>
              <p className="member-card__hint">
                пока нет загруженных профилей
              </p>
            </div>
          </article>
        )}
        {!isLoadingUsers &&
          !usersError &&
          users.length > 0 &&
          users.map((user) => {
            const displayName =
              user.full_name ||
              [user.first_name, user.last_name].filter(Boolean).join(" ") ||
              user.username ||
              "Без имени";
            const username = user.username ? `@${user.username}` : "—";
            const role = user.role || "—";
            const phone = user.number_of_phone || "—";
            const birthday = user.date_of_birthday || "—";
            const telegramId = user.telegram_id || "—";
            const accessLabel = user.may_to_open ? "открыт" : "закрыт";
            const isEditing = canEdit && editingId === telegramId;
            const editValues = editForm[telegramId] || {};

            return (
              <article className="member-card" key={telegramId}>
                <div className="member-card__main">
                  {isEditing ? (
                    <div className="member-card__edit">
                      <label>
                        ФИО
                        <input
                          value={editValues.full_name || ""}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "full_name",
                              event.target.value
                            )
                          }
                        />
                      </label>
                      <label>
                        Username
                        <input
                          value={editValues.username || ""}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "username",
                              event.target.value
                            )
                          }
                        />
                      </label>
                      <label className="member-card__edit--full">
                        Роли
                        <input
                          value={editValues.role || ""}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "role",
                              event.target.value
                            )
                          }
                        />
                      </label>
                      <label className="member-card__edit--full">
                        Уровень доступа
                        <select
                          value={editValues.access_level || "member"}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "access_level",
                              event.target.value
                            )
                          }
                        >
                          <option value="member">Участник</option>
                          <option value="moderator">Модератор</option>
                          <option value="admin">Админ</option>
                        </select>
                      </label>
                    </div>
                  ) : (
                    <>
                      <p className="member-card__name">{displayName}</p>
                      <p className="member-card__hint">{username}</p>
                    </>
                  )}
                </div>
                <div className="member-card__meta">
                  {isEditing ? (
                    <>
                      <label>
                        Телефон
                        <input
                          value={editValues.number_of_phone || ""}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "number_of_phone",
                              event.target.value
                            )
                          }
                        />
                      </label>
                      <label>
                        Дата рождения
                        <input
                          value={editValues.date_of_birthday || ""}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "date_of_birthday",
                              event.target.value
                            )
                          }
                        />
                      </label>
                      <label className="member-card__toggle">
                        <input
                          type="checkbox"
                          checked={Boolean(editValues.may_to_open)}
                          onChange={(event) =>
                            handleEditChange(
                              telegramId,
                              "may_to_open",
                              event.target.checked
                            )
                          }
                        />
                        Доступ
                      </label>
                    </>
                  ) : (
                    <>
                      <p>
                        <span>роль:</span> {role}
                      </p>
                      <p>
                        <span>телефон:</span> {phone}
                      </p>
                      <p>
                        <span>др:</span> {birthday}
                      </p>
                      <p>
                        <span>telegram:</span> {telegramId}
                      </p>
                      <p>
                        <span>доступ:</span> {accessLabel}
                      </p>
                    </>
                  )}
                </div>
                {canEdit && (
                  <div className="member-card__actions">
                    {isEditing ? (
                      <>
                        <button
                          type="button"
                          onClick={() => handleSave(telegramId)}
                        >
                          Сохранить
                        </button>
                        <button
                          type="button"
                          className="ghost"
                          onClick={() => setEditingId("")}
                        >
                          Отмена
                        </button>
                      </>
                    ) : (
                      <button
                        type="button"
                        onClick={() => setEditingId(telegramId)}
                      >
                        Редактировать
                      </button>
                    )}
                  </div>
                )}
              </article>
            );
          })}
      </div>
    </section>
  );
}
