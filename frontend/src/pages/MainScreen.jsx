import { useMemo } from "react";
import "../styles/MainScreen.scss";

function PlaceholderPanel({ title, cards }) {
  const placeholders = useMemo(
    () => Array.from({ length: cards }, (_, index) => index),
    [cards]
  );

  return (
    <section className="panel">
      <h3 className="panel__title">{title}</h3>
      <div className="panel__stack">
        {placeholders.map((index) => (
          <div className="card-placeholder" key={index} aria-hidden="true">
            <span />
            <span />
          </div>
        ))}
      </div>
    </section>
  );
}

export default function MainScreen() {
  return (
    <div className="dashboard">
      <PlaceholderPanel title="Текущие задачи:" cards={4} />
      <PlaceholderPanel title="Текущие проекты:" cards={3} />
      <PlaceholderPanel title="События:" cards={3} />
    </div>
  );
}
