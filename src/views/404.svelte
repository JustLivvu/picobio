<script lang="ts">
</script>

<div class="not_found_main">
    <div class="content">
        <span class="code">404</span>
        <div class="divider"></div>
        <span class="message">This page could not be found.</span>
    </div>
    <a href="/" class="back_link">← Return home</a>
</div>

<style>
    .not_found_main {
        width: 100%;
        height: 100vh;
        background-color: #11111b;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        gap: 20px;
        box-sizing: border-box;
    }

    .content {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 20px;
    }

    .code {
        color: #c9d2f0;
        font-size: 1.5rem;
        font-weight: 600;
        line-height: 1;
    }

    .divider {
        width: 1px;
        height: 30px;
        background-color: #3f4153;
    }

    .message {
        color: #8588a3;
        font-size: 0.95rem;
        font-weight: 300;
    }

    .back_link {
        color: #aeb9df;
        font-size: 0.85rem;
        font-weight: 400;
        text-decoration: none;
        transition: color 0.15s ease;
    }

    .back_link:hover {
        color: #ffffff;
        text-decoration: underline;
        text-underline-offset: 3px;
    }
</style>
