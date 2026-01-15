import { useMemo, useState } from "react";
import "../styles/ProfileScreen.scss";

const projectTemplate = {
  title: "Наша личная CRM-система",
  subtitle: "посмотреть задачи",
  status: "в разработке",
  participants: 3,
  role: "руководитель / разработчик",
};

const profileProjects = Array.from({ length: 4 }, (_, index) => ({
  ...projectTemplate,
  id: index,
}));

export default function ProfileScreen() {
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

  const profileFields = useMemo(
    () => [
      { label: "Имя", value: derivedProfile.firstName },
      { label: "Фамилия", value: derivedProfile.lastName },
      { label: "Дата рождения", value: derivedProfile.dateOfBirthday },
      { label: "Должность", value: derivedProfile.role },
      { label: "TelegramID", value: derivedProfile.telegramID },
    ],
    [derivedProfile]
  );

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
            {profileProjects.map((project) => (
              <article className="project-card" key={project.id}>
                <div className="project-card__main">
                  <p className="project-card__title">{project.title}</p>
                  <p className="project-card__subtitle">{project.subtitle}</p>
                </div>
                <div className="project-card__meta">
                  <p>
                    <span>статус:</span> {project.status}
                  </p>
                  <p>
                    <span>кол-во участников:</span> {project.participants}
                  </p>
                  <p>
                    <span>роль:</span> {project.role}
                  </p>
                </div>
              </article>
            ))}
          </div>
        </section>
      </div>
    </section>
  );
}
