<script lang="ts">
    interface Post {
        id: string;
        title: string;
        content: string;
        createdAt: string;
    }

    let { posts = [] }: { posts: Post[] } = $props();

    let searchQuery = $state('');

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
    <input
        type="text"
        class="search-input"
        placeholder="Search for posts"
        bind:value={searchQuery}
    />

    {#if filtered.length === 0}
        <div class="empty-state">
            <p class="title">No posts yet</p>
            <p class="subtitle">Check back later for new updates and posts.</p>
        </div>
    {:else}
        <div class="posts-list">
            {#each filtered as post (post.id)}
                <div class="post-card">
                    <div class="post-header">
                        <span class="post-title">{post.title}</span>
                        <span class="post-date">{post.createdAt}</span>
                    </div>
                    {#if post.content}
                        <p class="post-content">{post.content}</p>
                    {/if}
                </div>
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

    .search-input {
        width: 100%;
        box-sizing: border-box;
        padding: 10px 14px;
        background-color: #181825;
        border: 1px solid #3f4153;
        border-radius: 4px;
        color: #c9d2f0;
        font-size: 0.9rem;
        font-family: inherit;
        outline: none;
        transition: border-color 0.15s ease;
    }

    .search-input:focus {
        border-color: #9399b2;
    }

    .search-input::placeholder {
        color: #8588a3;
        font-weight: 300;
    }

    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        color: #8588a3;
    }

    .title {
        font-size: 1rem;
        font-weight: 600;
        color: #c9d2f0;
        margin: 0 0 6px 0;
    }

    .subtitle {
        font-size: 0.9rem;
        color: #8588a3;
        margin: 0;
    }

    .posts-list {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .post-card {
        padding: 16px;
        background-color: #181825;
        border: 1px solid #3f4153;
        border-radius: 4px;
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .post-header {
        display: flex;
        justify-content: space-between;
        align-items: baseline;
        gap: 12px;
    }

    .post-title {
        color: #c9d2f0;
        font-size: 0.95rem;
        font-weight: 500;
    }

    .post-date {
        color: #565970;
        font-size: 0.75rem;
        white-space: nowrap;
        flex-shrink: 0;
    }

    .post-content {
        color: #aeb9df;
        font-size: 0.85rem;
        font-weight: 300;
        margin: 0;
        line-height: 1.5;
        white-space: pre-wrap;
    }
</style>

