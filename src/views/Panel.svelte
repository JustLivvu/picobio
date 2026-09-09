<script lang="ts">
    import { api } from "../lib/api";
    import Icon from "@iconify/svelte";

    interface Post {
        id: string;
        title: string;
        content: string;
        createdAt: string;
    }

    type PanelTab = "profile" | "posts";

    let username = $state("");
    let password = $state("");
    let isLoading = $state(false);
    let errorMessage = $state("");
    let isAuthenticated = $state(false);
    let checkingAuth = $state(true);
    let activeTab = $state<PanelTab>("profile");
    let displayName = $state("");
    let pronouns = $state("");
    let description = $state("");
    let githubUrl = $state("");
    let joinedDate = $state("");
    let avatarUrl = $state("");
    let avatarFile = $state<File | null>(null);
    let avatarPreview = $state("");
    let posts = $state<Post[]>([]);
    let newPostTitle = $state("");
    let newPostContent = $state("");
    let isSavingProfile = $state(false);
    let isCreatingPost = $state(false);
    let statusMessage = $state("");

    async function checkAuth() {
        const token = sessionStorage.getItem("panel_token");
        if (!token) { checkingAuth = false; return; }
        try {
            const res = await api("/api/auth/verify", {
                headers: { Authorization: `Bearer ${token}` }
            });
            if (res.ok) {
                const data = await res.json();
                if (data.valid) { isAuthenticated = true; await loadData(token); }
                else { sessionStorage.removeItem("panel_token"); }
            } else { sessionStorage.removeItem("panel_token"); }
        } catch { sessionStorage.removeItem("panel_token"); }
        finally { checkingAuth = false; }
    }

    async function loadData(token: string) {
        try {
            const [profileRes, postsRes] = await Promise.all([api("/api/profile"), api("/api/posts")]);
            if (profileRes.ok) {
                const p = await profileRes.json();
                displayName = p.displayName || "";
                pronouns = p.pronouns || "";
                description = p.description || "";
                githubUrl = p.githubUrl || "";
                joinedDate = p.joinedDate || "";
                avatarUrl = p.avatarUrl || "";
                avatarPreview = p.avatarUrl ? `http://localhost:3000${p.avatarUrl}` : "";
            }
            if (postsRes.ok) { posts = await postsRes.json(); }
        } catch {}
    }

    async function handleLogin(e: SubmitEvent) {
        e.preventDefault();
        errorMessage = "";
        isLoading = true;
        try {
            const res = await api("/api/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password })
            });
            const data = await res.json();
            if (res.ok && data.success && data.token) {
                sessionStorage.setItem("panel_token", data.token);
                isAuthenticated = true;
                password = "";
                await loadData(data.token);
            } else { errorMessage = data.message || "Invalid credentials"; }
        } catch { errorMessage = "Server connection error"; }
        finally { isLoading = false; }
    }

    async function handleLogout() {
        const token = sessionStorage.getItem("panel_token");
        if (token) {
            try {
                await api("/api/auth/logout", {
                    method: "POST",
                    headers: { Authorization: `Bearer ${token}` }
                });
            } catch {}
        }
        sessionStorage.removeItem("panel_token");
        isAuthenticated = false;
        username = "";
        password = "";
    }

    function handleAvatarChange(e: Event) {
        const target = e.target as HTMLInputElement;
        if (target.files && target.files[0]) {
            avatarFile = target.files[0];
            avatarPreview = URL.createObjectURL(avatarFile);
        }
    }

    async function saveProfile(e: SubmitEvent) {
        e.preventDefault();
        const token = sessionStorage.getItem("panel_token");
        if (!token) return;
        isSavingProfile = true;
        statusMessage = "";
        try {
            if (avatarFile) {
                const formData = new FormData();
                formData.append("avatar", avatarFile);
                const avatarRes = await api("/api/profile/avatar", {
                    method: "POST",
                    headers: { Authorization: `Bearer ${token}` },
                    body: formData
                });
                if (avatarRes.ok) {
                    const aData = await avatarRes.json();
                    avatarUrl = aData.avatarUrl;
                    avatarFile = null;
                }
            }
            const res = await api("/api/profile", {
                method: "POST",
                headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
                body: JSON.stringify({ displayName, pronouns, description, githubUrl, joinedDate, avatarUrl })
            });
            if (res.ok) {
                statusMessage = "Profile saved";
                setTimeout(() => { statusMessage = ""; }, 3000);
            }
        } catch { statusMessage = "Failed to save"; }
        finally { isSavingProfile = false; }
    }

    async function createPost(e: SubmitEvent) {
        e.preventDefault();
        const token = sessionStorage.getItem("panel_token");
        if (!token || !newPostTitle.trim()) return;
        isCreatingPost = true;
        try {
            const res = await api("/api/posts", {
                method: "POST",
                headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
                body: JSON.stringify({ title: newPostTitle, content: newPostContent })
            });
            if (res.ok) {
                const created = await res.json();
                posts = [created, ...posts];
                newPostTitle = "";
                newPostContent = "";
            }
        } catch {}
        finally { isCreatingPost = false; }
    }

    async function deletePost(id: string) {
        const token = sessionStorage.getItem("panel_token");
        if (!token) return;
        try {
            const res = await api(`/api/posts?id=${id}`, {
                method: "DELETE",
                headers: { Authorization: `Bearer ${token}` }
            });
            if (res.ok) { posts = posts.filter(p => p.id !== id); }
        } catch {}
    }

    checkAuth();
