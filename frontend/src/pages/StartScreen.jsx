import { useEffect, useState } from "react";

export default function StartScreen() {
  const [asciiArt, setAsciiArt] = useState("");

  useEffect(() => {
    fetch(`${process.env.PUBLIC_URL}/ascii-art.txt`)
      .then((response) => response.text())
      .then((text) => setAsciiArt(text))
      .catch(() => setAsciiArt(""));
  }, []);

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
          <button className="btn btn--accent btn--cta" type="button">
            Авторизация
          </button>
        </div>
      </section>

      <aside className="ascii-panel" aria-hidden="true">
        <pre className="ascii-panel__art">{asciiArt}</pre>
      </aside>
    </div>
  );
}
