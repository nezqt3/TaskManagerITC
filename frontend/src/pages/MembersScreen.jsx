import { useEffect, useState } from "react";
import "../styles/MembersScreen.scss";

const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";

export default function MembersScreen() {
  const [users, setUsers] = useState([]);
  const [usersError, setUsersError] = useState("");
  const [isLoadingUsers, setIsLoadingUsers] = useState(true);

  useEffect(() => {
    let isActive = true;
    setIsLoadingUsers(true);
    setUsersError("");

    fetch(`${API_BASE}/get_users`)
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
        setUsers(Array.isArray(data) ? data : []);
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
  }, []);

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

            return (
              <article className="member-card" key={telegramId}>
                <div className="member-card__main">
                  <p className="member-card__name">{displayName}</p>
                  <p className="member-card__hint">{username}</p>
                </div>
                <div className="member-card__meta">
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
                </div>
              </article>
            );
          })}
      </div>
    </section>
  );
}
