import React from "react";
import ReactDOM from "react-dom/client";
import {
  Activity,
  ArrowDownToLine,
  ArrowUpDown,
  BadgeDollarSign,
  CheckCircle2,
  CircleAlert,
  Coins,
  KeyRound,
  LayoutDashboard,
  Loader2,
  LogOut,
  ShieldCheck,
  UserCircle,
  UserPlus,
  Wallet
} from "lucide-react";
import "./styles.css";

type ApiStatus = "checking" | "online" | "offline";
type Page = "overview" | "trade" | "wallet" | "account" | "identity";
type IdentityMode = "register" | "login";

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

function App() {
  const [apiStatus, setApiStatus] = React.useState<ApiStatus>("checking");
  const [page, setPage] = React.useState<Page>("overview");
  const [identityMode, setIdentityMode] = React.useState<IdentityMode>("login");
  const [registrationEmail, setRegistrationEmail] = React.useState("");
  const [registrationPassword, setRegistrationPassword] = React.useState("");
  const [loginEmail, setLoginEmail] = React.useState("");
  const [loginPassword, setLoginPassword] = React.useState("");
  const [isRegistering, setIsRegistering] = React.useState(false);
  const [isLoggingIn, setIsLoggingIn] = React.useState(false);
  const [registrationResult, setRegistrationResult] = React.useState<RegistrationResponse | null>(null);
  const [session, setSession] = React.useState<LoginResponse | null>(null);
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
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: registrationEmail, password: registrationPassword })
      });
      const payload = (await response.json()) as RegistrationResponse | ErrorResponse;

      if (!response.ok) {
        setNotice("message" in payload ? payload.message : "Registration failed");
        return;
      }

      const registeredUser = payload as RegistrationResponse;
      setRegistrationResult(registeredUser);
      setLoginEmail(registeredUser.email);
      setIdentityMode("login");
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

    try {
      const response = await fetch("/identity/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: loginEmail, password: loginPassword })
      });
      const payload = (await response.json()) as LoginResponse | ErrorResponse;

      if (!response.ok) {
        setNotice("message" in payload ? payload.message : "Login failed");
        return;
      }

      setSession(payload as LoginResponse);
      setLoginPassword("");
      setPage("account");
    } catch {
      setNotice("The API is unreachable.");
    } finally {
      setIsLoggingIn(false);
    }
  }

  function logout() {
    setSession(null);
    setPage("identity");
  }

  return (
    <main className="shell">
      <section className="workspace" aria-label="Goldin trading workspace">
        <aside className="sidebar">
          <div className="brand">
            <div className="brandMark">G</div>
            <div>
              <strong>Goldin</strong>
              <span>Gold Trading</span>
            </div>
          </div>

          <nav className="nav" aria-label="Primary">
            <NavButton icon={<LayoutDashboard size={18} />} label="Overview" page="overview" active={page} onClick={setPage} />
            <NavButton icon={<ArrowUpDown size={18} />} label="Trade" page="trade" active={page} onClick={setPage} />
            <NavButton icon={<Wallet size={18} />} label="Wallet" page="wallet" active={page} onClick={setPage} />
            <NavButton icon={<UserCircle size={18} />} label="Account" page="account" active={page} onClick={setPage} />
            <NavButton icon={<KeyRound size={18} />} label="Access" page="identity" active={page} onClick={setPage} />
          </nav>
        </aside>

        <section className="content">
          <header className="topbar">
            <div>
              <p className="eyebrow">{pageLabel(page)}</p>
              <h1>{pageTitle(page, session)}</h1>
            </div>
            <div className="topActions">
              {session && (
                <button className="ghostButton" type="button" onClick={logout}>
                  <LogOut size={18} />
                  <span>Logout</span>
                </button>
              )}
              <button className={`statusButton ${apiStatus}`} type="button" onClick={checkHealth}>
                {apiStatus === "checking" && <Loader2 className="spin" size={18} />}
                {apiStatus === "online" && <CheckCircle2 size={18} />}
                {apiStatus === "offline" && <CircleAlert size={18} />}
                <span>{statusLabel(apiStatus)}</span>
              </button>
            </div>
          </header>

          {page === "overview" && <OverviewPage apiStatus={apiStatus} session={session} />}
          {page === "trade" && <TradePage session={session} onAccess={() => setPage("identity")} />}
          {page === "wallet" && <WalletPage session={session} onAccess={() => setPage("identity")} />}
          {page === "account" && <AccountPage session={session} onAccess={() => setPage("identity")} />}
          {page === "identity" && (
            <IdentityPage
              mode={identityMode}
              setMode={setIdentityMode}
              apiStatus={apiStatus}
              registrationEmail={registrationEmail}
              registrationPassword={registrationPassword}
              loginEmail={loginEmail}
              loginPassword={loginPassword}
              isRegistering={isRegistering}
              isLoggingIn={isLoggingIn}
              registrationResult={registrationResult}
              session={session}
              notice={notice}
              setRegistrationEmail={setRegistrationEmail}
              setRegistrationPassword={setRegistrationPassword}
              setLoginEmail={setLoginEmail}
              setLoginPassword={setLoginPassword}
              submitRegistration={submitRegistration}
              submitLogin={submitLogin}
            />
          )}
        </section>
      </section>
    </main>
  );
}

