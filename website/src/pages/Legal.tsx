import { Link } from "react-router";
import { useI18n } from "../i18n";

export function Imprint() {
  const { t } = useI18n();
  return <article className="wrap section narrow prose">{t.legal.imprint()}</article>;
}

export function Privacy() {
  const { t } = useI18n();
  return (
    <article className="wrap section narrow prose">
      {t.legal.privacy((text) => (
        <Link to="/dashboard">{text}</Link>
      ))}
    </article>
  );
}

export function NotFound() {
  const { t } = useI18n();
  return (
    <div className="wrap section narrow">
      <h1>{t.legal.notFoundTitle}</h1>
      <p>
        {t.legal.notFound((text) => (
          <Link to="/">{text}</Link>
        ))}
      </p>
    </div>
  );
}
