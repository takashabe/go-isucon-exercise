-- user
create index idx_email_and_passhash on user(email, passhash);

-- follow
create index idx_follow_id on follow(follow_id);

-- tweet
create index idx_user_id_and_created_at on tweet(user_id, created_at);
