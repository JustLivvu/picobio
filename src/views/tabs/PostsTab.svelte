<script lang="ts">
    import Icon from "@iconify/svelte";

    interface Profile {
        displayName?: string;
        pronouns?: string;
        description?: string;
        joinedDate?: string;
        githubUrl?: string;
        avatarUrl?: string;
    }

    interface Post {
        id: string;
        title: string;
        content: string;
        createdAt: string;
    }

    let { posts = [], profile, avatarUrl = '' }: { posts: Post[]; profile?: Profile; avatarUrl?: string } = $props();

    let searchQuery = $state('');
    let hostname = $state(typeof window !== 'undefined' ? window.location.hostname : '');

    let filtered = $derived(
        searchQuery.trim()
            ? posts.filter(p =>
                p.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                p.content.toLowerCase().includes(searchQuery.toLowerCase())
              )
            : posts
    );
</script>

<div class="posts-container">
    <div class="search-wrapper">
        <Icon icon="material-symbols:search-rounded" class="search-icon" />
        <input
            type="text"
            class="search-input"
            placeholder="Search posts..."
            bind:value={searchQuery}
        />
    </div>

    {#if filtered.length === 0}
        <div class="empty-state">
            <Icon icon="material-symbols:post-add-rounded" class="empty-icon" />
            <span class="empty-title">No posts yet</span>
            <span class="empty-sub">Check back later for updates.</span>
        </div>
    {:else}
        <div class="posts-list">
            {#each filtered as post (post.id)}
                <article class="post-card">
                    <div class="post-header">
                        <div class="author-info">
                            <div class="author-avatar">
                                {#if avatarUrl}
                                    <img src={avatarUrl} alt={profile?.displayName || 'User'} />
                                {:else}
                                    <div class="avatar-placeholder">
                                        <Icon icon="material-symbols:person" />
                                    </div>
                                {/if}
                            </div>
                            <div class="author-meta">
                                <span class="author-name">{profile?.displayName || 'User'}</span>
                                <span class="author-handle">{(profile?.displayName || 'livvya').toLowerCase()}@{hostname}</span>
                            </div>
                        </div>
                        <span class="post-date">{post.createdAt}</span>
                    </div>

                    <div class="post-body">
                        <h3 class="post-title">{post.title}</h3>
                        {#if post.content}
                            <p class="post-content">{post.content}</p>
                        {/if}
                    </div>
                </article>
            {/each}
        </div>
    {/if}
</div>

<style>
    .posts-container {
        padding: 24px 0;
        width: 100%;
        display: flex;
        flex-direction: column;
        gap: 20px;
    }

    .search-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        width: 100%;
        background-color: #181825;
        border: 1px solid #3f4153;
        border-radius: 8px;
        padding: 0 12px;
        box-sizing: border-box;
        transition: border-color 0.15s ease;
    }

    .search-wrapper:focus-within {
        border-color: #9399b2;
    }

    .search-wrapper :global(.search-icon) {
        width: 18px;
        height: 18px;
        color: #565970;
        flex-shrink: 0;
    }

    .search-input {
        width: 100%;
        box-sizing: border-box;
        padding: 10px 8px;
        background-color: transparent;
        border: none;
        color: #c9d2f0;
        font-size: 0.9rem;
        font-family: inherit;
        outline: none;
    }

    .search-input::placeholder {
        color: #565970;
        font-weight: 300;
    }

    .empty-state {
        padding: 48px 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 8px;
    }

    .empty-state :global(.empty-icon) {
        width: 40px;
        height: 40px;
        color: #565970;
        margin-bottom: 4px;
    }

    .empty-title {
        color: #c9d2f0;
        font-size: 1rem;
        font-weight: 500;
    }

    .empty-sub {
        color: #565970;
        font-size: 0.85rem;
        font-weight: 300;
    }

    .posts-list {
        display: flex;
        flex-direction: column;
        gap: 16px;
    }

    .post-card {
        display: flex;
        flex-direction: column;
        gap: 14px;
        background-color: #181825;
        border-radius: 10px;
        padding: 18px 20px;
    }

    .post-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .author-info {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .author-avatar {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        overflow: hidden;
        flex-shrink: 0;
        background-color: #2a2b3d;
        border: 1px solid #3f4153;
    }

    .author-avatar img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        display: block;
    }

    .avatar-placeholder {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #8588a3;
    }

    .avatar-placeholder :global(svg) {
        width: 24px;
        height: 24px;
    }

    .author-meta {
        display: flex;
        flex-direction: column;
        gap: 1px;
    }

    .author-name {
        color: #c9d2f0;
        font-weight: 600;
        font-size: 0.95rem;
    }

    .author-handle {
        color: #8588a3;
        font-size: 0.75rem;
        font-weight: 400;
    }

    .post-date {
        color: #8588a3;
        font-size: 0.8rem;
        font-weight: 400;
    }

    .post-body {
        display: flex;
        flex-direction: column;
        gap: 6px;
    }

    .post-title {
        color: #f2f5f9;
        font-size: 1.05rem;
        font-weight: 600;
        margin: 0;
        line-height: 1.4;
    }

    .post-content {
        color: #bac2de;
        font-size: 0.9rem;
        font-weight: 300;
        margin: 0;
        line-height: 1.6;
        white-space: pre-wrap;
        word-break: break-word;
    }
</style>
