CREATE TABLE IF NOT EXISTS post_likes(
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE, --刪除貼文時自動清掉按讚記錄
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE  --刪除使用者時自動清掉按讚記錄
);