</script>

<div class="panel_main">
    {#if checkingAuth}
        <div class="status_text">Loading...</div>
    {:else if isAuthenticated}
        <div class="panel_dashboard">
            <div class="panel_header">
                <div class="nav_tabs">
                    <button type="button" class={activeTab === "profile" ? "tab-chosen" : "tab"} onclick={() => (activeTab = "profile")}>Profile</button>
                    <button type="button" class={activeTab === "posts" ? "tab-chosen" : "tab"} onclick={() => (activeTab = "posts")}>Posts ({posts.length})</button>
                </div>
                <div class="header_actions">
                    <a href="/" class="view_site_link">View site</a>
                    <button type="button" class="logout_btn" onclick={handleLogout}>Logout</button>
                </div>
            </div>

            <div class="separator"></div>

            {#if statusMessage}
                <div class="status_banner">{statusMessage}</div>
            {/if}

            {#if activeTab === "profile"}
                <form class="profile_form" onsubmit={saveProfile}>
                    <div class="avatar_section">
                        <div class="avatar_box">
                            {#if avatarPreview}
                                <img src={avatarPreview} alt="Avatar" />
                            {:else}
                                <div class="avatar_placeholder">
                                    <Icon icon="material-symbols:person-outline" />
                                </div>
                            {/if}
                        </div>
                        <div class="avatar_controls">
                            <label class="upload_btn">
                                <span>Change image</span>
                                <input type="file" accept="image/*" onchange={handleAvatarChange} style="display:none;" />
                            </label>
                            <span class="hint">Square JPG or PNG</span>
                        </div>
                    </div>

                    <div class="separator"></div>

                    <div class="form_grid">
                        <div class="field">
                            <label for="p_display">Display name</label>
                            <input id="p_display" type="text" bind:value={displayName} placeholder="Livvya" required />
                        </div>
                        <div class="field">
                            <label for="p_pronouns">Pronouns</label>
                            <input id="p_pronouns" type="text" bind:value={pronouns} placeholder="livvya · she/her" />
                        </div>
                        <div class="field">
                            <label for="p_github">GitHub URL</label>
                            <input id="p_github" type="text" bind:value={githubUrl} placeholder="https://github.com/..." />
                        </div>
                        <div class="field">
                            <label for="p_joined">Joined date</label>
                            <input id="p_joined" type="text" bind:value={joinedDate} placeholder="Joined on Jun 14, 2026" />
                        </div>
                    </div>

                    <div class="field">
                        <label for="p_desc">Description / Bio</label>
                        <textarea id="p_desc" rows="4" bind:value={description} placeholder="Tell something about yourself..."></textarea>
                    </div>

                    <button type="submit" class="submit_btn" disabled={isSavingProfile}>
                        {isSavingProfile ? "Saving..." : "Save changes"}
                    </button>
                </form>

            {:else if activeTab === "posts"}
                <div class="posts_management">
                    <form class="new_post_form" onsubmit={createPost}>
                        <div class="field">
                            <label for="post_title">Title</label>
                            <input id="post_title" type="text" bind:value={newPostTitle} placeholder="Post title..." required />
                        </div>
                        <div class="field">
                            <label for="post_content">Content</label>
                            <textarea id="post_content" rows="3" bind:value={newPostContent} placeholder="Write post content..."></textarea>
                        </div>
                        <button type="submit" class="submit_btn" disabled={isCreatingPost}>
                            {isCreatingPost ? "Publishing..." : "Add post"}
                        </button>
                    </form>

                    <div class="separator"></div>

                    <div class="posts_list">
                        {#if posts.length === 0}
                            <div class="empty_state">No posts yet.</div>
                        {:else}
                            {#each posts as post (post.id)}
                                <div class="post_card">
                                    <div class="post_info">
                                        <span class="post_title">{post.title}</span>
                                        <span class="post_date">{post.createdAt}</span>
                                        {#if post.content}
                                            <p class="post_desc">{post.content}</p>
                                        {/if}
                                    </div>
                                    <button type="button" class="delete_btn" onclick={() => deletePost(post.id)}>Delete</button>
                                </div>
                            {/each}
                        {/if}
                    </div>
                </div>
            {/if}
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
                    <input type="text" bind:value={username} placeholder="Username" required autocomplete="username" />
                </div>
                <div class="input_field">
                    <input type="password" bind:value={password} placeholder="Password" required autocomplete="current-password" />
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
    .status_text { color: #8588a3; font-size: 0.9rem; font-weight: 300; }
    .login_card {
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
    .card_header { display: flex; flex-direction: column; gap: 4px; }
    .title { color: #c9d2f0; font-size: 1.1rem; font-weight: 600; }
    .subtitle { color: #8588a3; font-size: 0.85rem; font-weight: 300; }
    .separator { width: 100%; height: 1px; background-color: #3f4153; flex-shrink: 0; }
    .error_msg { color: #f38ba8; font-size: 0.8rem; }
    .login_form { display: flex; flex-direction: column; gap: 10px; width: 100%; }
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
    .input_field input:focus { border-color: #9399b2; }
    .input_field input::placeholder { color: #565970; }
    .panel_dashboard {
        width: 680px;
        max-width: calc(100vw - 40px);
        background-color: #181825;
        border: 1px solid #3f4153;
        border-radius: 4px;
        padding: 24px 30px;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        gap: 20px;
    }
    .panel_header { display: flex; justify-content: space-between; align-items: center; }
    .nav_tabs { display: flex; gap: 16px; }
    .tab-chosen {
        padding: 4px 0;
        color: #c9d2f0;
        font-weight: 500;
        font-size: 0.9rem;
        border: none;
        border-bottom: 2px solid #9399b2;
        background: transparent;
        cursor: pointer;
        font-family: inherit;
    }
    .tab {
        padding: 4px 0;
        color: #8588a3;
        font-weight: 300;
        font-size: 0.9rem;
        border: none;
        background: transparent;
        cursor: pointer;
        font-family: inherit;
        transition: color 0.15s ease;
    }
    .tab:hover { color: #c9d2f0; }
    .header_actions { display: flex; align-items: center; gap: 12px; }
    .view_site_link { color: #8588a3; font-size: 0.85rem; text-decoration: none; transition: color 0.15s ease; }
    .view_site_link:hover { color: #c9d2f0; text-decoration: underline; }
    .logout_btn {
        padding: 4px 8px;
        background-color: transparent;
        border: 1px solid #3f4153;
        border-radius: 4px;
        color: #f38ba8;
        font-family: inherit;
        font-size: 0.8rem;
        cursor: pointer;
        transition: background-color 0.15s ease, border-color 0.15s ease;
    }
    .logout_btn:hover { background-color: #29293b; border-color: #f38ba8; }
    .status_banner { color: #a6e3a1; font-size: 0.85rem; }
    .profile_form { display: flex; flex-direction: column; gap: 18px; }
    .avatar_section { display: flex; align-items: center; gap: 20px; }
    .avatar_box {
        width: 72px;
        height: 72px;
        border-radius: 4px;
        border: 1px solid #3f4153;
        background-color: #11111b;
        overflow: hidden;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }
    .avatar_box img { width: 100%; height: 100%; object-fit: cover; }
    .avatar_placeholder { color: #8588a3; font-size: 2rem; display: flex; }
    .avatar_controls { display: flex; flex-direction: column; gap: 6px; }
    .upload_btn {
        display: inline-block;
        width: fit-content;
        padding: 6px 12px;
        background-color: #11111b;
        border: 1px solid #3f4153;
        border-radius: 4px;
        color: #c9d2f0;
        font-size: 0.8rem;
        cursor: pointer;
        transition: border-color 0.15s ease;
        font-family: inherit;
    }
    .upload_btn:hover { border-color: #9399b2; }
    .hint { color: #565970; font-size: 0.75rem; }
    .form_grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
    .field { display: flex; flex-direction: column; gap: 6px; }
    .field label { color: #8588a3; font-size: 0.8rem; font-weight: 300; }
    .field input, .field textarea {
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
    .field input:focus, .field textarea:focus { border-color: #9399b2; }
    .field textarea { resize: vertical; }
    .field input::placeholder, .field textarea::placeholder { color: #565970; }
    .submit_btn {
        padding: 8px 16px;
        background-color: #11111b;
        border: 1px solid #3f4153;
        border-radius: 4px;
        color: #c9d2f0;
        font-family: inherit;
        font-size: 0.85rem;
        cursor: pointer;
        align-self: flex-start;
        transition: background-color 0.15s ease, border-color 0.15s ease;
    }
    .submit_btn:hover:not(:disabled) { background-color: #29293b; border-color: #9399b2; }
    .submit_btn:disabled { opacity: 0.6; cursor: not-allowed; }
    .posts_management { display: flex; flex-direction: column; gap: 20px; }
    .new_post_form { display: flex; flex-direction: column; gap: 12px; }
    .posts_list { display: flex; flex-direction: column; gap: 10px; }
    .empty_state { color: #565970; font-size: 0.85rem; font-weight: 300; }
    .post_card {
        padding: 12px;
        background-color: #11111b;
        border: 1px solid #3f4153;
        border-radius: 4px;
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 16px;
    }
    .post_info { display: flex; flex-direction: column; gap: 4px; }
    .post_title { color: #c9d2f0; font-size: 0.9rem; font-weight: 500; }
    .post_date { color: #565970; font-size: 0.75rem; }
    .post_desc { color: #8588a3; font-size: 0.85rem; margin: 4px 0 0 0; white-space: pre-wrap; }
    .delete_btn {
        background: transparent;
        border: none;
        color: #f38ba8;
        font-size: 0.8rem;
        cursor: pointer;
        padding: 2px 0;
        font-family: inherit;
        opacity: 0.8;
        transition: opacity 0.15s ease;
        flex-shrink: 0;
    }
    .delete_btn:hover { opacity: 1; text-decoration: underline; }
    .back_link { color: #8588a3; font-size: 0.85rem; font-weight: 300; text-decoration: none; transition: color 0.15s ease; }
    .back_link:hover { color: #ffffff; text-decoration: underline; text-underline-offset: 3px; }
</style>