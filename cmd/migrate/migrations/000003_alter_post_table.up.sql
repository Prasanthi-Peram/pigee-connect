ALTER TABLE
    posts
ADD CONSTRAINT fk_user FOREIGN KEY (users_id) REFERENCES users(id);