function NavButton(props: {
  icon: React.ReactNode;
  label: string;
  page: Page;
  active: Page;
  onClick: (page: Page) => void;
}) {
  return (
    <button className={`navItem ${props.active === props.page ? "active" : ""}`} type="button" onClick={() => props.onClick(props.page)}>
      {props.icon}
      {props.label}
    </button>
  );
}

function OverviewPage({ apiStatus, session }: { apiStatus: ApiStatus; session: LoginResponse | null }) {
  return (
    <div className="pageStack">
      <section className="marketHero">
        <div>
          <p className="eyebrow">Goldin Markets</p>
          <h2>Physical gold trading workspace</h2>
        </div>
        <div className="heroMetric">
          <span>Session</span>
          <strong>{session ? "Authenticated" : "Guest"}</strong>
        </div>
      </section>

      <div className="metricsGrid">
        <MetricCard icon={<Activity size={20} />} label="API" value={statusLabel(apiStatus)} tone={apiStatus} />
        <MetricCard icon={<Coins size={20} />} label="Gold balance" value="-- g" />
        <MetricCard icon={<BadgeDollarSign size={20} />} label="Cash wallet" value="--" />
        <MetricCard icon={<ArrowUpDown size={20} />} label="Orders" value="--" />
      </div>

      <div className="grid">
        <section className="panel">
          <div className="panelHeader">
            <ArrowUpDown size={20} />
            <h2>Trading</h2>
          </div>
          <div className="emptyState compact">
            <Coins size={28} />
            <p>Market pricing will appear here after the pricing module is available.</p>
          </div>
        </section>

        <section className="panel">
          <div className="panelHeader">
            <Wallet size={20} />
            <h2>Wallets</h2>
          </div>
          <div className="emptyState compact">
            <Wallet size={28} />
            <p>Cash and gold wallets will appear here after the wallet module is available.</p>
          </div>
        </section>
      </div>
    </div>
  );
}

function TradePage({ session, onAccess }: { session: LoginResponse | null; onAccess: () => void }) {
  return (
    <div className="grid tradeGrid">
      <section className="panel">
        <div className="panelHeader">
          <Coins size={20} />
          <h2>Gold market</h2>
        </div>
        <div className="quoteBoard">
          <div><span>Bid</span><strong>--</strong></div>
          <div><span>Ask</span><strong>--</strong></div>
          <div><span>Spread</span><strong>--</strong></div>
        </div>
      </section>

      <section className="panel">
        <div className="panelHeader">
          <ArrowUpDown size={20} />
          <h2>Order ticket</h2>
        </div>
        {session ? (
          <div className="emptyState compact">
            <ArrowUpDown size={28} />
            <p>Order entry will activate after trading APIs are implemented.</p>
          </div>
        ) : (
          <AccessPrompt onAccess={onAccess} />
        )}
      </section>
    </div>
  );
}

