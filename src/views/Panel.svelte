<script lang="ts">
    import { api } from "../lib/api";

    let username = $state("");
    let password = $state("");
    let isLoading = $state(false);
    let errorMessage = $state("");
    let isAuthenticated = $state(false);
    let checkingAuth = $state(true);

    async function checkAuth() {
        const token = sessionStorage.getItem("panel_token");
        if (!token) {
            checkingAuth = false;
            return;
        }

        try {
            const res = await api("/api/auth/verify", {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });

            if (res.ok) {
                const data = await res.json();
                if (data.valid) {
                    isAuthenticated = true;
                } else {
                    sessionStorage.removeItem("panel_token");
                }
            } else {
                sessionStorage.removeItem("panel_token");
            }
        } catch {
            sessionStorage.removeItem("panel_token");
        } finally {
            checkingAuth = false;
        }
    }

    async function handleLogin(e: SubmitEvent) {
        e.preventDefault();
        errorMessage = "";
        isLoading = true;

        try {
            const res = await api("/api/auth/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ username, password })
            });

            const data = await res.json();

            if (res.ok && data.success && data.token) {
                sessionStorage.setItem("panel_token", data.token);
                isAuthenticated = true;
                password = "";
            } else {
                errorMessage = data.message || "Invalid credentials";
            }
        } catch {
            errorMessage = "Server connection error";
        } finally {
            isLoading = false;
        }
    }

    async function handleLogout() {
        const token = sessionStorage.getItem("panel_token");
        if (token) {
            try {
                await api("/api/auth/logout", {
                    method: "POST",
                    headers: {
                        Authorization: `Bearer ${token}`
                    }
                });
            } catch {}
        }
        sessionStorage.removeItem("panel_token");
        isAuthenticated = false;
        username = "";
        password = "";
    }

    checkAuth();
</script>

<div class="panel_main">
    {#if checkingAuth}
        <div class="status_text">Loading...</div>
    {:else if isAuthenticated}
        <div class="panel_card">
            <div class="panel_header">
                <span class="title">Admin Panel</span>
                <button type="button" class="logout_btn" onclick={handleLogout}>
                    Logout
                </button>
            </div>
            <div class="separator"></div>
            <div class="panel_body">
                <span class="desc">Logged in successfully.</span>
            </div>
            <div class="separator"></div>
            <a href="/" class="back_link">← Return home</a>
        </div>
    {:else}
        <div class="login_card">
            <div class="card_header">
                <span class="title">Panel</span>
                <span class="subtitle">Authentication required</span>
            </div>

            <div class="separator"></div>

            {#if errorMessage}
                <div class="error_msg">{errorMessage}</div>
            {/if}

            <form class="login_form" onsubmit={handleLogin}>
                <div class="input_field">
                    <input
                        type="text"
                        bind:value={username}
                        placeholder="Username"
                        required
                        autocomplete="username"
                    />
                </div>

                <div class="input_field">
                    <input
                        type="password"
                        bind:value={password}
                        placeholder="Password"
                        required
                        autocomplete="current-password"
                    />
                </div>

                <button type="submit" class="submit_btn" disabled={isLoading}>
                    {isLoading ? "Signing in..." : "Continue"}
                </button>
            </form>

            <div class="separator"></div>

            <a href="/" class="back_link">← Return home</a>
        </div>
    {/if}
</div>

<style>
    .panel_main {
        width: 100%;
        min-height: 100vh;
        background-color: #11111b;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        padding: 40px 20px;
        box-sizing: border-box;
    }

    .status_text {
        color: #8588a3;
        font-size: 0.9rem;
        font-weight: 300;
    }

    .login_card,
    .panel_card {
        width: 320px;
        border: 1px solid #3f4153;
        background-color: #181825;
        border-radius: 4px;
        padding: 24px;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        gap: 16px;
    }

    .card_header {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .panel_header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .title {
        color: #c9d2f0;
        font-size: 1.1rem;
        font-weight: 600;
    }

    .subtitle {
        color: #8588a3;
        font-size: 0.85rem;
        font-weight: 300;
    }

    .separator {
        width: 100%;
        height: 1px;
        background-color: #3f4153;
        flex-shrink: 0;
    }

    .error_msg {
        color: #f38ba8;
        font-size: 0.8rem;
        font-weight: 400;
    }

    .login_form {
        display: flex;
        flex-direction: column;
        gap: 10px;
        width: 100%;
    }

    .input_field input {
        width: 100%;
        box-sizing: border-box;
        background-color: #11111b;
        border: 1px solid #3f4153;
        border-radius: 4px;
        padding: 8px 12px;
        color: #c9d2f0;
        font-family: inherit;
        font-size: 0.85rem;
        outline: none;
        transition: border-color 0.15s ease;
    }

    .input_field input:focus {
        border-color: #9399b2;
    }

    .input_field input::placeholder {
        color: #565970;
    }

    .submit_btn {
        width: 100%;
        padding: 8px 12px;
        background-color: transparent;
        border: 1px solid #3f4153;
        border-radius: 4px;
        color: #c9d2f0;
        font-family: inherit;
        font-size: 0.85rem;
        font-weight: 400;
        cursor: pointer;
        transition: background-color 0.15s ease, border-color 0.15s ease, color 0.15s ease;
    }

    .submit_btn:hover:not(:disabled) {
        background-color: #29293b;
        color: #ffffff;
        border-color: #565970;
    }

    .submit_btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .logout_btn {
        padding: 4px 10px;
        background-color: transparent;
        border: 1px solid #3f4153;
        border-radius: 4px;
        color: #f38ba8;
        font-family: inherit;
        font-size: 0.8rem;
        cursor: pointer;
        transition: background-color 0.15s ease, border-color 0.15s ease;
    }

    .logout_btn:hover {
        background-color: #29293b;
        border-color: #f38ba8;
    }

    .panel_body {
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .desc {
        color: #aeb9df;
        font-size: 0.85rem;
        font-weight: 300;
    }

    .back_link {
        color: #8588a3;
        font-size: 0.85rem;
        font-weight: 300;
        text-decoration: none;
        transition: color 0.15s ease;
    }

    .back_link:hover {
        color: #ffffff;
        text-decoration: underline;
        text-underline-offset: 3px;
    }
</style>
