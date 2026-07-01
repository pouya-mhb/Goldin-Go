import React from "react";
import ReactDOM from "react-dom/client";
import { Activity, CheckCircle2, CircleAlert, Loader2, ShieldCheck, UserPlus } from "lucide-react";
import "./styles.css";

type ApiStatus = "checking" | "online" | "offline";

type RegistrationResponse = {
  user_id: string;
  email: string;
  registered_at: string;
};

type ErrorResponse = {
  error: string;
  message: string;
};

function App() {
  const [apiStatus, setApiStatus] = React.useState<ApiStatus>("checking");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [result, setResult] = React.useState<RegistrationResponse | null>(null);
  const [error, setError] = React.useState<string | null>(null);

  const checkHealth = React.useCallback(async () => {
    setApiStatus("checking");

    try {
      const response = await fetch("/health", { method: "GET" });
      setApiStatus(response.ok ? "online" : "offline");
    } catch {
      setApiStatus("offline");
    }
  }, []);

  React.useEffect(() => {
    void checkHealth();
  }, [checkHealth]);

  async function submitRegistration(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsSubmitting(true);
    setError(null);
    setResult(null);

    try {
      const response = await fetch("/identity/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ email, password })
      });

      const payload = (await response.json()) as RegistrationResponse | ErrorResponse;

      if (!response.ok) {
        const message = "message" in payload ? payload.message : "Registration failed";
        setError(message);
        return;
      }

      setResult(payload as RegistrationResponse);
      setPassword("");
    } catch {
      setError("The API is unreachable.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="shell">
      <section className="workspace" aria-label="Goldin operations console">
        <aside className="sidebar">
          <div className="brand">
            <div className="brandMark">G</div>
            <div>
              <strong>Goldin</strong>
              <span>Trading Platform</span>
            </div>
          </div>

          <nav className="nav" aria-label="Primary">
            <a className="navItem active" href="#identity">
              <UserPlus size={18} />
              Identity
            </a>
            <a className="navItem disabled" href="#health" aria-disabled="true">
              <Activity size={18} />
              Platform
            </a>
          </nav>
        </aside>

        <section className="content">
          <header className="topbar">
            <div>
              <p className="eyebrow">Identity</p>
              <h1>User registration</h1>
            </div>
            <button className={`statusButton ${apiStatus}`} type="button" onClick={checkHealth}>
              {apiStatus === "checking" && <Loader2 className="spin" size={18} />}
              {apiStatus === "online" && <CheckCircle2 size={18} />}
              {apiStatus === "offline" && <CircleAlert size={18} />}
              <span>{statusLabel(apiStatus)}</span>
            </button>
          </header>

          <div className="grid">
            <section className="panel" id="identity">
              <div className="panelHeader">
                <ShieldCheck size={20} />
                <h2>Create identity account</h2>
              </div>

              <form className="form" onSubmit={submitRegistration}>
                <label>
                  <span>Email</span>
                  <input
                    type="email"
                    value={email}
                    onChange={(event) => setEmail(event.target.value)}
                    autoComplete="email"
                    placeholder="user@example.com"
                    required
                  />
                </label>

                <label>
                  <span>Password</span>
                  <input
                    type="password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    autoComplete="new-password"
                    minLength={12}
                    placeholder="At least 12 characters"
                    required
                  />
                </label>

                <button className="primaryButton" type="submit" disabled={isSubmitting || apiStatus !== "online"}>
                  {isSubmitting ? <Loader2 className="spin" size={18} /> : <UserPlus size={18} />}
                  <span>{isSubmitting ? "Creating account" : "Create account"}</span>
                </button>
              </form>
            </section>

            <section className="panel detailPanel" aria-live="polite">
              <div className="panelHeader">
                <Activity size={20} />
                <h2>Registration result</h2>
              </div>

              {result && (
                <div className="result success">
                  <CheckCircle2 size={22} />
                  <dl>
                    <div>
                      <dt>User ID</dt>
                      <dd>{result.user_id}</dd>
                    </div>
                    <div>
                      <dt>Email</dt>
                      <dd>{result.email}</dd>
                    </div>
                    <div>
                      <dt>Registered</dt>
                      <dd>{formatDate(result.registered_at)}</dd>
                    </div>
                  </dl>
                </div>
              )}

              {error && (
                <div className="result danger">
                  <CircleAlert size={22} />
                  <p>{error}</p>
                </div>
              )}

              {!result && !error && (
                <div className="emptyState">
                  <ShieldCheck size={28} />
                  <p>No registration submitted yet.</p>
                </div>
              )}
            </section>
          </div>
        </section>
      </section>
    </main>
  );
}

function statusLabel(status: ApiStatus) {
  switch (status) {
    case "checking":
      return "Checking API";
    case "online":
      return "API online";
    case "offline":
      return "API offline";
  }
}

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(date);
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