function WalletPage({ session, onAccess }: { session: LoginResponse | null; onAccess: () => void }) {
  if (!session) {
    return <AccessPrompt onAccess={onAccess} />;
  }

  return (
    <div className="grid">
      <section className="panel walletPanel">
        <div className="panelHeader"><Wallet size={20} /><h2>Cash wallet</h2></div>
        <strong className="walletValue">--</strong>
        <div className="buttonRow">
          <button className="secondaryButton" type="button" disabled><ArrowDownToLine size={18} />Deposit</button>
          <button className="secondaryButton" type="button" disabled>Withdraw</button>
        </div>
      </section>

      <section className="panel walletPanel">
        <div className="panelHeader"><Coins size={20} /><h2>Gold balance</h2></div>
        <strong className="walletValue">-- g</strong>
        <div className="buttonRow">
          <button className="secondaryButton" type="button" disabled>Buy</button>
          <button className="secondaryButton" type="button" disabled>Sell</button>
        </div>
      </section>
    </div>
  );
}

function AccountPage({ session, onAccess }: { session: LoginResponse | null; onAccess: () => void }) {
  if (!session) {
    return <AccessPrompt onAccess={onAccess} />;
  }

  return (
    <div className="grid accountGrid">
      <section className="panel detailPanel">
        <div className="panelHeader"><UserCircle size={20} /><h2>Account</h2></div>
        <dl>
          <div><dt>User ID</dt><dd>{session.user_id}</dd></div>
          <div><dt>Email</dt><dd>{session.email}</dd></div>
          <div><dt>Token type</dt><dd>{session.token_type}</dd></div>
        </dl>
      </section>

      <section className="panel detailPanel">
        <div className="panelHeader"><ShieldCheck size={20} /><h2>Session</h2></div>
        <SessionResult session={session} />
      </section>
    </div>
  );
}

function IdentityPage(props: {
  mode: IdentityMode;
  setMode: (mode: IdentityMode) => void;
  apiStatus: ApiStatus;
  registrationEmail: string;
  registrationPassword: string;
  loginEmail: string;
  loginPassword: string;
  isRegistering: boolean;
  isLoggingIn: boolean;
  registrationResult: RegistrationResponse | null;
  session: LoginResponse | null;
  notice: string | null;
  setRegistrationEmail: (value: string) => void;
  setRegistrationPassword: (value: string) => void;
  setLoginEmail: (value: string) => void;
  setLoginPassword: (value: string) => void;
  submitRegistration: (event: React.FormEvent<HTMLFormElement>) => void;
  submitLogin: (event: React.FormEvent<HTMLFormElement>) => void;
}) {
  return (
    <div className="grid">
      <section className="panel" id="identity">
        <div className="panelHeader"><ShieldCheck size={20} /><h2>{props.mode === "register" ? "Create identity account" : "Authenticate user"}</h2></div>
        <div className="segmented" role="tablist" aria-label="Identity actions">
          <button className={props.mode === "register" ? "selected" : ""} type="button" onClick={() => props.setMode("register")}><UserPlus size={17} />Register</button>
          <button className={props.mode === "login" ? "selected" : ""} type="button" onClick={() => props.setMode("login")}><KeyRound size={17} />Login</button>
        </div>

        {props.mode === "register" && (
          <form className="form" onSubmit={props.submitRegistration}>
            <label><span>Email</span><input type="email" value={props.registrationEmail} onChange={(event) => props.setRegistrationEmail(event.target.value)} autoComplete="email" placeholder="user@example.com" required /></label>
            <label><span>Password</span><input type="password" value={props.registrationPassword} onChange={(event) => props.setRegistrationPassword(event.target.value)} autoComplete="new-password" minLength={12} placeholder="At least 12 characters" required /></label>
            <button className="primaryButton" type="submit" disabled={props.isRegistering || props.apiStatus !== "online"}>{props.isRegistering ? <Loader2 className="spin" size={18} /> : <UserPlus size={18} />}<span>{props.isRegistering ? "Creating account" : "Create account"}</span></button>
          </form>
        )}

        {props.mode === "login" && (
          <form className="form" onSubmit={props.submitLogin}>
            <label><span>Email</span><input type="email" value={props.loginEmail} onChange={(event) => props.setLoginEmail(event.target.value)} autoComplete="email" placeholder="user@example.com" required /></label>
            <label><span>Password</span><input type="password" value={props.loginPassword} onChange={(event) => props.setLoginPassword(event.target.value)} autoComplete="current-password" placeholder="Account password" required /></label>
            <button className="primaryButton" type="submit" disabled={props.isLoggingIn || props.apiStatus !== "online"}>{props.isLoggingIn ? <Loader2 className="spin" size={18} /> : <KeyRound size={18} />}<span>{props.isLoggingIn ? "Authenticating" : "Login"}</span></button>
          </form>
        )}
      </section>

      <section className="panel detailPanel" aria-live="polite">
        <div className="panelHeader"><Activity size={20} /><h2>{props.session ? "Authenticated session" : "Identity result"}</h2></div>
        {props.session && <SessionResult session={props.session} />}
        {!props.session && props.registrationResult && <RegistrationResult result={props.registrationResult} />}
        {props.notice && <div className="result danger"><CircleAlert size={22} /><p>{props.notice}</p></div>}
        {!props.registrationResult && !props.session && !props.notice && <div className="emptyState"><ShieldCheck size={28} /><p>No identity action submitted yet.</p></div>}
      </section>
    </div>
  );
}

