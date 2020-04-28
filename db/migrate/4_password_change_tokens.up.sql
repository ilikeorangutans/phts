create table user_password_change_tokens (
  user_id integer not null primary key references users(id) on delete cascade,
  created_at timestamp not null,
  token varchar(32) not null
);

create index on user_password_change_tokens(created_at);

