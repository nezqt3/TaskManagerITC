import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import "../styles/StartScreen.scss";

export default function StartScreen() {
  const [asciiArt, setAsciiArt] = useState("");
  const [authError, setAuthError] = useState("");
  const navigate = useNavigate();
  const widgetRef = useRef(null);
  const BOT_USERNAME =
    process.env.REACT_APP_TELEGRAM_BOT_USERNAME || "it_communitytest_bot";
  const API_BASE = process.env.REACT_APP_API_URL || "http://localhost:8080";
  const telegramLoginUrl = `https://t.me/${BOT_USERNAME}?start=login`;

  const handleTelegramAuth = useCallback(
    async (user) => {
      setAuthError("");
      try {
        const response = await fetch(`${API_BASE}/auth/telegram`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(user),
        });

        if (!response.ok) {
          throw new Error("Ошибка авторизации");
        }

        const payload = await response.json();
        if (payload?.jwt) {
          localStorage.setItem("jwt", payload.jwt);
        }
        if (payload?.profile) {
          localStorage.setItem("profile", JSON.stringify(payload.profile));
        }

        navigate("/profile");
      } catch (error) {
        setAuthError("Не удалось авторизоваться через Telegram");
      }
    },
    [API_BASE, navigate],
  );

  useEffect(() => {
    fetch(`${process.env.PUBLIC_URL}/ascii-art.txt`)
      .then((response) => response.text())
      .then((text) => setAsciiArt(text))
      .catch(() => setAsciiArt(""));
  }, []);

  useEffect(() => {
    window.onTelegramAuth = handleTelegramAuth;
    return () => {
      delete window.onTelegramAuth;
    };
  }, [handleTelegramAuth]);

  useEffect(() => {
    if (!widgetRef.current) {
      return;
    }

    widgetRef.current.innerHTML = "";
    const script = document.createElement("script");
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.async = true;
    script.setAttribute("data-telegram-login", BOT_USERNAME);
    script.setAttribute("data-size", "large");
    script.setAttribute("data-userpic", "true");
    script.setAttribute("data-request-access", "write");
    script.setAttribute("data-onauth", "onTelegramAuth(user)");
    widgetRef.current.appendChild(script);
  }, [BOT_USERNAME]);

  return (
    <div className="start-screen">
      <section className="hero">
        <div className="hero__top">
          <div className="hero__actions">
            <button className="btn btn--ghost" type="button">
              Поддержка
            </button>
            <button className="btn btn--primary" type="button">
              Соцсети
            </button>
          </div>
        </div>

        <div className="hero__content">
          <h1 className="hero__title">TaskManager</h1>
          <p className="hero__subtitle">it-сообщество</p>
          <div className="hero__auth">
            <div ref={widgetRef} className="hero__auth-widget" />
            {authError && <p className="hero__auth-error">{authError}</p>}
          </div>
        </div>
      </section>

      <aside className="ascii-panel" aria-hidden="true">
        <pre className="ascii-panel__art">{asciiArt}</pre>
      </aside>
    </div>
  );
}
