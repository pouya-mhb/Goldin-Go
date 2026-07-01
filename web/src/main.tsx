import React from "react";
import ReactDOM from "react-dom/client";
import {
  Activity,
  CheckCircle2,
  CircleAlert,
  KeyRound,
  Loader2,
  ShieldCheck,
  UserPlus
} from "lucide-react";
import "./styles.css";

type ApiStatus = "checking" | "online" | "offline";

type RegistrationResponse = {
  user_id: string;
  email: string;
  registered_at: string;
};

type LoginResponse = {
  user_id: string;
  email: string;
  access_token: string;
  refresh_token: string;
  token_type: string;
  access_token_expires_in: number;
  refresh_token_expires_in: number;
};

type ErrorResponse = {
  error: string;
  message: string;
};

type IdentityMode = "register" | "login";

function App() {
  const [apiStatus, setApiStatus] = React.useState<ApiStatus>("checking");
  const [mode, setMode] = React.useState<IdentityMode>("register");
  const [registrationEmail, setRegistrationEmail] = React.useState("");
  const [registrationPassword, setRegistrationPassword] = React.useState("");
  const [loginEmail, setLoginEmail] = React.useState("");
  const [loginPassword, setLoginPassword] = React.useState("");
  const [isRegistering, setIsRegistering] = React.useState(false);
  const [isLoggingIn, setIsLoggingIn] = React.useState(false);
  const [registrationResult, setRegistrationResult] = React.useState<RegistrationResponse | null>(null);
  const [loginResult, setLoginResult] = React.useState<LoginResponse | null>(null);
  const [notice, setNotice] = React.useState<string | null>(null);

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
    setIsRegistering(true);
    setNotice(null);
    setRegistrationResult(null);

    try {
      const response = await fetch("/identity/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ email: registrationEmail, password: registrationPassword })
      });

      const payload = (await response.json()) as RegistrationResponse | ErrorResponse;

      if (!response.ok) {
        const message = "message" in payload ? payload.message : "Registration failed";
        setNotice(message);
        return;
      }

      const registeredUser = payload as RegistrationResponse;
      setRegistrationResult(registeredUser);
      setLoginEmail(registeredUser.email);
      setRegistrationPassword("");
    } catch {
      setNotice("The API is unreachable.");
    } finally {
      setIsRegistering(false);
    }
  }

  async function submitLogin(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsLoggingIn(true);
    setNotice(null);
    setLoginResult(null);

    try {
      const response = await fetch("/identity/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ email: loginEmail, password: loginPassword })
      });

      const payload = (await response.json()) as LoginResponse | ErrorResponse;

      if (!response.ok) {
        const message = "message" in payload ? payload.message : "Login failed";
        setNotice(message);
        return;
      }

      setLoginResult(payload as LoginResponse);
      setLoginPassword("");
    } catch {
      setNotice("The API is unreachable.");
    } finally {
      setIsLoggingIn(false);
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
              <h1>Identity access</h1>
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
                <h2>{mode === "register" ? "Create identity account" : "Authenticate user"}</h2>
              </div>

              <div className="segmented" role="tablist" aria-label="Identity actions">
                <button
                  className={mode === "register" ? "selected" : ""}
                  type="button"
                  onClick={() => setMode("register")}
                >
                  <UserPlus size={17} />
                  Register
                </button>
                <button className={mode === "login" ? "selected" : ""} type="button" onClick={() => setMode("login")}>
                  <KeyRound size={17} />
                  Login
                </button>
              </div>

              {mode === "register" && (
                <form className="form" onSubmit={submitRegistration}>
                  <label>
                    <span>Email</span>
                    <input
                      type="email"
                      value={registrationEmail}
                      onChange={(event) => setRegistrationEmail(event.target.value)}
                      autoComplete="email"
                      placeholder="user@example.com"
                      required
                    />
                  </label>

                  <label>
                    <span>Password</span>
                    <input
                      type="password"
                      value={registrationPassword}
                      onChange={(event) => setRegistrationPassword(event.target.value)}
                      autoComplete="new-password"
                      minLength={12}
                      placeholder="At least 12 characters"
                      required
                    />
                  </label>

                  <button className="primaryButton" type="submit" disabled={isRegistering || apiStatus !== "online"}>
                    {isRegistering ? <Loader2 className="spin" size={18} /> : <UserPlus size={18} />}
                    <span>{isRegistering ? "Creating account" : "Create account"}</span>
                  </button>
                </form>
              )}

              {mode === "login" && (
                <form className="form" onSubmit={submitLogin}>
                  <label>
                    <span>Email</span>
                    <input
                      type="email"
                      value={loginEmail}
                      onChange={(event) => setLoginEmail(event.target.value)}
                      autoComplete="email"
                      placeholder="user@example.com"
                      required
                    />
                  </label>

                  <label>
                    <span>Password</span>
                    <input
                      type="password"
                      value={loginPassword}
                      onChange={(event) => setLoginPassword(event.target.value)}
                      autoComplete="current-password"
                      placeholder="Account password"
                      required
                    />
                  </label>

                  <button className="primaryButton" type="submit" disabled={isLoggingIn || apiStatus !== "online"}>
                    {isLoggingIn ? <Loader2 className="spin" size={18} /> : <KeyRound size={18} />}
                    <span>{isLoggingIn ? "Authenticating" : "Login"}</span>
                  </button>
                </form>
              )}
            </section>

            <section className="panel detailPanel" aria-live="polite">
              <div className="panelHeader">
                <Activity size={20} />
                <h2>{loginResult ? "Authenticated session" : "Identity result"}</h2>
              </div>

              {loginResult && (
                <div className="result success tokenResult">
                  <CheckCircle2 size={22} />
                  <dl>
                    <div>
                      <dt>User ID</dt>
                      <dd>{loginResult.user_id}</dd>
                    </div>
                    <div>
                      <dt>Email</dt>
                      <dd>{loginResult.email}</dd>
                    </div>
                    <div>
                      <dt>Access token</dt>
                      <dd>{compactToken(loginResult.access_token)}</dd>
                    </div>
                    <div>
                      <dt>Refresh token</dt>
                      <dd>{compactToken(loginResult.refresh_token)}</dd>
                    </div>
                    <div>
                      <dt>Access expires</dt>
                      <dd>{formatDuration(loginResult.access_token_expires_in)}</dd>
                    </div>
                  </dl>
                </div>
              )}

              {!loginResult && registrationResult && (
                <div className="result success">
                  <CheckCircle2 size={22} />
                  <dl>
                    <div>
                      <dt>User ID</dt>
                      <dd>{registrationResult.user_id}</dd>
                    </div>
                    <div>
                      <dt>Email</dt>
                      <dd>{registrationResult.email}</dd>
                    </div>
                    <div>
                      <dt>Registered</dt>
                      <dd>{formatDate(registrationResult.registered_at)}</dd>
                    </div>
                  </dl>
                </div>
              )}

              {notice && (
                <div className="result danger">
                  <CircleAlert size={22} />
                  <p>{notice}</p>
                </div>
              )}

              {!registrationResult && !loginResult && !notice && (
                <div className="emptyState">
                  <ShieldCheck size={28} />
                  <p>No identity action submitted yet.</p>
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

function formatDuration(seconds: number) {
  if (seconds < 60) {
    return `${seconds} seconds`;
  }

  const minutes = Math.round(seconds / 60);
  if (minutes < 60) {
    return `${minutes} minutes`;
  }

  const hours = Math.round(minutes / 60);
  return `${hours} hours`;
}

function compactToken(token: string) {
  if (token.length <= 24) {
    return token;
  }

  return `${token.slice(0, 14)}...${token.slice(-10)}`;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