function RegistrationResult({ result }: { result: RegistrationResponse }) {
  return (
    <div className="result success">
      <CheckCircle2 size={22} />
      <dl>
        <div><dt>User ID</dt><dd>{result.user_id}</dd></div>
        <div><dt>Email</dt><dd>{result.email}</dd></div>
        <div><dt>Registered</dt><dd>{formatDate(result.registered_at)}</dd></div>
      </dl>
    </div>
  );
}

function SessionResult({ session }: { session: LoginResponse }) {
  return (
    <div className="result success tokenResult">
      <CheckCircle2 size={22} />
      <dl>
        <div><dt>User ID</dt><dd>{session.user_id}</dd></div>
        <div><dt>Email</dt><dd>{session.email}</dd></div>
        <div><dt>Access token</dt><dd>{compactToken(session.access_token)}</dd></div>
        <div><dt>Refresh token</dt><dd>{compactToken(session.refresh_token)}</dd></div>
        <div><dt>Access expires</dt><dd>{formatDuration(session.access_token_expires_in)}</dd></div>
      </dl>
    </div>
  );
}

function MetricCard({ icon, label, value, tone }: { icon: React.ReactNode; label: string; value: string; tone?: ApiStatus }) {
  return <section className={`metricCard ${tone ?? ""}`}>{icon}<span>{label}</span><strong>{value}</strong></section>;
}

function AccessPrompt({ onAccess }: { onAccess: () => void }) {
  return (
    <section className="panel accessPrompt">
      <ShieldCheck size={30} />
      <h2>Account access required</h2>
      <button className="primaryButton" type="button" onClick={onAccess}><KeyRound size={18} /><span>Go to access</span></button>
    </section>
  );
}

function pageLabel(page: Page) {
  switch (page) {
    case "overview": return "Workspace";
    case "trade": return "Trading";
    case "wallet": return "Wallet";
    case "account": return "Account";
    case "identity": return "Identity";
  }
}

function pageTitle(page: Page, session: LoginResponse | null) {
  if (page === "account" && session) return session.email;
  switch (page) {
    case "overview": return "Gold trading overview";
    case "trade": return "Buy and sell gold";
    case "wallet": return "Wallet balances";
    case "account": return "Account access";
    case "identity": return "Register or login";
  }
}

function statusLabel(status: ApiStatus) {
  switch (status) {
    case "checking": return "Checking API";
    case "online": return "API online";
    case "offline": return "API offline";
  }
}

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat(undefined, { dateStyle: "medium", timeStyle: "short" }).format(date);
}

function formatDuration(seconds: number) {
  if (seconds < 60) return `${seconds} seconds`;
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) return `${minutes} minutes`;
  return `${Math.round(minutes / 60)} hours`;
}

function compactToken(token: string) {
  if (token.length <= 24) return token;
  return `${token.slice(0, 14)}...${token.slice(-10)}`;
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
