create table user_password_change_tokens (
  user_id integer not null primary key references users(id) on delete cascade,
  created_at timestamp not null,
  token varchar(32) not null,
  invite boolean not null default false
);

create index on user_password_change_tokens(invite, created_at);

