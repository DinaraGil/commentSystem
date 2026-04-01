CREATE TABLE IF NOT EXISTS user_data (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(32) NOT NULL
);

CREATE TABLE IF NOT EXISTS post (
    post_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    allow_comment BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES user_data (user_id)
);

CREATE TABLE IF NOT EXISTS comment (
    comment_id SERIAL PRIMARY KEY,
    reply_comment_id INTEGER,
    comment_level INTEGER,
    post_id INTEGER NOT NULL,
    user_id INTEGER,
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES user_data (user_id),
    FOREIGN KEY (post_id) REFERENCES post (post_id),
    FOREIGN KEY (reply_comment_id) REFERENCES comment (comment_id)
);