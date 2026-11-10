CREATE TABLE IF NOT EXISTS notifications(
    id UUID DEFAULT uuidv7(),
    actor_id UUID NOT NULL,
    recipient_id UUID NOT NULL,
    type TEXT NOT NULL,
    post_id UUID,           -- follow 沒有 post，可為 NULL
    read_at TIMESTAMPTZ,    -- NULL = 未讀
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY(id),
    FOREIGN KEY(actor_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(recipient_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CHECK(type IN ('post_liked', 'user_followed', 'new_post'))
);  